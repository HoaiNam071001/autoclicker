package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Nút "Thử điều kiện" (bước 20c): xét điều kiện của khối trên màn hình lúc này (một lần, không chờ), báo từng phần
// đúng / sai và vì sao – để người dùng biết điều kiện có đúng ý trước khi chạy. Logic so giống evalCond (blocks.go);
// {{biến}} lấy giá trị mặc định của input.

// CondPart – kết quả một điều kiện lá (ảnh / chữ / màu / input / bước trước).
type CondPart struct {
	Path   []int  `json:"path"` // vị trí trong conds lồng nhau ([] = chính điều kiện gốc)
	Kind   string `json:"kind"`
	Ok     bool   `json:"ok"`            // kết quả sau khi đã đảo (not)
	Known  bool   `json:"known"`         // false: chỉ biết lúc chạy (bước trước thành công)
	Detail string `json:"detail"`        // vd "thấy tại X 120 Y 40 · khớp 97%", "đọc được: …"
	Hit    *Rect  `json:"hit,omitempty"` // chỗ thấy (toạ độ robotgo) – overlay vẽ khung
}

type CondTest struct {
	Ok    bool       `json:"ok"`
	Known bool       `json:"known"` // false: có phần chỉ biết lúc chạy và nó quyết định kết quả
	Parts []CondPart `json:"parts"`
	Ms    int        `json:"ms"`
}

// findStepByID tìm bước ở mọi tầng của cây.
func findStepByID(steps []Step, id string) *Step {
	var out *Step
	walkSteps(steps, "", func(s *Step, _ string) error {
		if s.ID == id {
			out = s
			return errors.New("thấy")
		}
		return nil
	})
	return out
}

// TestCond xét điều kiện của bước stepID (path: điều kiện con, rỗng = cả điều kiện).
func (a *App) TestCond(sc Script, stepID string, path []int) (CondTest, error) {
	s := findStepByID(sc.Steps, stepID)
	if s == nil || s.Cond == nil {
		return CondTest{}, errors.New("bước không có điều kiện")
	}
	c := s.Cond
	for _, i := range path {
		if i < 0 || i >= len(c.Conds) {
			return CondTest{}, errors.New("không có điều kiện con này")
		}
		c = &c.Conds[i]
	}
	cond := cloneCond(*c)
	vars := func(t string) (string, error) {
		return expandVars(t, func(key string) (string, bool) {
			for _, in := range sc.Inputs {
				if in.Key == key && !isEmptyValue(in.Default) {
					return valueString(in.Default), true
				}
			}
			return "", false
		})
	}
	err := walkConds(&cond, func(x *Cond) error {
		for _, f := range []*string{&x.Left, &x.Right, &x.FindText} {
			t, err := vars(*f)
			if err != nil {
				return fmt.Errorf("thử điều kiện dùng giá trị mặc định của input: %w", err)
			}
			*f = t
		}
		return nil
	})
	if err != nil {
		return CondTest{}, err
	}
	e := newEngine(nil, sc, nil)
	e.screen, e.displays = mainScreen(), displays()
	var res CondTest
	err = a.withHidden(func() error {
		start := time.Now()
		var err error
		res.Ok, res.Known, err = e.testCond(cond, s.ID, path, &res.Parts)
		res.Ms = int(time.Since(start).Milliseconds())
		return err
	})
	return res, err
}

// testCond xét c một lần, ghi kết quả từng điều kiện lá vào parts. known = false khi kết quả phụ thuộc lúc chạy.
func (e *engine) testCond(c Cond, key string, path []int, parts *[]CondPart) (ok, known bool, err error) {
	part := CondPart{Path: append([]int{}, path...), Kind: c.Kind, Known: true}
	switch c.Kind {
	case CondAll, CondAny:
		ok, known = c.Kind == CondAll, true
		decided := false
		for i, sub := range c.Conds {
			v, k, err := e.testCond(sub, fmt.Sprintf("%s.%d", key, i), append(append([]int{}, path...), i), parts)
			if err != nil {
				return false, false, err
			}
			if decided {
				continue // đã biết kết quả – vẫn xét hết để báo từng phần
			}
			if !k {
				known = false
				continue
			}
			if v != (c.Kind == CondAll) { // all: gặp sai → sai; any: gặp đúng → đúng
				ok, known, decided = v, true, true
			}
		}
		return ok != c.Not, known, nil
	case CondResult:
		part.Known, part.Detail = false, "chỉ biết lúc chạy (bước chạy ngay trước khối)"
		*parts = append(*parts, part)
		return false, false, nil
	case CondInput:
		v, err := compareValues(c.Left, c.Op, c.Right)
		if err != nil {
			return false, false, err
		}
		ok = v
		part.Detail = fmt.Sprintf("“%s” %s “%s”", c.Left, map[string]string{"eq": "=", "ne": "≠", "contains": "chứa", "gt": ">", "lt": "<", "empty": "trống"}[c.Op], c.Right)
		if c.Op == "empty" {
			part.Detail = fmt.Sprintf("“%s” trống", c.Left)
		}
	case CondText:
		probe := c.probe(key)
		area := e.searchArea(probe)
		if err := captureArea(area, e.displays); err != nil {
			return false, false, err
		}
		hit, read, found, err := locateText(e.scr, e.ocr, area, probe)
		if err != nil {
			return false, false, err
		}
		ok = found
		if found {
			part.Hit, part.Detail = &hit, fmt.Sprintf("thấy “%s” tại X %d Y %d", c.FindText, hit.X+hit.W/2, hit.Y+hit.H/2)
		} else {
			part.Detail = "không thấy" + readHint(read)
		}
	case CondImage:
		probe := c.probe(key)
		if probe.Template == "" {
			return false, false, errors.New("điều kiện chưa có ảnh mẫu")
		}
		tpl, err := e.template(probe)
		if err != nil {
			return false, false, err
		}
		area := e.searchArea(probe)
		if err := captureArea(area, e.displays); err != nil {
			return false, false, err
		}
		img, err := e.scr.Capture(area)
		if err != nil {
			return false, false, err
		}
		m, k := matchScaled(img, area, tpl, probe.TemplateScale)
		score := int(math.Round(m.Score * 100))
		hit := Rect{X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k)}
		ok = score >= max(probe.Threshold, 1)
		part.Hit = &hit
		part.Detail = fmt.Sprintf("%s tại X %d Y %d · khớp %d%% (cần %d%%)", map[bool]string{true: "thấy", false: "giống nhất"}[ok], hit.X+hit.W/2, hit.Y+hit.H/2, score, probe.Threshold)
	case CondPixel:
		if len(c.Points) == 0 {
			return false, false, errors.New("điều kiện màu chưa chọn điểm")
		}
		want, err := parseColor(c.Color)
		if err != nil {
			return false, false, err
		}
		p := resolvePoint(c.Points[0], bounds(e.sc, e.screen))
		got, err := e.scr.Pixel(p.X, p.Y)
		if err != nil {
			return false, false, err
		}
		d := colorDiff(got, want)
		ok = d <= c.Tolerance
		part.Hit = &Rect{X: p.X - 6, Y: p.Y - 6, W: 12, H: 12}
		part.Detail = fmt.Sprintf("điểm X %d Y %d màu %s (cần %s ± %d, lệch %d)", p.X, p.Y, hexColor(got), strings.ToUpper(c.Color), c.Tolerance, d)
	default:
		return false, false, fmt.Errorf("loại điều kiện không hợp lệ %q", c.Kind)
	}
	part.Ok = ok != c.Not
	if c.Not {
		part.Detail += " → đảo"
	}
	*parts = append(*parts, part)
	return part.Ok, true, nil
}
