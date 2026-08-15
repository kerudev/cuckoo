package ui

import (
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawGrid() {
	// Set all values that depend on the previous frame
	cols := C_Grid.Cols
	if S_GroupBy.Eq(GroupByWdHour) {
		cols += 1
	}

	Cell.W = C_Zoom.Base * C_Zoom.Scale

	if UserOpt.DrawGrid {
		drawGridLines()
	}

	// Scissor Mode to prevent drawing pixels outside the grid
	rl.BeginScissorMode(Offset.X, Grid.Y, Grid.Width, Grid.Height)

	// Draw background line on Mouse over
	if S_IsOnWindow.Val {
		bgX := float32(Offset.X) - C_Zoom.Offset

		for range cols {
			mouseInX := bgX < S_Mouse.Val.X && S_Mouse.Val.X <= bgX+Cell.W
			mouseInY := float32(Grid.Y) < S_Mouse.Val.Y && S_Mouse.Val.Y <= float32(Footer.Y)

			if !(mouseInX && mouseInY) {
				bgX += Cell.W

				if bgX >= float32(Grid.Width+Offset.X) {
					break
				}
			}
		}

		bgRec := NewRectangleFromInt32(int32(bgX)+BoxBorder*2, Grid.Y, int32(Cell.W)-BoxBorder*2, Grid.Height)
		rl.DrawRectangleRec(bgRec, GridBGColor)
	}

	// Change slider based on mouse position to "follow" the cursor
	if S_Mouse.HasChanged() && S_Zoom.Eq(1) {
		S_ZoomSlider.Set(Clamp(S_Mouse.Val.X-Cell.W, 0, float32(Grid.Width)))
	}

	// Draw coordinates in layers by weekday
	for wd, dayCoords := range GridCoords {
		if S_Weekdays.Val[wd].Status != StatusOn {
			continue
		}

		if UserOpt.DrawLines {
			drawCoordsLines(dayCoords, S_Weekdays.Val[wd].Color)
		}

		if !UserOpt.DrawFade && !UserOpt.DrawCoords {
			continue
		}

		// Draw coordinates
		for i, coord := range dayCoords {
			if UserOpt.DrawFade {
				// Drawing gradient only before last coordinate
				if i+1 < len(dayCoords) {
					next := dayCoords[i+1]
					drawFade(coord, next, wd)
				}
			}

			if UserOpt.DrawCoords {
				// Skip if coord is off the Grid (left)
				if coord.X < float32(Offset.X) {
					continue
				}

				// Stop if coord is off the Grid (right)
				if coord.X > float32(Grid.Width+Offset.X) {
					break
				}

				rl.DrawCircle(int32(coord.X), int32(coord.Y), CoordRadius, S_Weekdays.Val[wd].Color)
			}
		}
	}

	rl.EndScissorMode()

	// Draw Zoom slider
	if S_Zoom.Val > 1 {
		scrollW := Grid.Width - GridBorder*2
		rg.SetStyle(rg.SLIDER, rg.SLIDER_WIDTH, rg.PropertyValue(float32(scrollW)/C_Zoom.Scale))

		zoomSliderRec := NewRectangleFromInt32(Offset.X+GridBorder, Footer.Y-ZoomSliderH-GridBorder, scrollW, ZoomSliderH)
		rg.Slider(zoomSliderRec, "", "", &S_ZoomSlider.Val, 0, float32(Grid.Width))
	}

	// Draw numbers on X axis
	for col := range cols {
		text := strconv.Itoa(col)

		textW := rl.MeasureTextEx(Font, text, float32(FontSize), 1).X
		textX := Cell.W*float32(col) - textW/2 + float32(Offset.X) - C_Zoom.Offset

		// Clamp number to the left side
		if textX < float32(Offset.X) {
			if textX+Cell.W <= float32(Offset.X+TextPad) {
				continue
			}

			textX = float32(Offset.X)
		}

		// Clamp number to the right side
		if textX > float32(Grid.Width+Offset.X)-textW/2 {
			if textX-Cell.W >= float32(Grid.Width+Offset.X-TextPad*3) {
				continue
			}

			textX = float32(Grid.Width+Offset.X) - textW/2
		}

		rl.DrawText(text, int32(textX), Footer.Y+2, FontSize, rl.Black)
	}

	// Draw numbers on Y axis
	textRect := rl.MeasureTextEx(Font, strconv.Itoa(cols), float32(FontSize), 1)

	nRow := 0
	step := C_Grid.HighestRow / C_Grid.Rows
	last := (C_Grid.HighestRow / step) * step

	for row := range C_Grid.HighestRow + 1 {
		if C_Grid.HighestRow > ROWS_CAP && row%step != 0 {
			continue
		}

		if row == last {
			row = C_Grid.HighestRow
		}

		text := strconv.Itoa(row)
		textSize := rl.MeasureTextEx(Font, strconv.Itoa(row), float32(FontSize), 1)

		textPos := rl.Vector2{
			X: textRect.X + rl.Lerp(0.0, textRect.X-textSize.X, 1),
			Y: textRect.Y + rl.Lerp(0.0, textRect.Y-textSize.Y, 0.5),
		}

		textY := float32(Footer.Y) - Cell.H*float32(nRow) - textPos.Y/2
		nRow++

		rl.DrawText(text, int32(textPos.X-float32(Offset.X)/2), int32(textY), FontSize, rl.Black)
	}

	// Draw Grid container
	rl.DrawRectangleLinesEx(Grid.ToFloat32(), 2, rl.Black)

	// Draw and save mouse over coordinates
	drawMouseOver()
}

func drawGridLines() {
	// Draw vertical lines
	colX := float32(Offset.X) - C_Zoom.Offset

	for range C_Grid.Cols {
		colX += Cell.W

		if colX < float32(Offset.X) {
			continue
		}

		if colX > float32(Grid.Width+Offset.X) {
			break
		}

		rl.DrawLineEx(
			rl.Vector2{X: colX, Y: float32(Grid.Y)},
			rl.Vector2{X: colX, Y: float32(Footer.Y)},
			float32(GridBorder),
			rl.LightGray,
		)
	}

	// Draw horizontal lines
	rowY := float32(Grid.Y)

	for range C_Grid.Rows {
		rowY += Cell.H
		rl.DrawLineEx(
			rl.Vector2{X: float32(Offset.X), Y: rowY},
			rl.Vector2{X: float32(Grid.Width + Offset.Y), Y: rowY},
			float32(GridBorder),
			rl.LightGray,
		)
	}
}

func drawCoordsLines(coords []GridCoord, color rl.Color) {
	// Draw lines that connect coordinates
	for k := 0; k < len(coords)-1; k++ {
		start := coords[k].Vector2()
		end := coords[k+1].Vector2()

		rl.DrawLineEx(start, end, float32(GridBorder), color)
	}
}

// TODO optimize this function, as it notably reduces FPS
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

		alpha0 *= coord.OrigY / float32(C_Grid.HighestY)
		alpha1 *= next.OrigY / float32(C_Grid.HighestY)
		alpha2 *= next.OrigY / float32(C_Grid.HighestY)

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

		alpha0 *= coord.OrigY / float32(C_Grid.HighestY)
		alpha1 *= coord.OrigY / float32(C_Grid.HighestY)
		alpha2 *= next.OrigY / float32(C_Grid.HighestY)

		recX = int32(coord.X)
		recY = int32(coord.Y)
		recAlpha = coord.OrigY
	}

	color := S_Weekdays.Val[wd].Color

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
	h := Footer.Y - int32(mid.Y)

	// Calculate rectangle fade based on highest coordinate
	recColor := rl.Fade(color, recAlpha*Cell.H/(float32(C_Grid.HighestY)*Cell.H))

	rl.DrawRectangleGradientV(recX, recY, w, h, recColor, S_Weekdays.Val[wd].Faded)
}

func drawMouseOver() {
	stateChanged := S_IsMouseLocked.HasChanged() ||
		S_Weekdays.HasChanged() ||
		S_Zoom.HasChanged() ||
		!S_IsMouseLocked.Val && S_Mouse.HasChanged()

	if !BlockUI && stateChanged && S_IsOverGrid.Val {
		return
	}

	MouseOver = [WEEKDAYS][]GridCoord{}
	TotalOver = 0

	// Get coords where Mouse is over
	for wd, dayCoords := range GridCoords {
		// If a day is not on, there are no coordinates to check
		if S_Weekdays.Val[wd].Status != StatusOn {
			continue
		}

		for _, coord := range dayCoords {
			// If the coordinate is not on the same Y range, skip it
			if !(S_MouseWithLock.Val.Y >= coord.Y-CoordRadius && S_MouseWithLock.Val.Y <= coord.Y+CoordRadius) {
				continue
			}

			// If the coordinate is behind the Mouse, don't check collisions
			if S_MouseWithLock.Val.X > coord.X+CoordRadius {
				continue
			}

			// If the coordinate is ahead the Mouse, don't keep iterating
			if S_MouseWithLock.Val.X+20 <= coord.X {
				break
			}

			if rl.CheckCollisionPointCircle(S_MouseWithLock.Val, coord.Vector2(), CoordRadius) {
				MouseOver[wd] = append(MouseOver[wd], coord)
				TotalOver++
			}
		}
	}
}
