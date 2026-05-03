package cuckoo

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants
const INITIAL_ROWS = 10
const INITIAL_COLS = 24
const ROWS_CAP = 30

// UI
var font = rl.Font{}
var fontSize = int32(12)
var textPad = int32(8)

var boxRoundness = float32(0.2)
var boxSegments = int32(8)
var boxSize = int32(20)
var boxBorder = int32(1)
var boxPadW = boxSize + boxBorder*2

var gridBorder = int32(2)

var footerW = int32(120)
var footerFontSize = int32(16)

var colors = []rl.Color{
	rl.Red,
	rl.Orange,
	rl.Gold,
	rl.Green,
	rl.Blue,
	rl.Purple,
	rl.Pink,
}

// Internal
var screen = Screen{W: 0, H: 0}
var offset = Vector2Int32{X: 20, Y: 20}
var cell = Cell{W: 0, H: 0}
var grid = Grid{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}

var zoom = float32(1.0)
var zoomSlider = float32(0.0)
var zoomOffset = float32(0.0)
var scale = float32(1.0)

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
var prevScreen = screen

var prevGroupBy = groupBy
var prevStepMin = stepMin

var prevZoom = zoom
var prevZoomSlider = zoomSlider

func drawGrid(gridCoords [][]GridCoord) {
	// Set all values that depend on the previous frame
	cols := grid.Cols
	if groupBy == GroupByWdHour {
		cols += 1
	}

	scroll := rl.GetMouseWheelMove()
	zoom = rl.Clamp(zoom+scroll, 1, 9)

	factor := float64(zoom-1) / 8.0
	base := float32(grid.W) / float32(grid.Cols)

	scale = float32(math.Pow(float64(grid.W)/float64(base), factor))
	zoomOffset = zoomSlider * (scale - 1)

	cell.W = base * scale

	// Draw lines vertically
	colX := float32(offset.X) - zoomOffset

	for range grid.Cols {
		colX += cell.W

		if colX < float32(offset.X) {
			continue
		}

		if colX > float32(grid.W+offset.X) {
			break
		}

		rl.DrawLineEx(
			rl.Vector2{X: colX, Y: float32(offset.Y)},
			rl.Vector2{X: colX, Y: float32(grid.H + offset.Y)},
			float32(gridBorder),
			rl.LightGray,
		)
	}

	// Draw lines horizontally
	rowY := float32(offset.Y)

	for range grid.Rows {
		rowY += cell.H
		rl.DrawLineEx(
			rl.Vector2{X: float32(offset.X), Y: rowY},
			rl.Vector2{X: float32(grid.W + offset.Y), Y: rowY},
			float32(gridBorder),
			rl.LightGray,
		)
	}

	// Draw background line on mouse over
	mouse := rl.GetMousePosition()
	bg := rl.NewColor(200, 230, 250, 80)
	bgX := float32(offset.X) - zoomOffset

	for range cols {
		mouseInX := bgX < mouse.X && mouse.X <= bgX+cell.W
		mouseInY := float32(offset.Y) < mouse.Y && mouse.Y <= float32(grid.H)

		if mouseInX && mouseInY {
			bgRec := rl.RectangleInt32{X: int32(bgX) + boxBorder, Y: offset.Y, Width: int32(cell.W) - boxBorder*2, Height: grid.H}
			rl.DrawRectangleRec(bgRec.ToFloat32(), bg)
		}

		bgX += cell.W

		if zoom == 1 {
			zoomSlider = rl.Clamp(mouse.X-cell.W, 0, float32(grid.W))
		}
	}

	// Draw zoom slider
	if zoom > 1 {
		scrollW := grid.W - gridBorder*2
		rg.SetStyle(rg.SLIDER, rg.SLIDER_WIDTH, rg.PropertyValue(float32(scrollW)/scale))

		zoomSliderRec := rl.RectangleInt32{X: offset.X + gridBorder, Y: grid.H + textPad - gridBorder, Width: scrollW, Height: 12}
		zoomSlider = rg.Slider(zoomSliderRec.ToFloat32(), "", "", zoomSlider, 0, float32(grid.W))
	}

	// Draw coordinates in layers by weekday
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
				start := dayCoords[k].Vector2()
				end := dayCoords[k+1].Vector2()

				switch drawMode {
				case DrawLines:
					rl.DrawLineEx(start, end, float32(gridBorder), colors[day])
				case DrawBezier:
					rl.DrawLineBezier(start, end, float32(gridBorder), colors[day])
				}
			}
		}

		// Draw coordinates
		if !drawCoords {
			continue
		}

		for _, coord := range dayCoords {
			if coord.X < float32(offset.X) {
				continue
			}

			if coord.X > float32(grid.W+offset.X) {
				break
			}

			rl.DrawCircle(int32(coord.X), int32(coord.Y), 4, colors[day])
		}
	}

	// Draw rectangles on left and right so lines are hidden
	// TODO optimize so this is not needed
	rl.DrawRectangle(0, offset.Y, offset.X, grid.H, rl.RayWhite)
	rl.DrawRectangle(grid.W+offset.X, offset.Y, offset.X, grid.H, rl.RayWhite)

	// Draw text on X axis
	for col := range cols {
		text := strconv.Itoa(col)

		textW := rl.MeasureTextEx(font, text, float32(fontSize), 1).X
		textX := cell.W*float32(col) - textW/2 + float32(offset.X) - zoomOffset

		// Clamp number to the left side
		if textX < float32(offset.X) {
			if textX+cell.W > float32(offset.X+textPad) {
				textX = float32(offset.X)
			} else {
				continue
			}
		}

		// Clamp number to the right side
		if textX > float32(grid.W+offset.X)-textW/2 {
			if textX-cell.W < float32(grid.W+offset.X-textPad*3) {
				textX = float32(grid.W+offset.X) - textW/2
			} else {
				continue
			}
		}

		rl.DrawText(text, int32(textX), grid.H+offset.Y+2, fontSize, rl.Black)
	}

	// Draw text on Y axis
	textRect := rl.MeasureTextEx(font, strconv.Itoa(cols), float32(fontSize), 1)

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
		textSize := rl.MeasureTextEx(font, strconv.Itoa(row), float32(fontSize), 1)

		textPos := rl.Vector2{
			X: textRect.X + rl.Lerp(0.0, textRect.X-textSize.X, 1),
			Y: textRect.Y + rl.Lerp(0.0, textRect.Y-textSize.Y, 0.5),
		}

		textY := float32(grid.H+offset.Y) - cell.H*float32(nRow) - textPos.Y/2
		nRow++

		rl.DrawText(text, int32(textPos.X-float32(offset.X)/2), int32(textY), fontSize, rl.Black)
	}

	// Draw grid container
	gridRec := rl.RectangleInt32{X: offset.X, Y: offset.Y, Width: grid.W, Height: grid.H}
	rl.DrawRectangleLinesEx(gridRec.ToFloat32(), 2, rl.Black)
}

func drawOptions(groupByScroll *int32) {
	// Draw option - DrawMode
	rl.DrawText("Draw mode", offset.X, grid.H+offset.Y*2+textPad, fontSize, rl.Black)
	drawModeRec := rl.RectangleInt32{X: offset.X, Y: grid.H + offset.Y*3, Width: boxSize, Height: boxSize}
	drawModeIdx := rg.ToggleGroup(drawModeRec.ToFloat32(), "#113#;#127#;#125#", int32(drawMode))
	drawMode = DrawMode(drawModeIdx)

	// Draw option - DrawCoords
	drawCoordsIcon := ""
	if drawCoords {
		drawCoordsIcon = "#212#"
	} else {
		drawCoordsIcon = "#213#"
	}

	drawCoordsRec := rl.RectangleInt32{X: offset.X + boxPadW*3, Y: grid.H + offset.Y*3, Width: boxSize, Height: boxSize}
	drawCoords = rg.Toggle(drawCoordsRec.ToFloat32(), drawCoordsIcon, drawCoords)

	if drawMode == DrawNone && !drawCoords {
		drawCoords = true
	}

	// Draw option - GroupBy
	rl.DrawText("Group by", offset.X, grid.H+offset.Y*4+textPad, fontSize, rl.Black)
	groupByRec := rl.RectangleInt32{X: offset.X, Y: grid.H + offset.Y*5, Width: 100, Height: 31*2 + 1}
	groupByIdx := rg.ListView(groupByRec.ToFloat32(), "Wd+Hour;Wd+Hour+Min", groupByScroll, int32(groupBy))

	if groupByIdx >= 0 {
		groupBy = GroupBy(groupByIdx)
	}

	// Draw option - Weekdays

	// Check the implementation of GuiLoadStyleDefault for additional keys
	// https://github.com/raysan5/raygui/blob/master/src/raygui.h

	rl.DrawText("Weekdays", 120+offset.X, grid.H+offset.Y*4+textPad, fontSize, rl.Black)

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

		button := rl.RectangleInt32{
			X:      120 + offset.X + boxPadW*int32(i),
			Y:      grid.H + offset.Y*5,
			Width:  boxSize,
			Height: boxSize,
		}

		active := rg.Toggle(button.ToFloat32(), strconv.Itoa(i), status.Bool())

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
		rl.DrawText("Step of x minutes", 120+offset.X, grid.H+offset.Y*6+textPad, fontSize, rl.Black)
		stepMinRec := rl.RectangleInt32{X: 120 + offset.X, Y: grid.H + offset.Y*7, Width: boxSize, Height: boxSize}
		stepMinIdx := rg.ToggleGroup(stepMinRec.ToFloat32(), "#113#;5;10;15;20;30", int32(stepMin))
		stepMin = StepMin(stepMinIdx)
	}
}

func drawMouseOver(gridCoords [][]GridCoord) {
	mouseOver := []GridCoord{}
	mouse := rl.GetMousePosition()

	// Get coords where mouse is over
	for day, dayCoords := range gridCoords {
		for _, coord := range dayCoords {
			if weekdaysToggle[day] != StatusOn {
				continue
			}

			if rl.CheckCollisionPointCircle(mouse, coord.Vector2(), 4) {
				mouseOver = append(mouseOver, coord)
			}
		}
	}

	if len(mouseOver) == 0 {
		return
	}

	names := []string{}
	for _, coord := range mouseOver {
		for _, job := range coord.Jobs {
			names = append(names, job.Name)
		}
	}

	duplicates := countDuplicates(names)

	finalNames := []string{}
	maxW := int32(0)

	for name, count := range duplicates {
		s := fmt.Sprintf("%s (%d)", name, count)
		finalNames = append(finalNames, s)

		if w := rl.MeasureText(s, fontSize); w > maxW {
			maxW = w
		}
	}

	sort.Slice(finalNames, func(i, j int) bool {
		return sortAlphabetically(finalNames[i], finalNames[j])
	})

	// Prepare tooltip
	base := mouseOver[0]

	tooltip := rl.RectangleInt32{
		X:      int32(base.X) + textPad,
		Y:      int32(base.Y) - textPad,
		Width:  maxW + textPad*2,
		Height: fontSize*int32(len(finalNames)) + textPad*2,
	}

	// Move tooltip to the left when it renders out of the grid
	if tooltip.X+tooltip.Width > offset.X+grid.W {
		tooltip.X = int32(base.X) - textPad - tooltip.Width
	}

	tooltipRec := tooltip.ToFloat32()
	rl.DrawRectangleRounded(tooltipRec, boxRoundness, boxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(tooltipRec, boxRoundness, boxSegments, 2, rl.Black)

	// Draw text on tooltip
	for i, name := range finalNames {
		rl.DrawText(
			name,
			tooltip.X+textPad,
			tooltip.Y+textPad+fontSize*int32(i),
			fontSize,
			rl.Black,
		)
	}
}

func drawFooter() {
	footerX := screen.W - offset.X - footerW
	footerY := grid.H + offset.Y*2 + fontSize*2

	text := "Drop file to change sample"
	textW := rl.MeasureText(text, footerFontSize)

	rl.DrawText(text, screen.W-textW-offset.X, grid.H+offset.Y*2, footerFontSize, rl.Black)

	texts := []string{
		fmt.Sprintf("Scale: x%.2f", scale),
		fmt.Sprintf("Cell.W: %.2f", cell.W),
		fmt.Sprint("Cell.H: ", cell.H),
	}

	rl.DrawText(strings.Join(texts, "\n"), footerX+textPad, footerY+textPad, footerFontSize, rl.Black)
	rl.DrawRectangleLines(footerX, footerY, footerW, int32(len(texts)+1)*footerFontSize+textPad*2, rl.Black)
}

func DrawLoop(sample map[string]string) {
	crons := stringsToCrons(sample)
	coords := cronsToCoords(crons)

	gridCoords := [][]GridCoord{}

	groupByScroll := int32(0)

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 700, "Cuckoo")
	rl.SetWindowMinSize(800, 700)

	font = rl.GetFontDefault()

	for !rl.WindowShouldClose() {
		screen.W = int32(rl.GetScreenWidth())
		screen.H = int32(rl.GetScreenHeight())

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
		}

		// Recalculate grid and coordinates only when screen changes size
		if screen.H != prevScreen.H || screen.W != prevScreen.W {
			grid.W = screen.W - 40
			grid.H = screen.H - 240

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

		if zoom != prevZoom || zoomSlider != prevZoomSlider && zoom > 1 {
			zoomOffset = zoomSlider * (scale - 1)

			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Save state for next frame
		prevScreen.W = screen.W
		prevScreen.H = screen.H

		prevGroupBy = groupBy
		prevStepMin = stepMin

		prevZoom = zoom
		prevZoomSlider = zoomSlider
	}

	rl.CloseWindow()
}
