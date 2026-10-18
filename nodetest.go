package main

import (
	"errors"
	"fmt"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
)

// Nút "Thử node" (bước 20 F3): xét node Nếu / Switch / Switch theo màn hình trên màn hình lúc này (một lần, không chờ),
// báo cổng sẽ đi + từng phần đúng / sai và vì sao. Cùng logic với lúc chạy (cond.go); {{key}} của input lấy giá trị
// mặc định, biến node ({{n2…}}) chưa biết → known = false.

type NodeTest struct {
	Port  string     `json:"port"`  // cổng sẽ đi lúc này ("" = chưa biết)
	Known bool       `json:"known"` // false: phụ thuộc biến node, chỉ biết lúc chạy
	Parts []CondPart `json:"parts"`
	Value string     `json:"value,omitempty"` // switch: giá trị đã xét
	Ms    int        `json:"ms"`
}

// TestNode thử node stepID của kịch bản sc.
func (a *App) TestNode(sc model.Script, stepID string) (NodeTest, error) {
	var s *model.Step
	for i := range sc.Steps {
		if sc.Steps[i].ID == stepID {
			s = &sc.Steps[i]
		}
	}
	if s == nil {
		return NodeTest{}, errors.New("không có node này")
	}
	n, err := withInputDefaults(sc, *s)
	if err != nil {
		return NodeTest{}, err
	}
	e := newEngine(nil, sc, nil)
	e.screen, e.displays = device.MainScreen(), device.Displays()
	var res NodeTest
	run := func() error {
		start := time.Now()
		var err error
		res, err = e.testNode(n)
		res.Ms = int(time.Since(start).Milliseconds())
		return err
	}
	if n.Type == model.StepSwitch { // không nhìn màn hình – khỏi ẩn cửa sổ
		return res, run()
	}
	return res, a.withHidden(run)
}

// withInputDefaults – bản sao node với {{key}} của input thay bằng giá trị mặc định (Thử tìm / Thử node).
func withInputDefaults(sc model.Script, s model.Step) (model.Step, error) {
	s = cloneStep(s)
	for _, f := range varFields(&s) {
		t, err := expandVars(*f.p, func(key string) (string, bool) {
			for _, in := range sc.Inputs {
				if in.Key == key && !isEmptyValue(in.Default) {
					return valueString(in.Default), true
				}
			}
			return "", false
		})
		if err != nil {
			return s, fmt.Errorf("thử dùng giá trị mặc định của input: %w", err)
		}
		*f.p = t
	}
	return s, nil
}

func (e *engine) testNode(s model.Step) (NodeTest, error) {
	res := NodeTest{Known: true, Parts: []CondPart{}}
	switch s.Type {
	case model.StepIf:
		if s.Cond == nil {
			return res, errors.New("node chưa có điều kiện")
		}
		ok, known, err := e.evalCond(*s.Cond, s.ID+"#cond", nil, newSnapshot(e.scr, e.ocr), &res.Parts)
		if err != nil {
			return res, err
		}
		res.Known = known
		if known {
			res.Port = map[bool]string{true: model.PortTrue, false: model.PortFalse}[ok]
		}
	case model.StepSwitchScreen:
		snap := newSnapshot(e.scr, e.ocr)
		res.Port = model.PortNone
		for i, c := range s.Cases {
			if c.Cond == nil {
				continue
			}
			p0 := len(res.Parts)
			ok, _, err := e.evalCond(*c.Cond, fmt.Sprintf("%s#case%d", s.ID, i), []int{i}, snap, &res.Parts)
			if err != nil {
				return res, fmt.Errorf("case %s: %w", caseName(c), err)
			}
			for j := p0; j < len(res.Parts); j++ {
				res.Parts[j].Case = c.Key
			}
			if ok && res.Port == model.PortNone {
				res.Port = model.CasePrefix + c.Key // vẫn xét hết các case để báo từng cái
			}
		}
	case model.StepSwitch:
		res.Value = s.Value
		if nodeVarRe.MatchString(s.Value) {
			res.Known = false
		}
		res.Port = model.PortDefault
		for i, c := range s.Cases {
			part := CondPart{Path: []int{i}, Case: c.Key, Kind: model.CondValue, Known: res.Known}
			op := map[string]string{"eq": "bằng", "contains": "chứa", "regex": "khớp mẫu"}[c.Op]
			part.Detail = fmt.Sprintf("“%s” %s “%s”", s.Value, op, c.Match)
			if res.Known {
				ok, err := matchCase(s.Value, c)
				if err != nil {
					return res, err
				}
				part.Ok = ok
				if ok && res.Port == model.PortDefault {
					res.Port = model.CasePrefix + c.Key
				}
			}
			res.Parts = append(res.Parts, part)
		}
		if !res.Known {
			res.Port = ""
		}
	default:
		return res, errors.New("chỉ thử được node Nếu / Switch")
	}
	return res, nil
}
