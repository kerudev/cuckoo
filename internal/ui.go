package cuckoo

import (
	"fmt"
	"maps"
	"math"
	"slices"
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
var fontRadius = fontSize / 2
var textPad = int32(8)

var boxRadius = float32(8)
var boxSegments = int32(8)
var boxSize = int32(20)
var boxBorder = int32(1)
var boxPad = boxSize + boxBorder*2
var coordRadius = float32(4)

var gridBorder = int32(2)

var footerW = int32(120)
var footerFontSize = int32(16)

var weekdays = []Weekday{
	NewWeekday(rl.Red),
	NewWeekday(rl.Orange),
	NewWeekday(rl.Gold),
	NewWeekday(rl.Green),
	NewWeekday(rl.Blue),
	NewWeekday(rl.Purple),
	NewWeekday(rl.Pink),
}

// Internal
var offset = Vector2Int32{X: 20, Y: 20}
var cell = Cell{W: 0, H: 0}

var grid = Grid{Cols: INITIAL_COLS, Rows: INITIAL_ROWS}
var gridHighestY = float32(0.0)

var zoomOffset = float32(0.0)
var zoomBase = float32(0.0)
var zoomFactor = float32(1.0)
var zoomScale = float32(1.0)

var mouseOver = make([][]GridCoord, 7)
var totalOver = 0

// State
var screen = NewState(Screen{W: 0, H: 0})
var mouse = NewState(rl.Vector2{X: 0, Y: 0})
var isMouseLocked = NewState(false)

var groupBy = NewState(GroupByWdHourMin)
var stepMin = NewState(StepMin1)

var zoom = NewState(float32(1.0))
var zoomSlider = NewState(float32(0.0))

// User options
var userOptions = UserOptions{
	drawCoords: true,
	drawLines:  true,
	drawGrid:   true,
	drawFade:   true,
}

var position = PositionGrid

func drawGridLines() {
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
}

func drawCoordsLines(coords []GridCoord, color rl.Color) {
	// Draw lines that connect coordinates
	for k := 0; k < len(coords)-1; k++ {
		start := coords[k].Vector2()
		end := coords[k+1].Vector2()

		rl.DrawLineEx(start, end, float32(gridBorder), color)
	}
}

func drawFade(coord GridCoord, next GridCoord, wd int) {
	mid := Vector2Int32{}

	alpha0 := float32(255)
	alpha1 := float32(255)
	alpha2 := float32(255)

	recX := int32(0)
	recY := int32(0)
	recAlpha := float32(0)

	if coord.Y < next.Y {
		/**
		 * (0) x
		 *     |\
		 *     | \
		 *     |  \
		 * (1) x---x (2)
		 *
		 * - 0: coord
		 * - 1: mid
		 * - 2: next
		 */

		mid.X = int32(coord.X)
		mid.Y = int32(next.Y)

		alpha0 *= coord.OrigY / float32(gridHighestY)
		alpha1 *= next.OrigY / float32(gridHighestY)
		alpha2 *= next.OrigY / float32(gridHighestY)

		recX = int32(mid.X)
		recY = int32(mid.Y)
		recAlpha = next.OrigY
	} else {
		/**
		 *         x (2)
		 *        /|
		 *       / |
		 *      /  |
		 * (0) x---x (1)
		 *
		 * - 0: coord
		 * - 1: mid
		 * - 2: next
		 */

		mid.X = int32(next.X)
		mid.Y = int32(coord.Y)

		alpha0 *= coord.OrigY / float32(gridHighestY)
		alpha1 *= coord.OrigY / float32(gridHighestY)
		alpha2 *= next.OrigY / float32(gridHighestY)

		recX = int32(coord.X)
		recY = int32(coord.Y)
		recAlpha = coord.OrigY
	}

	color := weekdays[wd].color

	// All draw calls use integers to avoid:
	// - Drawing the same pixel twice (darker color)
	// - Not drawing a pixel (white pixel)

	// Draw triangle with faded vertices
	rl.Begin(rl.Triangles)
	rl.Color4ub(color.R, color.G, color.B, uint8(alpha0))
	rl.Vertex2i(int32(coord.X), int32(coord.Y))
	rl.Color4ub(color.R, color.G, color.B, uint8(alpha1))
	rl.Vertex2i(mid.X, mid.Y)
	rl.Color4ub(color.R, color.G, color.B, uint8(alpha2))
	rl.Vertex2i(int32(next.X), int32(next.Y))
	rl.End()

	// Draw gradient below graph
	w := int32(next.X) - int32(coord.X)
	h := grid.H + offset.Y - int32(mid.Y)

	// Calculate rectangle fade based on highest coordinate
	recColor := rl.Fade(color, recAlpha*cell.H/(float32(gridHighestY)*cell.H))

	rl.DrawRectangleGradientV(recX, recY, w, h, recColor, weekdays[wd].faded)
}

func drawGrid(gridCoords [][]GridCoord) {
	// Set all values that depend on the previous frame
	cols := grid.Cols
	if groupBy.Equals(GroupByWdHour) {
		cols += 1
	}

	scroll := rl.GetMouseWheelMove()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		// Horizontal scroll
		calc := float32(cell.W) / (zoomScale * zoomFactor * 2)

		if scroll > 0 {
			zoomSlider.Val += calc
		} else if scroll < 0 {
			zoomSlider.Val -= calc
		}
	} else {
		// Vertical scroll
		zoom.Set(rl.Clamp(zoom.Val+scroll, 1, 9))
		zoomBase = float32(grid.W) / float32(grid.Cols)

		zoomFactor = (zoom.Val - 1) / 8.0
		zoomScale = float32(math.Pow(float64(grid.W)/float64(zoomBase), float64(zoomFactor)))

		zoomOffset = zoomSlider.Val * (zoomScale - 1)
	}

	cell.W = zoomBase * zoomScale

	if userOptions.drawGrid {
		drawGridLines()
	}

	// Draw background line on mouse over
	bg := rl.NewColor(200, 230, 250, 80)
	bgX := float32(offset.X) - zoomOffset

	for range cols {
		mouseInX := bgX < mouse.Val.X && mouse.Val.X <= bgX+cell.W
		mouseInY := float32(offset.Y) < mouse.Val.Y && mouse.Val.Y <= float32(grid.H)

		if mouseInX && mouseInY {
			bgRec := rl.RectangleInt32{X: int32(bgX) + boxBorder, Y: offset.Y, Width: int32(cell.W) - boxBorder*2, Height: grid.H}
			rl.DrawRectangleRec(bgRec.ToFloat32(), bg)
		}

		bgX += cell.W

		if zoom.Equals(1) {
			zoomSlider.Set(rl.Clamp(mouse.Val.X-cell.W, 0, float32(grid.W)))
		}
	}

	// Draw coordinates in layers by weekday
	for wd, dayCoords := range gridCoords {
		if weekdays[wd].status != StatusOn {
			continue
		}

		if userOptions.drawLines {
			drawCoordsLines(dayCoords, weekdays[wd].color)
		}

		if !userOptions.drawCoords && !userOptions.drawFade {
			continue
		}

		// Draw coordinates
		for i, coord := range dayCoords {
			if userOptions.drawFade {
				// Skip drawing gradient after last coordinate
				if i+1 >= len(dayCoords) {
					continue
				}

				next := dayCoords[i+1]

				drawFade(coord, next, wd)
			}

			if userOptions.drawCoords {
				// Skip if coord is off the grid (left)
				if coord.X < float32(offset.X) {
					continue
				}

				// Stop if coord is off the grid (right)
				if coord.X > float32(grid.W+offset.X) {
					break
				}

				rl.DrawCircle(int32(coord.X), int32(coord.Y), coordRadius, weekdays[wd].color)
			}
		}
	}

	// Draw zoom slider
	if zoom.Val > 1 {
		scrollW := grid.W - gridBorder*2
		rg.SetStyle(rg.SLIDER, rg.SLIDER_WIDTH, rg.PropertyValue(float32(scrollW)/zoomScale))

		zoomSliderRec := rl.RectangleInt32{X: offset.X + gridBorder, Y: grid.H + textPad, Width: scrollW, Height: 10}
		rg.Slider(zoomSliderRec.ToFloat32(), "", "", &zoomSlider.Val, 0, float32(grid.W))
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
	step := grid.HighestRow / grid.Rows
	last := (grid.HighestRow / step) * step

	for row := range grid.HighestRow + 1 {
		if grid.HighestRow > ROWS_CAP && row%step != 0 {
			continue
		}

		if row == last {
			row = grid.HighestRow
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

func drawUIOptions(groupByScroll *int32) {
	// Draw option - GroupBy
	rl.DrawText("Group by", offset.X, grid.H+offset.Y*2+textPad, fontSize, rl.Black)
	groupByRec := rl.RectangleInt32{X: offset.X, Y: grid.H + offset.Y*3, Width: 100, Height: 31*2 + 1}
	groupByIdx := int32(groupBy.Val)
	rg.ListView(groupByRec.ToFloat32(), "Wd+Hour;Wd+Hour+Min", groupByScroll, &groupByIdx)

	// Prevent ListView from having nothing selected
	if groupByIdx >= 0 {
		groupBy.Set(GroupBy(groupByIdx))
	}

	// Draw option - Weekdays

	// Check the implementation of GuiLoadStyleDefault for additional keys
	// https://github.com/raysan5/raygui/blob/master/src/raygui.h

	rl.DrawText("Weekdays", 120+offset.X, grid.H+offset.Y*2+textPad, fontSize, rl.Black)

	def_BORDER_WIDTH := rg.GetStyle(rg.BUTTON, rg.BORDER_WIDTH)

	def_TEXT_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED)
	def_TEXT_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED)

	def_BORDER_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL)
	def_BASE_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL)
	def_BORDER_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED)
	def_BASE_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED)
	def_BORDER_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED)
	def_BASE_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED)

	rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, 1)

	for wd := range weekdays {
		status := weekdays[wd].status
		hex := rg.NewColorPropertyValue(weekdays[wd].color)
		black := rg.NewColorPropertyValue(rl.Black)

		// Set styles based on status
		switch status {
		case StatusDisabled:
			base := rg.NewColorPropertyValue(rl.RayWhite)

			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, def_BORDER_COLOR_NORMAL)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, base)

		case StatusOff:
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, black)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, black)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, lerpHex(hex, 0.9))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, lerpHex(hex, 0.8))

		case StatusOn:
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, black)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, black)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, lerpHex(hex, 0.8))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, lerpHex(hex, 0.7))
		}

		button := rl.RectangleInt32{
			X:      120 + offset.X + boxPad*int32(wd),
			Y:      grid.H + offset.Y*3,
			Width:  boxSize,
			Height: boxSize,
		}

		active := status.Bool()
		rg.Toggle(button.ToFloat32(), strconv.Itoa(wd), &active)

		if status != StatusDisabled {
			weekdays[wd].status = StatusFromBool(active)

			if !active && all(weekdays, func(wd Weekday) bool { return wd.status != StatusOn }) {
				weekdays[wd].status = StatusOn
			}
		}

		// Reset style to defaults
		rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, def_BORDER_WIDTH)

		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, def_TEXT_COLOR_PRESSED)
		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, def_TEXT_COLOR_FOCUSED)

		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, def_BORDER_COLOR_NORMAL)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, def_BASE_COLOR_NORMAL)
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, def_BASE_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_PRESSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, def_BASE_COLOR_PRESSED)
	}

	// Draw option - StepMin
	if groupBy.Equals(GroupByWdHourMin) {
		rl.DrawText("Step of x minutes", 120+offset.X, grid.H+offset.Y*4+textPad, fontSize, rl.Black)
		stepMinRec := rl.RectangleInt32{X: 120 + offset.X, Y: grid.H + offset.Y*5, Width: boxSize, Height: boxSize}

		stepMinIdx := int32(stepMin.Val)
		rg.ToggleGroup(stepMinRec.ToFloat32(), "#113#;5;10;15;20;30", &stepMinIdx)

		stepMin.Set(StepMin(stepMinIdx))
	}
}

func drawUserOptions(positionScroll *int32) {
	// User option - TooltipPosition
	rl.DrawText("Tooltip position", offset.X, grid.H+offset.Y*7+textPad, fontSize, rl.Black)
	positionRec := rl.RectangleInt32{X: offset.X, Y: grid.H + offset.Y*8, Width: 100, Height: 31*2 + 1}

	positionIdx := int32(position)
	rg.ListView(positionRec.ToFloat32(), "Grid;Coordinate", positionScroll, &positionIdx)

	// Prevent ListView from having nothing selected
	if positionIdx >= 0 {
		position = TooltipPosition(positionIdx)
	}

	// User option - Draw options
	rl.DrawText("Draw options", 120+offset.X, grid.H+offset.Y*7+textPad, fontSize, rl.Black)

	drawCoordsIcon := ""
	if userOptions.drawCoords {
		drawCoordsIcon = "#212#"
	} else {
		drawCoordsIcon = "#213#"
	}

	options := []ToggleParams{
		{drawCoordsIcon, &userOptions.drawCoords},
		{"#127#", &userOptions.drawLines},
		{"#94#", &userOptions.drawFade},
		{"#97#", &userOptions.drawGrid},
	}

	toggleRec := rl.RectangleInt32{X: 120 + offset.X, Y: grid.H + offset.Y*8, Width: boxSize, Height: boxSize}

	for _, params := range options {
		rg.Toggle(toggleRec.ToFloat32(), params.Icon, params.Ptr)
		toggleRec.X += boxPad
	}

	if !userOptions.drawCoords && !userOptions.drawLines && !userOptions.drawFade {
		userOptions.drawCoords = true
	}
}

func drawTooltip(gridCoords [][]GridCoord) {
	if isMouseLocked.HasChanged() || zoom.HasChanged() || !isMouseLocked.Val && mouse.HasChanged() {
		mouseOver = make([][]GridCoord, 7)
		totalOver = 0

		// Get coords where mouse is over
		for wd, dayCoords := range gridCoords {
			// If a day is not on, there are no coordinates to check
			if weekdays[wd].status != StatusOn {
				continue
			}

			for _, coord := range dayCoords {
				// If the coordinate is not on the same Y range, skip it
				if !(mouse.Val.Y >= coord.Y-coordRadius && mouse.Val.Y <= coord.Y+coordRadius) {
					continue
				}

				// If the coordinate is behind the mouse, don't check collisions
				if mouse.Val.X > coord.X+coordRadius {
					continue
				}

				// If the coordinate is ahead the mouse, don't keep iterating
				if mouse.Val.X+20 <= coord.X {
					break
				}

				if rl.CheckCollisionPointCircle(mouse.Val, coord.Vector2(), coordRadius) {
					mouseOver[wd] = append(mouseOver[wd], coord)
					totalOver++
				}
			}
		}
	}

	// If mouse is not over any coordinate, return
	if totalOver == 0 {
		return
	}

	maxCronW := int32(0)
	maxNameW := int32(0)

	nRows := 0

	// TODO the code below this point to the end of the function can be heavily
	// optimized and reduced

	crons := map[string]map[string][]string{}
	for _, coords := range mouseOver {
		for _, coord := range coords {
			for _, job := range coord.Jobs {
				time := job.AsTime()

				if _, ok := crons[time]; !ok {
					crons[time] = make(map[string][]string)
				}

				crons[time][job.Cron] = append(crons[time][job.Cron], job.Name)
			}
		}
	}

	result := map[string]map[string][]string{}
	for time, cronKeys := range crons {
		for cron, names := range cronKeys {
			if w := rl.MeasureText(cron, fontSize) + textPad + fontRadius; w > maxCronW {
				maxCronW = w
			}

			counts := countDuplicates(names)

			for name, count := range counts {
				s := fmt.Sprintf("%s (%d)", name, count)

				if _, ok := result[time]; !ok {
					result[time] = make(map[string][]string)
				}

				result[time][cron] = append(result[time][cron], s)

				if w := rl.MeasureText(s, fontSize); w > maxNameW {
					maxNameW = w
				}
			}

			// Add one line for the cron string and another for spacing
			nRows += len(counts) + 1 + 1
		}

		nRows += 2
	}

	maxW := max(maxCronW, maxNameW)

	// Prepare tooltip
	tooltip := rl.RectangleInt32{
		Width:  maxW + textPad*2,
		Height: fontSize * int32(nRows),
	}

	switch position {
	case PositionGrid:
		pad := offset.X * 2

		tooltip.X = pad
		tooltip.Y = pad

		// Move tooltip to the right when coordinates are on the left side
		if !isMouseLocked.Val && tooltip.Width > int32(mouse.Val.X)-pad-offset.X {
			tooltip.X = screen.Val.W - pad - tooltip.Width
		}

	case PositionCoord:
		var base GridCoord

		for _, coord := range mouseOver {
			if len(coord) > 0 {
				base = coord[0]
				break
			}
		}

		tooltip.X = int32(base.X) + textPad
		tooltip.Y = int32(base.Y) - textPad

		// Move tooltip to the left when it renders out of the grid
		if tooltip.X+tooltip.Width > offset.X+grid.W {
			tooltip.X = int32(base.X) - textPad - tooltip.Width
		}
	}

	drawTooltipRec(tooltip.ToFloat32())

	row := int32(0)

	// TODO optimize so this doesn't have to run every time
	times := slices.Collect(maps.Keys(result))
	sort.Slice(times, func(i, j int) bool {
		return sortAlphabetically(times[i], times[j])
	})

	for _, time := range times {
		// Draw text on tooltip
		rg.DrawIcon(
			rg.ICON_CLOCK,
			tooltip.X+textPad,
			tooltip.Y+textPad+2+fontSize*row,
			1,
			rl.Black,
		)

		rl.DrawText(
			time,
			tooltip.X+textPad*4,
			tooltip.Y+textPad+2+fontSize*row,
			16,
			rl.Black,
		)

		row += 2

		// TODO optimize so this doesn't have to run every time
		crons := slices.Collect(maps.Keys(result[time]))
		sort.Slice(crons, func(i, j int) bool {
			return sortAlphabetically(crons[i], crons[j])
		})

		for _, cron := range crons {
			wds := parseCronField(strings.Split(cron, " ")[4], 0, 6)

			segments := float32(len(wds))
			for _, wd := range wds {
				if weekdays[wd].status != StatusOn {
					segments--
				}
			}

			angleFactor := float32(360) / segments
			angle := float32(0)

			for _, wd := range wds {
				if weekdays[wd].status != StatusOn {
					continue
				}

				rl.DrawCircleSector(
					rl.Vector2{
						X: float32(tooltip.X + textPad + fontRadius),
						Y: float32(tooltip.Y + textPad + fontSize*row + fontRadius),
					},
					float32(fontRadius),
					angle,
					angle+angleFactor,
					8,
					weekdays[wd].color,
				)

				angle += angleFactor
			}

			// Draw crons and their count
			rl.DrawText(
				cron,
				tooltip.X+textPad+4*4,
				tooltip.Y+textPad+fontSize*row,
				fontSize,
				rl.Black,
			)

			// TODO optimize so this doesn't have to run every time
			sort.Slice(result[time][cron], func(i, j int) bool {
				return sortAlphabetically(result[time][cron][i], result[time][cron][j])
			})

			row++

			for i, name := range result[time][cron] {
				rl.DrawText(
					name,
					tooltip.X+textPad,
					tooltip.Y+textPad+fontSize*(int32(i)+row),
					fontSize,
					rl.Black,
				)
			}

			row += int32(len(result[time][cron])) + 1
		}
	}
}

func drawTooltipRec(rec rl.Rectangle) {
	// Raylib computes the radius using the formula:
	// float radius = (rec.width > rec.height)? (rec.height*roundness)/2 : (rec.width*roundness)/2;
	//
	// The radius depends on the "roundness", which must be known beforehand so
	// the radius is always the same.
	boxRoundness := 2 * boxRadius / minF32(rec.Height, rec.Width)

	rl.DrawRectangleRounded(rec, boxRoundness, boxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, boxSegments, 2, rl.Black)
}

func drawFooter() {
	footerX := screen.Val.W - offset.X - footerW
	footerY := grid.H + offset.Y*2 + fontSize*2

	text := "Drop file to change sample"
	textW := rl.MeasureText(text, footerFontSize)

	rl.DrawText(text, screen.Val.W-textW-offset.X, grid.H+offset.Y*2, footerFontSize, rl.Black)

	texts := []string{
		fmt.Sprintf("Scale: x%.2f", zoomScale),
		fmt.Sprintf("Cell.W: %.2f", cell.W),
		fmt.Sprintf("Cell.H: %.2f", cell.H),
		fmt.Sprint("[L]ocked: ", isMouseLocked.Val),
	}

	rl.DrawText(strings.Join(texts, "\n"), footerX+textPad, footerY+textPad, footerFontSize, rl.Black)
	rl.DrawRectangleLines(footerX, footerY, footerW, int32(len(texts)+1)*footerFontSize+textPad*2, rl.Black)
}

func DrawLoop(sample map[string]string) {
	crons := stringsToCrons(sample)
	coords := cronsToCoords(crons)

	gridCoords := [][]GridCoord{}

	groupByScroll := int32(0)
	positionScroll := int32(0)

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 700, "Cuckoo")
	rl.SetWindowMinSize(800, 700)

	font = rl.GetFontDefault()

	for !rl.WindowShouldClose() {
		screen.Val.W = int32(rl.GetScreenWidth())
		screen.Val.H = int32(rl.GetScreenHeight())

		mouse.Set(rl.GetMousePosition())

		// Recalculate grid and coordinates only when screen changes size
		if screen.HasChanged() {
			grid.W = screen.Val.W - 40
			grid.H = screen.Val.H - 240

			gridCoords = coordToGrid(coords, &grid)
		}

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

		// Input event handling (keyboard)
		if rl.IsKeyPressed(rl.KeyL) {
			isMouseLocked.Set(!isMouseLocked.Val)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawGrid(gridCoords)

		drawUIOptions(&groupByScroll)

		// Vertical line
		lineY := grid.H + offset.Y*7 - textPad/2
		rl.DrawLine(offset.X, lineY, 290, lineY, rl.Gray)

		drawUserOptions(&positionScroll)

		// Horizontal line
		lineX := 150 + offset.X + boxPad*6
		rl.DrawLine(lineX, grid.H+offset.Y*2+textPad, lineX, screen.Val.H-offset.Y, rl.Gray)

		drawFooter()
		drawTooltip(gridCoords)

		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if stepMin.HasChanged() {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Recalculate coordinates based on group by
		if groupBy.HasChanged() {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		if zoom.HasChanged() || zoomSlider.HasChanged() && zoom.Val > 1 {
			// NOTE: unlock mouse as mouseOver as coordinates are recalculated
			// when zoom changes. Might be good to change this at some point.
			isMouseLocked.Set(false)

			zoomOffset = zoomSlider.Val * (zoomScale - 1)

			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Save state for next frame
		screen.Update()
		mouse.Update()
		isMouseLocked.Update()

		groupBy.Update()
		stepMin.Update()

		zoom.Update()
		zoomSlider.Update()
	}

	rl.CloseWindow()
}
