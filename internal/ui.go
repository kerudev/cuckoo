package cuckoo

import (
	"fmt"
	"sort"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	Name    string
	Min     int // 0-59
	Hour    int // 0-23
	Day     int // 1-31
	Month   int // 1-12
	Weekday int // 0-6
}

type Coord struct {
	Name string
	X    float32
	Y    float32
}

type GridCoord struct {
	Names []string
	X     float32
	Y     float32
}

type Cell struct {
	W float32
	H float32
}

type Grid struct {
	W        float32
	H        float32
	Rows     int
	Cols     int
	HighestY int
}

type DrawMode int32

const (
	DrawNone DrawMode = iota
	DrawLines
	DrawBezier
)

type GroupBy int32

const (
	GroupByWdHour GroupBy = iota
	GroupByWdHourMin
	// GroupByMin
)

type StepMin int32

const (
	StepMin1 StepMin = iota
	StepMin5
	StepMin10
	StepMin15
	StepMin20
	StepMin30
)

// Constants
const INITIAL_ROWS = 10
const INITIAL_COLS = 24
const ROWS_CAP = 30

// Internal
var offset = rl.Vector2{X: 20, Y: 20}
var cell = Cell{W: 0, H: 0}
var grid = Grid{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}

var fontSize = float32(12)
var font = rl.Font{}

var boxRoundness = float32(0.2)
var boxSegments = int32(8)
var boxPadX = float32(16)

var colors = []rl.Color{
	rl.Red,
	rl.Green,
	rl.Blue,
	rl.Purple,
	rl.Beige,
	rl.Pink,
	rl.Orange,
}

// User options
var drawCoords = true
var drawMode = DrawLines
var stepMin = StepMin1
var groupBy = GroupByWdHourMin

var weekdaysToggle = []bool{
	true, // rl.Red
	true, // rl.Green
	true, // rl.Blue
	true, // rl.Purple
	true, // rl.Beige
	true, // rl.Pink
	true, // rl.Orange
}

func coordToVec2(coord GridCoord) rl.Vector2 {
	return rl.Vector2{X: coord.X, Y: coord.Y}
}

func coordToGrid(coords [][]Coord, grid *Grid) [][]GridCoord {
	result := make([][]GridCoord, 7)

	grid.Rows = INITIAL_ROWS
	grid.Cols = INITIAL_COLS

	for day, coordDay := range coords {
		for _, coord := range coordDay {
			found := false

			for i := range result[day] {
				if coord.X == result[day][i].X {
					found = true
					result[day][i].Names = append(result[day][i].Names, coord.Name)
				}

				if len(result[day][i].Names) >= grid.Rows {
					grid.Rows = len(result[day][i].Names) + 2
				}
			}

			if !found {
				result[day] = append(result[day], GridCoord{
					Names: []string{coord.Name},
					X:     coord.X,
					Y:     coord.Y,
				})
			}
		}
	}

	grid.HighestY = grid.Rows

	// Remove the last column, as it makes no sense when grouping by hour
	if groupBy == GroupByWdHour {
		grid.Cols -= 1
	}

	if grid.HighestY > ROWS_CAP {
		grid.Rows = INITIAL_ROWS
	}

	cell.W = grid.W / float32(grid.Cols)
	cell.H = grid.H / float32(grid.Rows)

	for day := range 7 {
		for i := range result[day] {
			result[day][i].X = result[day][i].X/float32(grid.Cols)*grid.W + offset.X
			result[day][i].Y = grid.H + offset.Y - (grid.H / float32(grid.HighestY) * float32(len(result[day][i].Names)))
		}
	}

	return result
}

func drawGrid(gridCoords [][]GridCoord) {
	// Draw lines vertically
	for col := range grid.Cols {
		x := cell.W*float32(col+1) + offset.X

		rl.DrawLineEx(
			rl.Vector2{X: x, Y: offset.Y},
			rl.Vector2{X: x, Y: grid.H + offset.Y},
			2,
			rl.LightGray,
		)
	}

	// Draw lines horizontally
	for row := range grid.Rows {
		y := cell.H*float32(row+1) + offset.Y

		rl.DrawLineEx(
			rl.Vector2{X: offset.X, Y: y},
			rl.Vector2{X: grid.W + offset.Y, Y: y},
			2,
			rl.LightGray,
		)
	}

	// Draw grid container
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: offset.X, Y: offset.Y, Width: grid.W, Height: grid.H},
		2,
		rl.Black,
	)

	// Draw text on X axis
	cols := grid.Cols
	if groupBy == GroupByWdHour {
		cols += 1
	}

	for col := range cols {
		text := strconv.Itoa(col)

		textW := rl.MeasureTextEx(font, text, fontSize, 1).X
		textX := cell.W*float32(col) - textW/2 + offset.X

		rl.DrawText(text, int32(textX), int32(grid.H+offset.Y+2), int32(fontSize), rl.Black)
	}

	// Draw text on Y axis
	textRect := rl.MeasureTextEx(font, strconv.Itoa(cols), fontSize, 1)

	nRow := 0
	isCapped := false
	last := (grid.HighestY / (grid.HighestY / grid.Rows)) * (grid.HighestY / grid.Rows)

	for row := range grid.HighestY + 1 {
		if grid.HighestY > ROWS_CAP && row%(grid.HighestY/grid.Rows) != 0 {
			isCapped = true
			continue
		}

		if isCapped && row == last {
			row = grid.HighestY
		}

		text := strconv.Itoa(row)
		textSize := rl.MeasureTextEx(font, strconv.Itoa(row), fontSize, 1)

		textPos := rl.Vector2{
			X: textRect.X + rl.Lerp(0.0, textRect.X-textSize.X, 1),
			Y: textRect.Y + rl.Lerp(0.0, textRect.Y-textSize.Y, 0.5),
		}

		textY := -cell.H*float32(nRow) - textPos.Y/2 + offset.Y + grid.H
		nRow++

		rl.DrawText(text, int32(textPos.X-offset.X/2), int32(textY), int32(fontSize), rl.Black)
	}

	for day, dayCoords := range gridCoords {
		if !weekdaysToggle[day] {
			continue
		}

		if drawMode != DrawNone {
			// Sort coordinates to draw line in order
			sort.Slice(dayCoords, func(i, j int) bool {
				return dayCoords[i].X < dayCoords[j].X
			})

			// Draw lines that connect coordinates
			for k := 0; k < len(dayCoords)-1; k++ {
				start := coordToVec2(dayCoords[k])
				end := coordToVec2(dayCoords[k+1])

				switch drawMode {
				case DrawLines:
					rl.DrawLineEx(start, end, 2, colors[day])
				case DrawBezier:
					rl.DrawLineBezier(start, end, 2, colors[day])
				}
			}
		}

		// Draw coordinates
		if drawCoords {
			for _, coord := range dayCoords {
				rl.DrawCircle(int32(coord.X), int32(coord.Y), 4, colors[day])
			}
		}
	}
}

func drawOptions(groupByScroll *int32) {
	// Draw option - DrawMode
	rl.DrawText("Draw mode", int32(offset.X), int32(grid.H+offset.Y*2+6), 12, rl.Black)
	drawModeIdx := int32(drawMode)
	drawModeIdx = rg.ToggleGroup(
		rl.Rectangle{X: offset.X, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
		"#113#;#127#;#125#",
		drawModeIdx,
	)
	drawMode = DrawMode(drawModeIdx)

	// Draw option - GroupBy
	rl.DrawText("Group by", int32(offset.X), int32(grid.H+offset.Y*4+6), 12, rl.Black)
	groupByIdx := int32(groupBy)
	groupByIdx = rg.ListView(
		rl.Rectangle{X: offset.X, Y: grid.H + offset.Y*5, Width: 100, Height: 31*2 + 1},
		"Wd+Hour;Wd+Hour+Min",
		groupByScroll,
		groupByIdx,
	)
	groupBy = GroupBy(groupByIdx)

	if groupBy == GroupByWdHourMin {
		// Check the implementation of GuiLoadStyleDefault for additional keys
		// https://github.com/raysan5/raygui/blob/master/src/raygui.h

		rl.DrawText("Weekdays", int32(120+offset.X), int32(grid.H+offset.Y*4+6), 12, rl.Black)

		def_BORDER_WIDTH := rg.GetStyle(rg.BUTTON, rg.BORDER_WIDTH)

		def_BORDER_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED)
		def_BASE_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED)
		def_BORDER_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED)
		def_BASE_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED)

		rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, 1)

		for i := range weekdaysToggle {
			color := colors[i]
			hexColor := rg.NewColorPropertyValue(color)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hexColor)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, LerpColorToHex(color, 0.8))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hexColor)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, LerpColorToHex(color, 0.7))

			rec := rl.Rectangle{
				X:      120 + offset.X + float32(22*i),
				Y:      grid.H + offset.Y*5,
				Width:  20,
				Height: 20,
			}

			active := rg.Toggle(rec, strconv.Itoa(i), weekdaysToggle[i])
			weekdaysToggle[i] = active
		}

		// Reset style to defaults
		rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, def_BORDER_WIDTH)

		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, def_BASE_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_PRESSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, def_BASE_COLOR_PRESSED)
	}

	// Draw option - StepMin
	if groupBy == GroupByWdHourMin {
		rl.DrawText("Step of x minutes", int32(120+offset.X), int32(grid.H+offset.Y*6+6), 12, rl.Black)
		stepMinIdx := int32(stepMin)
		stepMinIdx = rg.ToggleGroup(
			rl.Rectangle{X: 120 + offset.X, Y: grid.H + offset.Y*7, Width: 20, Height: 20},
			"#113#;5;10;15;20;30",
			stepMinIdx,
		)
		stepMin = StepMin(stepMinIdx)
	}

	// Draw option - DrawCoords
	drawCoords = rg.CheckBox(
		rl.Rectangle{X: 120 + offset.X, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
		"Draw coordinates",
		drawCoords,
	)

	if drawMode == DrawNone && !drawCoords {
		drawCoords = true
	}
}

func drawMouseOver(gridCoords [][]GridCoord) {
	// Draw coordinate information on mouse hover
	for _, coordDay := range gridCoords {
		for _, coord := range coordDay {
			mouseOverCoord := rl.CheckCollisionPointCircle(rl.GetMousePosition(), coordToVec2(coord), 4)

			if mouseOverCoord {
				maxW := float32(0)

				// Calculate max text size
				for _, name := range coord.Names {
					textW := float32(rl.MeasureText(name, int32(fontSize)))

					if textW > maxW {
						maxW = textW
					}
				}

				rec := rl.Rectangle{
					X:      coord.X + 8,
					Y:      coord.Y - 8,
					Width:  maxW + 16,
					Height: fontSize*float32(len(coord.Names)) + 16,
				}

				rl.DrawRectangleRounded(rec, boxRoundness, boxSegments, rl.White)
				rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, boxSegments, 2, rl.Black)

				sort.Slice(coord.Names, func(i, j int) bool {
					return sortAlphabetically(coord.Names[i], coord.Names[j])
				})

				for i, name := range coord.Names {
					spacingY := float32(i) * fontSize

					rl.DrawText(
						name,
						int32(coord.X+boxPadX),
						int32(coord.Y+spacingY),
						int32(fontSize),
						rl.Black,
					)
				}
			}
		}
	}
}

func DrawLoop(sample map[string]string) {
	crons := stringsToCrons(sample)
	coords := cronsToCoords(crons)

	gridCoords := [][]GridCoord{}

	groupByScroll := int32(0)

	// Previous state
	prevScreenW := int32(0)
	prevScreenH := int32(0)

	prevGroupBy := groupBy
	prevStepMin := stepMin

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)

	rl.InitWindow(800, 700, "Cuckoo")
	rl.SetWindowMinSize(800, 700)

	font = rl.GetFontDefault()

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		// Check if a file was dropped and reload coords

		// TODO this sometimes crashes when the file is being dragged over the
		// window. This may be a problem of Go's bindings
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			sample = map[string]string{}
			err := ReadPath(droppedFiles[0], &sample)

			if err != nil {
				fmt.Println(err)
			} else {
				crons = stringsToCrons(sample)
				coords = cronsToCoords(crons)
				gridCoords = coordToGrid(coords, &grid)
			}

			// rl.UnloadDroppedFiles();
		}

		// Recalculate grid and coordinates only when screen changes size
		if screenH != float32(prevScreenH) || screenW != float32(prevScreenW) {
			grid.W = screenW - 40
			grid.H = screenH - 240

			gridCoords = coordToGrid(coords, &grid)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawGrid(gridCoords)
		drawOptions(&groupByScroll)
		drawMouseOver(gridCoords)

		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if prevStepMin != stepMin {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Recalculate coordinates based on group by
		if prevGroupBy != groupBy {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		prevScreenW = int32(screenW)
		prevScreenH = int32(screenH)

		prevGroupBy = groupBy
		prevStepMin = stepMin
	}

	rl.CloseWindow()
}
