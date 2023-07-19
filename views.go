package cbgo

import (
	"fmt"
	"sort"
	"strconv"
)

type ViewState struct {
	Serial	int
	WS	int
	View	int
}

func NewViewState(serial, ws, view int) *ViewState {
	return &ViewState{
		Serial: serial,
		WS: ws,
		View: view,
	}
}

func (v *ViewState) String() string {
	return fmt.Sprintf("[%d:%d] %d", v.WS, v.View, v.Serial)
}

func (v *ViewState) SwitchCmd() string {
	ret := "workspace " + strconv.Itoa(v.WS)
	if v.View != 0 {
		ret += "\n" + "view " + strconv.Itoa(v.View)
	}
	return ret
}

func (v *ViewState) Equal(ws, view int) bool {
	return v.WS == ws && v.View == view
}

type ViewStates []*ViewState

func (vs ViewStates) String() string {
	ret := "[ "

	for _, v := range(vs) {
		ret += fmt.Sprintf("%s ", v)
	}
	ret += "]"
	return ret
}

func (vs ViewStates) Len() int {
	return len(vs)
}

func (vs ViewStates) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ViewStates) Less(i, j int) bool {
	return vs[i].Serial < vs[j].Serial
}

func (vs ViewStates) Find(ws, view int) (ret *ViewState, idx int) {
	idx = -1
	for i, v := range(vs) {
		if v.Equal(ws, view) {
			idx = i
			ret = v
			break
		}
	}
	return
}

func (vs *ViewStates) Append(v *ViewState) bool {
	if v.View != 0 {
		vs.Pop(v.WS, 0)
	}
	*vs = append(*vs, v)
	return true
}

func (vs *ViewStates) Pop(ws, view int) {
	_, idx := vs.Find(ws, view); if idx != -1 {
		*vs = append((*vs)[:idx], (*vs)[idx+1:]...)
	}
}

func (vs *ViewStates) Push(serial, ws, view int) int {
	fmt.Printf("push(%d, %d, %d)\n", serial, ws, view)

	if len(*vs) == 0 {
		fmt.Printf("empty\n")
		if vs.Append(NewViewState(serial, ws, view)) {
			return serial + 1
		} else {
			return serial
		}
	} else {
		cur := (*vs)[len(*vs)-1]
		fmt.Printf("cur=%s\n", cur)

		if cur.Equal(ws, view) {
			return serial
		}

		prev, _ := vs.Find(ws, view); if prev == nil {
			if !vs.Append(NewViewState(serial, ws, view)) {
				return serial
			}
		} else {
			prev.Serial = serial
		}

		sort.Sort(vs)
		return serial + 1
	}
}

func (vs ViewStates) Prev() *ViewState {
	if len(vs) < 0 {
		return nil
	} else if len(vs) < 2 {
		return vs[len(vs)-1]
	} else {
		return vs[len(vs)-2]
	}
}
