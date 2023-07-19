package cbgo

import (
	_ "encoding/json"
)

type Coords struct {
	X	float64	`json:"x"`
	Y	float64	`json:"y"`
}

type Size struct {
	Width	int	`json:"width"`
	Height	int	`json:"height"`
}

type View struct {
	ID	int	`json:"id"`
	PID	int	`json:"pid"`
	Coords	*Coords	`json:"coords"`
	Type	string	`json:"type"`
	Title	string	`json:"title,omitempty"`
}

type Views []View

func (v Views) Len() int {
	return len(v)
}

func (v Views) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Views) Less(i, j int) bool {
	return v[i].ID < v[j].ID
}

type Tile struct {
	ID	int	`json:"id"`
	Coords	*Coords	`json:"coords"`
	Size	*Size	`json:"size"`
	View	int	`json:"view_id"`
}

type Workspace struct {
	Views		Views		`json:"views"`
	Tiles		[]*Tile		`json:"tiles"`
}

type Output struct {
	Priority	int		`json:"priority"`
	Coords		*Coords		`json:"coords"`
	Size		*Size		`json:"size"`
	RefreshRate	float64		`json:"refresh_rate"`
	CurWS		int		`json:"curr_workspace"`
	Workspaces	[]*Workspace	`json:"workspaces"`
}

type Keyboard struct {
	CommandsEnabled	int		`json:"commands_enabled"`
	RepeatDelay	int		`json:"repeat_delay"`
	RepeatRate	int		`json:"repeat_rate"`
}

type InputDevice struct {
	IsVirtual	int		`json:"is_virtual"`
	Type		string		`json:"type"`
}

type DumpEvent struct {
	EventName	string			`json:"event_name"`
	NumWorkspaces	int			`json:"nws"`
	BgColor		[]float64		`json:"bg_color"`
	CurView		int			`json:"views_curr_id"`
	CurTile		int			`json:"tiles_curr_id"`
	CurOutput	string			`json:"curr_output"`
	Defaultmode	string			`json:"default_mode"`
	Modes		[]string		`json:"modes"`
	Outputs		map[string]*Output	`json:"outputs"`
	Keyboards	map[string]*Keyboard	`json:"keyboards"`
	InputDevices	map[string]*InputDevice	`json:"input_devices"`
	CursorCoords	*Coords		`json:"cursor_coords"`
}

type Event struct {
	Name		string	`json:"event_name"`
}

type CycleViewsEvent struct {
	Name		string	`json:"event_name"`
	OldID		int	`json:"old_view_id"`
	OldPID		int	`json:"old_view_pid"`
	NewID		int	`json:"new_view_id"`
	NewPID		int	`json:"new_view_pid"`
	TileID		int	`json:"tile_id"`
	Workspace	int	`json:"workspace"`
	Output		string	`json:"output"`
	OutputID	int	`json:"output_id"`
}

type SwitchWSEvent struct {
	Name		string	`json:"event_name"`
	OldWorkspace	int	`json:"old_workspace"`
	NewWorkspace	int	`json:"new_workspace"`
	Output		string	`json:"output"`
	OutputID	int	`json:"output_id"`
	FocusedTile	int	`json:"focused_tile_id"`
	FocusedView	int	`json:"focused_view_id"`
}

type CursorSwitchTileEvent struct {
	Name		string	`json:"event_name"`
	OldOutput	string	`json:"old_output"`
	OldOutputID	int	`json:"old_output_id"`
	OldTile		int	`json:"old_tile"`
	NewOutput	string	`json:"new_output"`
	NewOutputID	int	`json:"new_output_id"`
	NewTile		int	`json:"new_tile"`
}

type ViewMapUnmapEvent struct {
	Name		string	`json:"event_name"`
	View		int	`json:"view_id"`
	Tile		int	`json:"tile_id"`
	Workspace	int	`json:"workspace"`
	Output		string	`json:"output"`
	OutputID	int	`json:"output_id"`
	ViewPID		int	`json:"view_pid"`
}

type MoveViewToWSEvent struct {
	Name		string	`json:"event_name"`
	View		int	`json:"view_id"`
	OldWorkspace	int	`json:"old_workspace"`
	NewWorkspace	int	`json:"new_workspace"`
	Output		string	`json:"output"`
	OutputID	int	`json:"output_id"`
	ViewPID		int	`json:"view_pid"`
}

type CustomEvent struct {
	Name		string		`json:"event_name"`
	Message		string		`json:"message"`
}

type CustomCmd struct {
	Command		string		`json:"command"`
}
