package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"cbgo"
)

func evDecodeErr(d []byte, err error) {
	fmt.Printf("couldn't decode event data '%s': %s\n", d, err)
}

type State struct {
	ipc	cbgo.CbIpc
	dump struct {
		initial		bool
		requested	bool
		currentWS	bool
	}
	Serial	int
	Access	cbgo.ViewStates
}

func (s *State) pushView(ws, view int) {
	fmt.Printf("pushView(%d, %d) in serial=%d access=%s\n", ws, view, s.Serial, s.Access)
	s.Serial = s.Access.Push(s.Serial, ws, view)
	fmt.Printf("pushView out serial=%d access=%s\n", s.Serial, s.Access)
}

func (s *State) changeViewWS(oldWs, view, newWs int) {
	s.Access.Pop(oldWs, view)
	s.pushView(newWs, view)
}

func (s *State) initViews(dump *cbgo.DumpEvent) {
	cur_output, ok := dump.Outputs[dump.CurOutput]
	if !ok {
		fmt.Println("invalid curr_output id:", dump.CurOutput)
		return
	}

	for n, ws := range(cur_output.Workspaces) {
		ws_id := n + 1
		for _, v := range(ws.Views) {
			if v.ID == 0 {
				continue
			}
			if cur_output.CurWS == ws_id && dump.CurView == v.ID {
				continue
			}
			s.pushView(ws_id, v.ID)
		}
	}

	s.pushView(cur_output.CurWS, dump.CurView)
}

func (s *State) listViews(dump *cbgo.DumpEvent, currentWS bool) {
	cur_output, ok := dump.Outputs[dump.CurOutput]
	if !ok {
		fmt.Println("invalid curr_output id:", dump.CurOutput)
		return
	}

	var workspaces []*cbgo.Workspace
	if currentWS && cur_output.CurWS > 0 && cur_output.CurWS <= len(cur_output.Workspaces) {
		workspaces = []*cbgo.Workspace{cur_output.Workspaces[cur_output.CurWS-1]}
	} else {
		workspaces = cur_output.Workspaces
	}

	view_names := []string{}

	for n, ws := range(workspaces) {
		sort.Sort(ws.Views)
		for _, v := range(ws.Views) {
			if v.ID == 0 {
				continue
			}

			var title string

			if v.Title == "" {
				p := filepath.Join("/proc", strconv.Itoa(v.PID), "comm")
				c, err := os.ReadFile(p)

				if err != nil {
					title = "unknown"
				} else {
					title = strings.TrimSpace(string(c))
				}
			} else {
				title = v.Title
			}

			var sep string
			if v.ID == dump.CurView {
				sep = "+"
			} else {
				sep = "-"
			}

			vn := strconv.Itoa(v.ID) + sep + title
			if ! currentWS {
				vn = strconv.Itoa(n+1) + ":" + vn
			}
			view_names = append(view_names, vn)
		}
	}
	m := "message " + strings.Join(view_names, "\r")
	s.ipc.SendCmd(m)
}

func (s *State) handleDump(d []byte) {
	if !s.dump.requested {
		return
	}

	var ev cbgo.DumpEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		if s.dump.initial {
			s.initViews(&ev)
			s.dump.initial = false
		} else {
			s.listViews(&ev, s.dump.currentWS)
			s.dump.currentWS = false
		}
	}
	s.dump.requested = false
}

func (s *State) handleCycleViews(d []byte) {
	var ev cbgo.CycleViewsEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		fmt.Printf("%s: %#v\n", ev.Name, ev)
		s.pushView(ev.Workspace, ev.NewID)
	}
}

func (s *State) handleSwitchWS(d []byte) {
	var ev cbgo.SwitchWSEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		fmt.Printf("%s: %#v\n", ev.Name, ev)
		s.pushView(ev.NewWorkspace, ev.FocusedView)
	}
}

func (s *State) handleMap(d []byte) {
	var ev cbgo.ViewMapUnmapEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		fmt.Printf("%s: %#v\n", ev.Name, ev)
		s.pushView(ev.Workspace, ev.View)
	}
}

func (s *State) handleUnmap(d []byte) {
	var ev cbgo.ViewMapUnmapEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		fmt.Printf("%s: %#v\n", ev.Name, ev)
		s.Access.Pop(ev.Workspace, ev.View)
	}
}

func (s *State) handleMoveViewToWS(d []byte) {
	var ev cbgo.MoveViewToWSEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
	} else {
		fmt.Printf("%s: %#v\n", ev.Name, ev)
		s.changeViewWS(ev.OldWorkspace, ev.View, ev.NewWorkspace)
	}
}

func (s *State) handleCustomEvent(d []byte) {
	var ev cbgo.CustomEvent
	if err := json.Unmarshal(d, &ev); err != nil {
		evDecodeErr(d, err)
		return
	}

	fmt.Printf("customEvent %v\n", ev)

	cmdMsg := strings.ReplaceAll(ev.Message, `\"`, `"`)
	var cmd cbgo.CustomCmd
	err := json.Unmarshal([]byte(cmdMsg), &cmd)
	if err != nil {
		fmt.Printf("couldn't decode event data '%s': %s\n", string(ev.Message), err)
		return
	}
	switch cmd.Command {
	case "other_view":
		prev := s.Access.Prev()
		fmt.Printf("other: access=%s, prev=%s\n", s.Access, prev)

		if prev == nil {
			s.ipc.SendCmd("message no other view")
		} else {
			s.ipc.SendCmd(prev.SwitchCmd())
		}
	case "list_ws_views":
		s.dump.currentWS = true
		fallthrough
	case "list_views":
		s.dump.requested = true
		s.ipc.SendCmd("dump")
	}
}

func (s *State) Run() {
	s.dump.requested = true
	s.dump.initial = true
	s.ipc.SendCmd("dump")

	for {
		d, err := s.ipc.ReadEvent()
		if err != nil {
			fmt.Println("can't read event from IPC socket", err)
			os.Exit(-1)
		}

		var ev cbgo.Event
		err = json.Unmarshal(d, &ev)
		if err != nil {
			evDecodeErr(d, err)
			continue
		}

		switch ev.Name {
		case "dump":
			s.handleDump(d)

		case "custom_event":
			s.handleCustomEvent(d)

		case "cycle_views":
			s.handleCycleViews(d)

		case "move_view_to_ws":
			s.handleMoveViewToWS(d)

		case "switch_ws":
			s.handleSwitchWS(d)

		case "view_map":
			s.handleMap(d)

		case "view_unmap":
			s.handleUnmap(d)

		default:
			fmt.Printf("unhandled event %s\n", ev.Name)
		}
	}
}

func main() {
	var state State

	if err := state.ipc.Open(); err != nil {
		fmt.Println("can't open IPC socket:", err)
		os.Exit(-1)
	}
	defer state.ipc.Close()

	state.Run()
}
