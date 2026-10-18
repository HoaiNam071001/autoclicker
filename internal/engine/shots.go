package engine

import (
	"fmt"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

const (
	maxShots = 5    // ảnh lúc lỗi tối đa mỗi lần chạy (vòng thử lại qua cổng fail không chụp mãi)
	shotMaxW = 1600 // ảnh rộng hơn thì thu nhỏ (Retina) cho file nhẹ
)

// shoot chụp màn chứa vùng kịch bản lúc node thất bại (tối đa maxShots mỗi lần chạy); lỗi chụp thì bỏ qua.
func (e *Engine) shoot(res *model.StepResult) {
	if e.shots == nil || len(e.shots) >= maxShots || res.Error == nil || res.Error.Code == model.ErrPermission || res.Error.Code == model.ErrEnd {
		return
	}
	img, err := e.scr.Capture(e.searchArea(model.Step{SearchArea: "screen"}))
	if err != nil {
		return
	}
	res.Shot = fmt.Sprintf("%d.png", len(e.shots)+1)
	e.shots[res.Shot] = vision.Shrink(img, shotMaxW)
}
