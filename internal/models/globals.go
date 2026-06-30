package models

import rl "github.com/gen2brain/raylib-go/raylib"

// Constants
const INITIAL_ROWS = 10
const INITIAL_COLS = 24
const ROWS_CAP = 30
const WD_COUNT = 7

// UI
var Font = rl.Font{}
var FontSize = int32(12)
var FontRadius = FontSize / 2
var TextPad = int32(8)

var BoxRadius = float32(8.0)
var BoxDiameter = 2 * BoxRadius
var BoxSegments = int32(8)
var BoxSize = int32(20)
var BoxBorder = int32(1)
var BoxPad = BoxSize + BoxBorder*2
var CoordRadius = float32(4.0)

var GridBorder = int32(2)

var TooltipTimeFontSize = int32(16)
var TooltipScrollW = int32(10)

var FooterW = int32(120)
var FooterFontSize = int32(16)

// Internal
var Offset = Vector2Int32{X: 20, Y: 20}
var Cell = Rec[float32]{}
var Grid = rl.RectangleInt32{X: Offset.X, Y: Offset.Y}
var Tooltip = rl.RectangleInt32{}

var MouseOver = make([][]GridCoord, WD_COUNT)
var TotalOver = 0

// Context
var C_Grid = GridContext{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}
var C_Zoom = ZoomContext{Factor: 1, Scale: 1}

// State
var S_Screen = NewState(Rec[int32]{})
var S_Mouse = NewState(rl.Vector2{})
var S_IsMouseLocked = NewState(false)

var S_Weekdays = NewState([WD_COUNT]Weekday{
	NewWeekday(rl.Red),
	NewWeekday(rl.Orange),
	NewWeekday(rl.Gold),
	NewWeekday(rl.Green),
	NewWeekday(rl.Blue),
	NewWeekday(rl.Purple),
	NewWeekday(rl.Pink),
})

var S_GroupBy = NewState(GroupByWdHourMin)
var S_StepMin = NewState(StepMin1)

var S_Zoom = NewState(float32(1.0))
var S_ZoomSlider = NewState(float32(0.0))

var S_TooltipScroll = NewState(int32(0))
var S_TooltipScrollMax = NewState(int32(0))
var S_TooltipHasOverflow = NewState(false)

// User options
var UserOpt = UserOptions{
	DrawCoords: true,
	DrawLines:  true,
	DrawGrid:   true,
	DrawFade:   true,
}

var Position = PositionGrid
