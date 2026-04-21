package cuckoo

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
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

var zoom = float32(0)
var zoomSliderVal = float32(0)

var colors = []rl.Color{
	rl.Red,
	rl.Orange,
	rl.Gold,
	rl.Green,
	rl.Blue,
	rl.Purple,
	rl.Pink,
}

// User options
var drawCoords = true
var drawMode = DrawLines
var stepMin = StepMin1
var groupBy = GroupByWdHourMin

var weekdaysToggle = []Status{
	StatusOn, // rl.Red
	StatusOn, // rl.Orange
	StatusOn, // rl.Gold
	StatusOn, // rl.Green
	StatusOn, // rl.Blue
	StatusOn, // rl.Purple
	StatusOn, // rl.Pink
}

// Previous state
var prevScreenW = int32(0)
var prevScreenH = int32(0)

var prevGroupBy = groupBy
var prevStepMin = stepMin

var prevZoom = zoom

func drawGrid(gridCoords [][]GridCoord) {
	cols := grid.Cols
	if groupBy == GroupByWdHour {
		cols += 1
	}

	// Draw lines vertically
	colX := offset.X
	for range grid.Cols {
		colX += cell.W
		rl.DrawLineEx(
			rl.Vector2{X: colX, Y: offset.Y},
			rl.Vector2{X: colX, Y: grid.H + offset.Y},
			2,
			rl.LightGray,
		)
	}

	// Draw lines horizontally
	rowY := offset.Y
	for range grid.Rows {
		rowY += cell.H
		rl.DrawLineEx(
			rl.Vector2{X: offset.X, Y: rowY},
			rl.Vector2{X: grid.W + offset.Y, Y: rowY},
			2,
			rl.LightGray,
		)
	}

	// Draw background line on mouse over
	mouse := rl.GetMousePosition()
	bg := rl.NewColor(200, 230, 250, 80)
	bgX := offset.X

	for range cols {
		mouseInX := bgX < mouse.X && mouse.X <= bgX+cell.W
		mouseInY := offset.Y < mouse.Y && mouse.Y <= grid.H

		if mouseInX && mouseInY {
			rec := rl.Rectangle{X: bgX + 1, Y: offset.Y, Width: cell.W - 2, Height: grid.H}
			rl.DrawRectangleRec(rec, bg)
		}

		bgX += cell.W

		if zoom == 0 {
			zoomSliderVal = rl.Clamp(mouse.X-cell.W, 0, grid.W-4)
		}
	}

	// Draw grid container
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: offset.X, Y: offset.Y, Width: grid.W, Height: grid.H},
		2,
		rl.Black,
	)

	// Draw zoom slider
	scroll := rl.GetMouseWheelMove()
	zoom = rl.Clamp(zoom+scroll, 0, 8)

	t := zoom / 8.0
	base := grid.W / float32(grid.Cols)

	factor := float32(math.Pow(float64(grid.W/base), float64(t)))
	cell.W = base * factor

	if zoom != 0 {
		zoomSliderVal = rg.Slider(
			rl.Rectangle{X: offset.X + 2, Y: grid.H + 6, Width: grid.W - 4, Height: 12},
			"",
			"",
			zoomSliderVal,
			0,
			grid.W-4,
		)
	}

	// Draw text on X axis
	for col := range cols {
		text := strconv.Itoa(col)

		textW := rl.MeasureTextEx(font, text, fontSize, 1).X
		textX := cell.W*float32(col) - textW/2 + offset.X

		rl.DrawText(text, int32(textX), int32(grid.H+offset.Y+2), int32(fontSize), rl.Black)
	}

	// Draw text on Y axis
	textRect := rl.MeasureTextEx(font, strconv.Itoa(cols), fontSize, 1)

	nRow := 0
	step := grid.HighestY / grid.Rows
	last := (grid.HighestY / step) * step

	for row := range grid.HighestY + 1 {
		if grid.HighestY > ROWS_CAP && row%step != 0 {
			continue
		}

		if row == last {
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
		if weekdaysToggle[day] != StatusOn {
			continue
		}

		if drawMode != DrawNone {
			// Sort coordinates to draw line in order
			sort.Slice(dayCoords, func(i, j int) bool {
				return dayCoords[i].X < dayCoords[j].X
			})

			// Draw lines that connect coordinates
			for k := 0; k < len(dayCoords)-1; k++ {
				start := dayCoords[k].Vec2()
				end := dayCoords[k+1].Vec2()

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
				rl.DrawCircle(int32(coord.X*(zoom+1)), int32(coord.Y), 4, colors[day])
			}
		}
	}
}

func drawOptions(groupByScroll *int32) {
	// Draw option - DrawMode
	rl.DrawText("Draw mode", int32(offset.X), int32(grid.H+offset.Y*2+8), 12, rl.Black)
	drawModeIdx := int32(drawMode)
	drawModeIdx = rg.ToggleGroup(
		rl.Rectangle{X: offset.X, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
		"#113#;#127#;#125#",
		drawModeIdx,
	)
	drawMode = DrawMode(drawModeIdx)

	// Draw option - DrawCoords
	drawCoordsIcon := ""
	if drawCoords {
		drawCoordsIcon = "#212#"
	} else {
		drawCoordsIcon = "#213#"
	}

	drawCoords = rg.Toggle(
		rl.Rectangle{X: offset.X + 22*3, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
		drawCoordsIcon,
		drawCoords,
	)

	if drawMode == DrawNone && !drawCoords {
		drawCoords = true
	}

	// Draw option - GroupBy
	rl.DrawText("Group by", int32(offset.X), int32(grid.H+offset.Y*4+8), 12, rl.Black)
	groupByIdx := int32(groupBy)
	groupByIdx = rg.ListView(
		rl.Rectangle{X: offset.X, Y: grid.H + offset.Y*5, Width: 100, Height: 31*2 + 1},
		"Wd+Hour;Wd+Hour+Min",
		groupByScroll,
		groupByIdx,
	)

	if groupByIdx >= 0 {
		groupBy = GroupBy(groupByIdx)
	}

	// Draw option - Weekdays

	// Check the implementation of GuiLoadStyleDefault for additional keys
	// https://github.com/raysan5/raygui/blob/master/src/raygui.h

	rl.DrawText("Weekdays", int32(120+offset.X), int32(grid.H+offset.Y*4+8), 12, rl.Black)

	def_BORDER_WIDTH := rg.GetStyle(rg.BUTTON, rg.BORDER_WIDTH)

	def_BASE_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL)

	def_BORDER_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED)
	def_BASE_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED)
	def_BORDER_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED)
	def_BASE_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED)

	rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, 1)

	for i, status := range weekdaysToggle {
		color := colors[i]
		hexColor := rg.NewColorPropertyValue(color)

		if status == StatusDisabled {
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_FOCUSED)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, def_BASE_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_PRESSED)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, def_BASE_COLOR_NORMAL)
		} else {
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hexColor)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, lerpColorToHex(color, 0.8))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hexColor)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, lerpColorToHex(color, 0.7))
		}

		rec := rl.Rectangle{
			X:      120 + offset.X + float32(22*i),
			Y:      grid.H + offset.Y*5,
			Width:  20,
			Height: 20,
		}

		active := rg.Toggle(rec, strconv.Itoa(i), status.Bool())

		if status != StatusDisabled {
			weekdaysToggle[i] = StatusFromBool(active)

			if !active && all(weekdaysToggle, StatusOff) {
				weekdaysToggle[i] = StatusOn
			}
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
		rl.DrawText("Step of x minutes", int32(120+offset.X), int32(grid.H+offset.Y*6+8), 12, rl.Black)
		stepMinIdx := int32(stepMin)
		stepMinIdx = rg.ToggleGroup(
			rl.Rectangle{X: 120 + offset.X, Y: grid.H + offset.Y*7, Width: 20, Height: 20},
			"#113#;5;10;15;20;30",
			stepMinIdx,
		)
		stepMin = StepMin(stepMinIdx)
	}
}

func drawMouseOver(gridCoords [][]GridCoord) {
	// Draw coordinate information on mouse hover
	for _, coordDay := range gridCoords {
		for _, coord := range coordDay {
			mouseOverCoord := rl.CheckCollisionPointCircle(rl.GetMousePosition(), coord.Vec2(), 4)

			if mouseOverCoord {
				maxW := float32(0)

				// Calculate max text size
				names := countDuplicates(coord.Names)
				fmtNames := []string{}
				
				for name, count := range names {
					fmtNames = append(fmtNames, fmt.Sprintf("%s (%d)", name, count))
				}

				for _, name := range fmtNames {
					textW := float32(rl.MeasureText(name, int32(fontSize)))

					if textW > maxW {
						maxW = textW
					}
				}

				rec := rl.Rectangle{
					X:      coord.X + 8,
					Y:      coord.Y - 8,
					Width:  maxW + 16,
					Height: fontSize*float32(len(fmtNames)) + 16,
				}

				rl.DrawRectangleRounded(rec, boxRoundness, boxSegments, rl.White)
				rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, boxSegments, 2, rl.Black)

				sort.Slice(fmtNames, func(i, j int) bool {
					return sortAlphabetically(fmtNames[i], fmtNames[j])
				})

				for i, name := range fmtNames {
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

func drawFooter() {
	screenW := int32(rl.GetScreenWidth())

	footerW := int32(120)
	footerH := int32(100)
	footerX := int32(screenW-int32(offset.X)) - footerW
	footerY := int32(grid.H + offset.Y*2 + 8)

	textPad := int32(8)

	rl.DrawRectangleLines(footerX, footerY, footerW, footerH, rl.Black)
	// rl.DrawText("Zoom  : "+zoom.String(), footerX+textPad, footerY+textPad, 16, rl.Black)
	rl.DrawText(fmt.Sprint("Cell.W: ", cell.W), footerX+textPad, footerY+textPad*2+16, 16, rl.Black)
	rl.DrawText(fmt.Sprint("Cell.H: ", cell.H), footerX+textPad, footerY+textPad*4+16, 16, rl.Black)
}

func DrawLoop(sample map[string]string) {
	crons := stringsToCrons(sample)
	jobs := cronsToJobs(crons)
	coords := jobsToCoords(jobs)

	gridCoords := [][]GridCoord{}

	groupByScroll := int32(0)

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
				jobs = cronsToJobs(crons)
				coords = cronsToCoords(crons)
				gridCoords = coordToGrid(coords, &grid)
			}
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
		drawFooter()
		drawMouseOver(gridCoords)

		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if stepMin != prevStepMin {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Recalculate coordinates based on group by
		if groupBy != prevGroupBy {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		prevScreenW = int32(screenW)
		prevScreenH = int32(screenH)

		prevGroupBy = groupBy
		prevStepMin = stepMin

		prevZoom = zoom
	}

	rl.CloseWindow()
}
