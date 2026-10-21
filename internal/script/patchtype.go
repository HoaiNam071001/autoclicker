package script

import (
	"encoding/json"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// Đổi loại node qua update_script (46 B2): dựng node mới theo loại mới (mặc định của loại đó) thay vì đè trường lên node
// cũ – không để sót trường của loại cũ. Giữ trường chung của node và trường cũ nào loại mới cũng dùng (points, selectors…).

// keepOnRetype – trường chung giữ khi đổi loại (commonFieldsDoc trong model/steptypes.go).
var keepOnRetype = []string{"name", "pos", "enabled", "delayAfter", "jitterMs", "retry", "breakpoint"}

// retypeStep – node old đổi sang loại trong raw (raw có "type" khác old.Type). id / ref giữ nguyên.
func retypeStep(old model.Step, raw json.RawMessage, newType string) (model.Step, error) {
	keep := slices.Clone(keepOnRetype)
	for _, d := range model.StepTypeDocs() {
		if d.Type == newType {
			for _, f := range d.Fields {
				keep = append(keep, strings.SplitN(f, "[", 2)[0]) // "points[1]" → points
			}
		}
	}
	b, err := json.Marshal(old)
	if err != nil {
		return model.Step{}, err
	}
	var from, patch map[string]json.RawMessage
	if err := json.Unmarshal(b, &from); err != nil {
		return model.Step{}, err
	}
	if err := json.Unmarshal(raw, &patch); err != nil {
		return model.Step{}, err
	}
	for _, k := range keep {
		if _, sent := patch[k]; !sent && from[k] != nil {
			patch[k] = from[k]
		}
	}
	patch["id"], patch["ref"] = from["id"], from["ref"]
	merged, err := json.Marshal(patch)
	if err != nil {
		return model.Step{}, err
	}
	return ExpandStep(merged, old.ID, old.Ref)
}

// dropBadPorts bỏ cạnh đi ra từ node s theo cổng loại mới không có (if → true / false, click → success / fail…).
func dropBadPorts(sc *model.Script, s model.Step) []model.Edge {
	var dropped []model.Edge
	sc.Edges = slices.DeleteFunc(sc.Edges, func(e model.Edge) bool {
		bad := e.From == s.ID && !model.HasPort(s, e.Port)
		if bad {
			dropped = append(dropped, e)
		}
		return bad
	})
	return dropped
}
