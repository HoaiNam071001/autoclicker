//go:build probe

package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-vgo/robotgo"

	"ghosthand/internal/device"
	"ghosthand/internal/model"
)

// checkLag – sau một click (trang đổi màu nền ngay trong handler), chụp liên tục tới khi ảnh chụp thấy màu mới:
// đo độ trễ "ảnh chụp cũ" (bước 33: mẫu game bấm lại chỗ nút đã nhảy đi).
func checkLag(ctx context.Context, cal calibration) error {
	scr := device.RobotgoScreen{}
	c := point{200, 300}
	s := cal.toScreen(c)
	x, y := int(math.Round(s.X)), int(math.Round(s.Y))
	area := model.Rect{X: x + 40, Y: y, W: 20, H: 20} // cạnh điểm bấm, không dính con trỏ
	if *lagFull {
		area = device.MainScreen() // như OCR: chụp cả màn rồi xem điểm cạnh chỗ bấm
	}
	for i := 0; i < 8; i++ {
		before, err := sample(scr, area, x+50, y+10)
		if err != nil {
			return err
		}
		robotgo.Move(x, y)
		robotgo.Click("left", false)
		t0 := time.Now()
		got, polls := -1.0, 0
		for time.Since(t0) < 2*time.Second {
			now, err := sample(scr, area, x+50, y+10)
			if err != nil {
				return err
			}
			polls++
			if now != before {
				got = float64(time.Since(t0).Microseconds()) / 1000
				break
			}
		}
		fmt.Printf("  lần %d: ảnh chụp thấy màu mới sau %.0f ms (%d lần chụp)\n", i+1, got, polls)
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func sample(scr device.RobotgoScreen, r model.Rect, x, y int) (uint32, error) {
	img, err := scr.Capture(r)
	if err != nil {
		return 0, err
	}
	k := float64(img.Bounds().Dx()) / float64(r.W) // ảnh Retina = 2 px / điểm
	px := img.Bounds().Min.X + int(float64(x-r.X)*k)
	py := img.Bounds().Min.Y + int(float64(y-r.Y)*k)
	cr, cg, cb, _ := img.At(px, py).RGBA()
	return cr>>8<<16 | cg>>8<<8 | cb>>8, nil
}

// checkJump – nút nhảy chỗ sau mỗi lần trúng (như game.html): đọc rect qua CDP, di chuột + bấm (chờ settle), đếm trúng / trượt.
func checkJump(ctx context.Context, cal calibration) error {
	if _, err := evalJS[bool](ctx, `(() => { const b = document.createElement('button'); b.id = 'j'; b.textContent = 'BẮT TÔI';
		b.style.cssText = 'position:fixed;left:300px;top:300px;padding:18px 26px;font-size:24px';
		window.__hit = 0; window.__miss = 0;
		b.addEventListener('click', e => { e.stopPropagation(); __hit++; b.style.left = (100 + Math.random() * 800) + 'px'; b.style.top = (100 + Math.random() * 500) + 'px'; });
		document.body.appendChild(b); document.addEventListener('click', () => __miss++); return true; })()`); err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		c, err := evalJS[point](ctx, `(() => { const r = document.getElementById('j').getBoundingClientRect(); return {X: r.left + r.width/2, Y: r.top + r.height/2}; })()`)
		if err != nil {
			return err
		}
		s := cal.toScreen(c)
		robotgo.Move(int(math.Round(s.X)), int(math.Round(s.Y)))
		time.Sleep(time.Duration(*settleMs) * time.Millisecond)
		robotgo.Click("left", false)
		time.Sleep(400 * time.Millisecond)
	}
	r, err := evalJS[struct{ Hit, Miss int }](ctx, `({Hit: __hit, Miss: __miss})`)
	fmt.Printf("  settle %d ms: trúng %d, trượt %d\n", *settleMs, r.Hit, r.Miss)
	return err
}
