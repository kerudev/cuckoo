package models

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants
const VERSION = "1.0.0"

const INITIAL_ROWS = 10
const INITIAL_COLS = 24
const ROWS_CAP = 30
const ROWS_RATIO = ROWS_CAP / INITIAL_ROWS
const WEEKDAYS = 7

const MIN_ZOOM = 1
const MAX_ZOOM = 9

const MIN_FILES = 0
const MAX_FILES = 6

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

const ModalLineH = int32(24)
const ModalLineEmptyH = int32(12)

const CoordRadius = float32(4.0)
const CoordDiameter = CoordRadius * 2

const GridBorder = int32(2)

var GridBGColor = rl.NewColor(200, 230, 250, 80)

const ZoomSliderH = int32(10)

const TooltipTimeFontSize = int32(16)
const TooltipScrollW = int32(10)

const FooterW = int32(120)
const FooterFontSize = int32(16)

// UI (components)
var Offset = Vector2Int32{X: 30, Y: 20}
var Cell = Rec[float32]{}
var Style = map[string]rg.PropertyValue{}

var BackButton = rl.Rectangle{
	X:      float32(Offset.X),
	Y:      float32(Offset.Y),
	Width:  float32(BoxSize),
	Height: float32(BoxSize),
}
var FileButton = rl.Rectangle{
	X:      BackButton.X + float32(BoxPad),
	Y:      BackButton.Y,
	Width:  0, // Depends on S_Screen.Val.W
	Height: BackButton.Height,
}
var FileButtonText = rl.Rectangle{
	X:      FileButton.X + float32(BoxPad)/2,
	Y:      FileButton.Y,
	Width:  0, // Depends on S_Screen.Val.W
	Height: FileButton.Height,
}

var LockButton = rl.Rectangle{
	X:      0, // Depends on S_Screen.Val.W
	Y:      float32(Offset.Y),
	Width:  float32(BoxSize),
	Height: float32(BoxSize),
}
var HelpButton = rl.Rectangle{
	X:      0, // Depends on S_Screen.Val.W
	Y:      float32(Offset.Y),
	Width:  float32(BoxSize),
	Height: float32(BoxSize),
}
var AboutButton = rl.Rectangle{
	X:      0, // Depends on S_Screen.Val.W
	Y:      float32(Offset.Y),
	Width:  float32(BoxSize),
	Height: float32(BoxSize),
}

var CloseErrorButton = rl.Rectangle{
	X:      BackButton.X,
	Y:      BackButton.Y + float32(Offset.Y),
	Width:  BackButton.Width,
	Height: BackButton.Height,
}
var ErrorBox = rl.Rectangle{
	X:      FileButton.X,
	Y:      float32(Offset.Y + BoxSize),
	Width:  0, // Depends on S_Screen.Val.W
	Height: float32(BoxSize),
}
var ErrorMessageText = rl.Rectangle{
	X:      ErrorBox.X + float32(BoxPad)/2,
	Y:      ErrorBox.Y,
	Width:  0, // Depends on S_Screen.Val.W
	Height: ErrorBox.Height,
}

var Grid = rl.RectangleInt32{X: Offset.X, Y: Offset.Y * 3}
var Footer = rl.RectangleInt32{}
var Tooltip = rl.RectangleInt32{}
var HelpWindow = rl.RectangleInt32{}
var AboutWindow = rl.RectangleInt32{}

// Internal data
var Crons = []Cron{}
var Coords = [][]Coord{}
var GridCoords = [][]GridCoord{}
var TooltipLines = map[string]map[string]JobsCountsByWd{}

var DirFiles = []string{}
var DirFilesCount = int32(0)

var WdCounts = [WEEKDAYS]CronCountsByWd{}
var MouseOver = [WEEKDAYS][]GridCoord{}
var TotalOver = 0
var TooltipHasOverflow = false

var ErrorText = ""
var ShowHelp = false
var ShowAbout = false
var BlockUI = false

// Context
var C_Grid = GridContext{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}
var C_Zoom = ZoomContext{Scale: 1}

// State
var S_Screen = NewState(Rec[int32]{})
var S_Mouse = NewState(rl.Vector2{})
var S_MouseWithLock = NewState(rl.Vector2{})
var S_IsMouseLocked = NewState(false)
var S_IsOnWindow = NewState(false)
var S_IsOverGrid = NewState(false)

var S_Zoom = NewState(float32(1.0))
var S_ZoomSlider = NewState(float32(0.0))

var S_TooltipScroll = NewState(int32(0))
var S_TooltipScrollMax = NewState(int32(0))

var S_FileScroll = NewState(int32(-1))
var S_FileActive = NewState(int32(-1))
var S_FileFocused = NewState(int32(-1))

var S_FilePath = NewState("")
var S_FileName = NewState("")
var S_FilePicker = NewState(false)

var S_DirLastUpdate = NewState(int64(0))
var S_FileLastUpdate = NewState(int64(0))

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
