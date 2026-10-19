package engine

import (
	"errors"
	"fmt"
)

// errDangerNotAllowed – thao tác nguy hiểm (script.DangerOps) khi người dùng chưa xác nhận: lỗi code permission (errCode).
var errDangerNotAllowed = errors.New("chưa cho phép")

// dangerErr – lỗi cho thao tác what khi chưa có Options.AllowDanger.
func dangerErr(what string) error {
	return fmt.Errorf("%w %s: chưa được xác nhận – bấm Chạy trong app để xem và đồng ý", errDangerNotAllowed, what)
}
