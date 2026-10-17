package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Khối điều kiện (bước 20): kịch bản là cây – bước if / group có bước con Children (và Else). Đánh số theo đường dẫn
// "3.2": nhánh Else đánh số tiếp sau Children (if 3 có 2 bước đúng → bước đầu nhánh sai là 3.3).

// isBlock – loại bước có bước con.
func isBlock(t string) bool { return t == StepIf || t == StepGroup }

// stepPath – đường dẫn của bước thứ i (từ 0) trong danh sách có tiền tố prefix ("" ở gốc, "3." trong khối 3).
func stepPath(prefix string, i int) string { return prefix + strconv.Itoa(i+1) }

// walkSteps gọi f cho mọi bước (kể cả bước con, theo thứ tự duyệt sâu), dừng ở lỗi đầu tiên.
func walkSteps(steps []Step, prefix string, f func(s *Step, path string) error) error {
	for i := range steps {
		s := &steps[i]
		path := stepPath(prefix, i)
		if err := f(s, path); err != nil {
			return err
		}
		if err := walkSteps(s.Children, path+".", f); err != nil {
			return err
		}
		if err := walkSteps(s.Else, path+".", func(c *Step, p string) error { return f(c, shiftPath(p, path, len(s.Children))) }); err != nil {
			return err
		}
	}
	return nil
}

// shiftPath cộng off vào thành phần ngay sau tiền tố parent của p ("3.1.2" với parent "3", off 2 → "3.3.2").
func shiftPath(p, parent string, off int) string {
	if off == 0 {
		return p
	}
	rest := strings.TrimPrefix(p, parent+".")
	head, tail, _ := strings.Cut(rest, ".")
	n, _ := strconv.Atoi(head)
	out := parent + "." + strconv.Itoa(n+off)
	if tail != "" {
		out += "." + tail
	}
	return out
}

// cloneSteps sao chép sâu (bước con, điểm, điều kiện) – resolveInputs sửa bản sao không đụng kịch bản gốc.
func cloneSteps(steps []Step) []Step {
	if steps == nil {
		return nil
	}
	out := make([]Step, len(steps))
	for i, s := range steps {
		s.Points = slices.Clone(s.Points)
		s.Keys = slices.Clone(s.Keys)
		s.Modifiers = slices.Clone(s.Modifiers)
		s.Children = cloneSteps(s.Children)
		s.Else = cloneSteps(s.Else)
		if s.Cond != nil {
			c := cloneCond(*s.Cond)
			s.Cond = &c
		}
		if s.OnFail != nil {
			f := *s.OnFail
			s.OnFail = &f
		}
		out[i] = s
	}
	return out
}

func cloneCond(c Cond) Cond {
	c.Points = slices.Clone(c.Points)
	if c.Conds != nil {
		cs := make([]Cond, len(c.Conds))
		for i, x := range c.Conds {
			cs[i] = cloneCond(x)
		}
		c.Conds = cs
	}
	return c
}

// walkConds gọi f cho c và mọi điều kiện con.
func walkConds(c *Cond, f func(*Cond) error) error {
	if err := f(c); err != nil {
		return err
	}
	for i := range c.Conds {
		if err := walkConds(&c.Conds[i], f); err != nil {
			return err
		}
	}
	return nil
}

// ---- Engine ----

// errBreak – bước Thoát khối lặp: khối group gần nhất bắt lấy.
var errBreak = errors.New("thoát khối lặp")

// stopSignal – bước Dừng kịch bản.
type stopSignal struct {
	id, path, msg string
	fail          bool
}

func (s *stopSignal) Error() string { return s.msg }

// stepFail – bước thất bại không được xử lý (không có / hết OnFail). Đi ngược lên qua các khối, giữ đường dẫn của
// bước gốc gây lỗi; khối cha có OnFail thì khối cha xử lý (vd thử lại cả nhóm).
type stepFail struct {
	id, path string
	err      error
}

func (f *stepFail) Error() string { return fmt.Sprintf("bước %s: %v", f.path, f.err) }
func (f *stepFail) Unwrap() error { return f.err }

// notFound – chờ ảnh / chữ / màu hết giờ: điều kiện coi là "sai", bước thường coi là thất bại.
type notFound struct{ msg string }

func (n *notFound) Error() string { return n.msg }

func isSignal(err error) bool {
	var st *stopSignal
	return errors.Is(err, errBreak) || errors.As(err, &st)
}

// runList chạy một danh sách bước (gốc hoặc thân khối). off: số cộng vào thứ tự để đánh số (nhánh Else).
func (e *engine) runList(ctx context.Context, steps []Step, prefix string, off int) error {
	jumps := 0
	for i := 0; i < len(steps); i++ {
		s := steps[i]
		if !s.Enabled {
			continue
		}
		path := stepPath(prefix, off+i)
		if prefix == "" {
			e.rootStep = e.rootOrder[i]
		}
		err := e.runWithRetry(ctx, s, path)
		if err == nil {
			if err = sleep(ctx, e.jitterMs(s.DelayAfter, s.JitterMs)); err != nil {
				return err
			}
			e.lastOK = true
			continue
		}
		if ctx.Err() != nil || isSignal(err) {
			return err
		}
		e.lastOK = false
		if f := s.OnFail; f != nil && f.Action != FailStop {
			e.miss = nil // đã xử lý theo "khi thất bại"
		}
		if f := s.OnFail; f != nil {
			switch f.Action {
			case FailNext:
				continue
			case FailElse:
				if err := e.runList(ctx, s.Else, path+".", 0); err != nil {
					return err
				}
				continue
			case FailGoto:
				if jumps++; jumps > maxJumps {
					return &stepFail{s.ID, path, fmt.Errorf("nhảy bước quá %d lần – kiểm tra “khi thất bại” có tạo vòng lặp không dừng không", maxJumps)}
				}
				i = slices.IndexFunc(steps, func(x Step) bool { return x.ID == f.Goto }) - 1 // vòng for tăng lên 1
				continue
			}
		}
		var sf *stepFail
		if errors.As(err, &sf) {
			return err // lỗi của bước con – giữ đường dẫn bước gốc
		}
		return &stepFail{s.ID, path, err}
	}
	return nil
}

// execStep chạy một bước: khối thì chạy bước con, bước thường thì runStep.
func (e *engine) execStep(ctx context.Context, s Step, path string) error {
	switch s.Type {
	case StepIf:
		ok, err := e.evalCond(ctx, *s.Cond, s.ID)
		if err != nil {
			return err
		}
		if ok {
			return e.runList(ctx, s.Children, path+".", 0)
		}
		return e.runList(ctx, s.Else, path+".", len(s.Children))
	case StepGroup:
		return e.runGroup(ctx, s, path)
	case StepBreak:
		return errBreak
	case StepStop:
		msg := s.Message
		if msg == "" {
			msg = "dừng theo bước Dừng kịch bản"
		}
		return &stopSignal{id: s.ID, path: path, msg: msg, fail: s.StopError}
	}
	return e.runStep(ctx, s)
}

func (e *engine) runGroup(ctx context.Context, s Step, path string) error {
	maxL := s.MaxLoops
	if maxL <= 0 {
		maxL = defaultMaxLoops
	}
	total := "∞"
	if s.LoopMode == LoopCount {
		total = strconv.Itoa(max(s.Repeat, 1))
	}
	name := s.Name
	if name == "" {
		name = "Bước " + path
	}
	e.blocks = append(e.blocks, "")
	defer func() {
		e.blocks = e.blocks[:len(e.blocks)-1]
		e.update(true, func(p *Progress) { p.Block = strings.Join(e.blocks, " › ") })
	}()
	for r := 1; ; r++ {
		if s.LoopMode == LoopCount {
			if r > max(s.Repeat, 1) {
				return nil
			}
		} else {
			ok, err := e.evalCond(ctx, *s.Cond, s.ID)
			if err != nil {
				return err
			}
			if s.LoopMode == LoopUntil && ok || s.LoopMode == LoopWhile && !ok {
				return nil
			}
			if r > maxL {
				want := map[string]string{LoopUntil: "vẫn chưa đúng", LoopWhile: "vẫn còn đúng"}[s.LoopMode]
				return fmt.Errorf("đã lặp %d lần (tối đa) mà điều kiện %s", maxL, want)
			}
		}
		e.blocks[len(e.blocks)-1] = fmt.Sprintf("%s · lặp %d / %s", name, r, total)
		e.update(true, func(p *Progress) { p.Block = strings.Join(e.blocks, " › ") })
		err := e.runList(ctx, s.Children, path+".", 0)
		if errors.Is(err, errBreak) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// evalCond xét điều kiện. Chờ ảnh / chữ / màu hết giờ = sai; lỗi khác (thiếu quyền, ảnh hỏng…) = lỗi của bước.
// key: id bước chứa điều kiện – để nhớ ảnh mẫu đã giải mã.
func (e *engine) evalCond(ctx context.Context, c Cond, key string) (bool, error) {
	var ok bool
	switch c.Kind {
	case CondImage, CondText, CondPixel:
		probe := c.probe(key)
		var err error
		if c.Kind == CondPixel {
			err = e.waitPixel(ctx, probe)
		} else {
			_, err = e.findImage(ctx, probe)
		}
		var nf *notFound
		switch {
		case err == nil:
			ok = true
		case errors.As(err, &nf):
			ok, e.miss = false, nil // "không thấy" của điều kiện không phải lỗi – đừng báo chỗ giống nhất cho bước khác
		default:
			return false, err
		}
	case CondInput:
		var err error
		if ok, err = compareValues(c.Left, c.Op, c.Right); err != nil {
			return false, err
		}
	case CondResult:
		ok = e.lastOK
	case CondAll, CondAny:
		ok = c.Kind == CondAll
		for i, sub := range c.Conds {
			v, err := e.evalCond(ctx, sub, fmt.Sprintf("%s.%d", key, i))
			if err != nil {
				return false, err
			}
			if v != ok { // all: gặp sai → sai; any: gặp đúng → đúng
				ok = v
				break
			}
		}
	default:
		return false, fmt.Errorf("loại điều kiện không hợp lệ %q", c.Kind)
	}
	return ok != c.Not, nil
}

// probe – bước Chờ ảnh / Chờ chữ / Chờ pixel tương ứng với điều kiện (dùng lại engine tìm ảnh).
func (c Cond) probe(key string) Step {
	s := Step{
		ID: key + "#cond", Type: StepWaitImage, Template: c.Template, TemplateScale: c.TemplateScale, Threshold: c.Threshold,
		SearchArea: c.SearchArea, FindText: c.FindText, TextExact: c.TextExact, TextNoAccent: c.TextNoAccent, Occurrence: c.Occurrence,
		Points: c.Points, Color: c.Color, Tolerance: c.Tolerance, TimeoutMs: c.WaitMs, IntervalMs: c.IntervalMs,
	}
	if s.IntervalMs <= 0 {
		s.IntervalMs = 250
	}
	switch c.Kind {
	case CondText:
		s.FindBy = FindByText
	case CondPixel:
		s.Type = StepWaitPixel
	}
	return s
}

// compareValues so 2 vế (đã thay {{biến}}). eq / ne / contains không phân biệt hoa thường, bỏ khoảng trắng 2 đầu;
// gt / lt so số.
func compareValues(left, op, right string) (bool, error) {
	l, r := strings.TrimSpace(left), strings.TrimSpace(right)
	switch op {
	case "eq":
		return strings.EqualFold(l, r), nil
	case "ne":
		return !strings.EqualFold(l, r), nil
	case "contains":
		return strings.Contains(strings.ToLower(l), strings.ToLower(r)), nil
	case "empty":
		return l == "", nil
	case "gt", "lt":
		a, err1 := strconv.ParseFloat(strings.ReplaceAll(l, ",", "."), 64)
		b, err2 := strconv.ParseFloat(strings.ReplaceAll(r, ",", "."), 64)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("so lớn / nhỏ cần 2 số, có “%s” và “%s”", l, r)
		}
		if op == "gt" {
			return a > b, nil
		}
		return a < b, nil
	}
	return false, fmt.Errorf("phép so không hợp lệ %q", op)
}
