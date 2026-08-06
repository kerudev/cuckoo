package models

import rl "github.com/gen2brain/raylib-go/raylib"

// Constants
const INITIAL_ROWS = 10
const INITIAL_COLS = 24
const ROWS_CAP = 30
const ROWS_RATIO = ROWS_CAP / INITIAL_ROWS
const WEEKDAYS = 7

// UI (runtime)
var Font = rl.Font{}
var ListViewItemH = int32(0)

// UI (constants)
const FontSize = int32(12)
const FontRadius = FontSize / 2
const TextPad = int32(8)

const BoxRadius = float32(8.0)
const BoxDiameter = 2 * BoxRadius
const BoxSegments = int32(8)
const BoxSize = int32(20)
const BoxBorder = int32(1)
const BoxPad = BoxSize + BoxBorder*2

const HelpLineH = int32(24)
const HelpLineEmptyH = int32(14)

const CoordRadius = float32(4.0)

const GridBorder = int32(2)

var GridBGColor = rl.NewColor(200, 230, 250, 80)

const ZoomSliderH = int32(10)

const TooltipTimeFontSize = int32(16)
const TooltipScrollW = int32(10)

const FooterW = int32(120)
const FooterFontSize = int32(16)

var AllowedExt = []string{".json"}

// Internal
var Offset = Vector2Int32{X: 20, Y: 20}
var Cell = Rec[float32]{}

var Grid = rl.RectangleInt32{X: Offset.X, Y: Offset.Y * 3}
var Footer = rl.RectangleInt32{}
var Tooltip = rl.RectangleInt32{}
var HelpWindow = rl.RectangleInt32{}

var Sample = map[string]string{}
var Crons = []Cron{}
var Coords = [][]Coord{}
var GridCoords = [][]GridCoord{}

var WdCounts = [WEEKDAYS]CronCountsByWd{}
var MouseOver = [WEEKDAYS][]GridCoord{}
var TotalOver = 0
var TooltipHasOverflow = false

var ErrorText = ""
var ShowHelp = false

// Context
var C_Grid = GridContext{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}
var C_Zoom = ZoomContext{Factor: 1, Scale: 1}

// State
var S_Screen = NewState(Rec[int32]{})
var S_Mouse = NewState(rl.Vector2{})
var S_MouseWithLock = NewState(rl.Vector2{})
var S_IsMouseLocked = NewState(false)
var S_IsOverGrid = NewState(false)

var S_Zoom = NewState(float32(1.0))
var S_ZoomSlider = NewState(float32(0.0))

var S_TooltipScroll = NewState(int32(0))
var S_TooltipScrollMax = NewState(int32(0))

var S_FileScroll = NewState(int32(-1))
var S_FileActive = NewState(int32(-1))
var S_FileFocused = NewState(int32(-1))

var S_FileName = NewState("")
var S_LastFile = NewState("")
var S_FilePicker = NewState(false)

// User options
var S_Weekdays = NewState(Weekdays{
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

var UserOpt = UserOptions{
	DrawCoords: true,
	DrawLines:  true,
	DrawGrid:   true,
	DrawFade:   true,
}

var S_Position = NewState(PositionGrid)
