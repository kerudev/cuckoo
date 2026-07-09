package ui

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawTooltip(gridCoords [][]GridCoord) {
	stateChanged := S_IsMouseLocked.HasChanged() ||
		S_Weekdays.HasChanged() ||
		S_Zoom.HasChanged() ||
		!S_IsMouseLocked.Val && S_Mouse.HasChanged()

	if stateChanged && !rg.IsLocked() {
		MouseOver = [WEEKDAYS][]GridCoord{}
		TotalOver = 0

		// Get coords where Mouse is over
		for wd, dayCoords := range gridCoords {
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

	// If Mouse is not over any coordinate, return
	if TotalOver == 0 {
		return
	}

	for wd, dayCoords := range MouseOver {
		// If a day is not on, there are no coordinates to check
		if S_Weekdays.Val[wd].Status != StatusOn {
			continue
		}

		if len(dayCoords) == 0 {
			continue
		}

		for _, coord := range dayCoords {
			faded := rl.ColorLerp(S_Weekdays.Val[wd].Color, rl.White, 0.3)
			rl.DrawCircle(int32(coord.X), int32(coord.Y), CoordRadius, faded)

			// https://github.com/raysan5/raylib/blob/6735907/src/rshapes.c#L1463
			rl.DrawRing(coord.Vector2(), CoordRadius-1, CoordRadius+1, 0, 360, 36, rl.Black)
		}
	}

	nRows := 0
	maxW := int32(0)

	// Extract data from MouseOver
	// Time (HH:MM) -> Cron string -> Job names
	crons := map[string]map[string][]string{}
	for _, coords := range MouseOver {
		for _, coord := range coords {
			for _, job := range coord.Jobs {
				time := fmt.Sprintf("%s (%d)", job.AsTime(), int32(coord.OrigY))

				if _, ok := crons[time]; !ok {
					crons[time] = make(map[string][]string)
				}

				crons[time][job.Cron] = append(crons[time][job.Cron], job.Name)
			}
		}
	}

	schedule := map[string]map[string][]string{}
	for time, crons := range crons {
		for cron, jobs := range crons {
			if w := rl.MeasureText(cron, FontSize) + TextPad + FontRadius; w > maxW {
				maxW = w
			}

			counts := CountDuplicates(jobs)

			for job, count := range counts {
				s := fmt.Sprintf("%s (%d)", job, count)

				if _, ok := schedule[time]; !ok {
					schedule[time] = make(map[string][]string)
				}

				schedule[time][cron] = append(schedule[time][cron], s)

				if w := rl.MeasureText(s, FontSize); w > maxW {
					maxW = w
				}
			}

			// Add one line for the cron string and another for spacing
			nRows += len(counts) + 1 + 1
		}

		// Spacing
		nRows += 2
	}

	// Resize tooltip when mouse is unlocked or weekdays changed
	if !S_IsMouseLocked.Val || S_Weekdays.HasChanged() {
		S_TooltipHasOverflow.Set(false)

		// Prepare Tooltip
		Tooltip.Width = maxW + TextPad*2
		Tooltip.Height = FontSize * int32(nRows)

		S_TooltipScrollMax.Set(Tooltip.Height)

		padX := Offset.X * 2
		padY := Offset.Y * 2

		switch Position {
		case PositionGrid:
			Tooltip.X = padX
			Tooltip.Y = padY

			// Clamp height when it's too large
			if Tooltip.Height > Grid.Height-padY {
				Tooltip.Width += TooltipScrollW
				Tooltip.Height = Grid.Height - padY

				S_TooltipScrollMax.Val -= Tooltip.Height
				S_TooltipHasOverflow.Set(true)
			}

			// Move to the right when coordinates are on the left side
			if Tooltip.Width > int32(S_Mouse.Val.X)-padX-Offset.X {
				Tooltip.X = S_Screen.Val.W - padX - Tooltip.Width
			}

		case PositionCoord:
			var base GridCoord

			for _, coords := range MouseOver {
				if len(coords) > 0 {
					base = coords[0]
					break
				}
			}

			Tooltip.X = int32(base.X) + Offset.X
			Tooltip.Y = int32(base.Y)

			// Do stuff when the rectangle gets out of the grid (below)
			if Tooltip.Height+Tooltip.Y > Grid.Height-TextPad {
				Tooltip.Y -= Tooltip.Height

				// Move upwards
				if Tooltip.Y < Offset.Y {
					Tooltip.Y = Offset.Y * 2
					Tooltip.Height = Clamp(Tooltip.Height, padY, Tooltip.Height)
				}

				// Clamp height when it's too large
				if Tooltip.Height > Grid.Height-padY {
					Tooltip.Width += TooltipScrollW
					Tooltip.Height = Grid.Height - padY

					S_TooltipScrollMax.Val -= Tooltip.Height
					S_TooltipHasOverflow.Set(true)
				}
			}

			// Move to the left when it renders out of the Grid
			if Tooltip.X+Tooltip.Width > Grid.Width {
				Tooltip.X = int32(base.X) - Offset.X - Tooltip.Width
			}
		}
	}

	drawTooltipRec()

	// TODO optimize to reduce draw calls when text is out of the tooltip
	rl.BeginScissorMode(Tooltip.X, Tooltip.Y, Tooltip.Width, Tooltip.Height)
	drawTooltipText(schedule)
	rl.EndScissorMode()
}

func drawTooltipRec() {
	rec := Tooltip.ToFloat32()

	// Raylib computes the radius using the formula:
	// float radius = (rec.width > rec.height)? (rec.height*roundness)/2 : (rec.width*roundness)/2;
	//
	// The radius depends on the "roundness", which must be known beforehand so
	// the radius is always the same.
	boxRoundness := BoxDiameter / min(rec.Height, rec.Width)

	rl.DrawRectangleRounded(rec, boxRoundness, BoxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, BoxSegments, 2, rl.Black)

	if S_TooltipHasOverflow.Val {
		rg.SetStyle(rg.SCROLLBAR, rg.BORDER_WIDTH, rg.GetStyle(rg.SLIDER, rg.BORDER_WIDTH))

		rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_NORMAL, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_NORMAL))
		rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_FOCUSED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_FOCUSED))
		rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_PRESSED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_PRESSED))
		rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_DISABLED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_DISABLED))

		rg.SetStyle(rg.BUTTON, rg.BASE_COLOR_NORMAL, rg.GetStyle(rg.SLIDER, rg.BASE_COLOR_NORMAL))

		tooltipScrollRec := rl.RectangleInt32{
			X:      Tooltip.X + Tooltip.Width - TooltipScrollW,
			Y:      Tooltip.Y + int32(BoxRadius),
			Width:  TooltipScrollW,
			Height: Tooltip.Height - int32(BoxDiameter),
		}

		// Allow scroll just when mouse is over tooltip
		if rl.CheckCollisionPointRec(S_Mouse.Val, rec) {
			scroll := int32(rl.GetMouseWheelMove()) * int32(16)

			if scroll != 0 {
				S_TooltipScroll.Val -= scroll
				S_TooltipScroll.Set(Clamp(S_TooltipScroll.Val, 0, S_TooltipScrollMax.Val))
			}
		}

		S_TooltipScroll.Set(rg.ScrollBar(tooltipScrollRec.ToFloat32(), S_TooltipScroll.Val, 0, S_TooltipScrollMax.Val))
	}
}

func drawTooltipText(data map[string]map[string][]string) {
	row := int32(0)

	// Sort HH:MM keys
	times := slices.Collect(maps.Keys(data))
	sort.Slice(times, func(i, j int) bool {
		return SortAlphabetically(times[i], times[j])
	})

	for _, time := range times {
		// Draw clock icon and time
		rg.DrawIcon(
			rg.ICON_CLOCK,
			Tooltip.X+TextPad,
			Tooltip.Y+TextPad+2+FontSize*row-S_TooltipScroll.Val,
			1,
			rl.Black,
		)

		rl.DrawText(
			time,
			Tooltip.X+TextPad*4,
			Tooltip.Y+TextPad+2+FontSize*row-S_TooltipScroll.Val,
			TooltipTimeFontSize,
			rl.Black,
		)

		// Spacing
		row += 2

		// Sort cron strings
		crons := slices.Collect(maps.Keys(data[time]))
		sort.Slice(crons, func(i, j int) bool {
			return SortAlphabetically(crons[i], crons[j])
		})

		for _, cron := range crons {
			wds := ParseCronField(strings.Split(cron, " ")[4], 0, 6)

			segments := float32(len(wds))
			for _, wd := range wds {
				if S_Weekdays.Val[wd].Status != StatusOn {
					segments--
				}
			}

			angleFactor := float32(360) / segments
			angle := float32(0)

			for _, wd := range wds {
				if S_Weekdays.Val[wd].Status != StatusOn {
					continue
				}

				rl.DrawCircleSector(
					rl.Vector2{
						X: float32(Tooltip.X + TextPad + FontRadius),
						Y: float32(Tooltip.Y + TextPad + FontSize*row + FontRadius - S_TooltipScroll.Val),
					},
					float32(FontRadius),
					angle,
					angle+angleFactor,
					8,
					S_Weekdays.Val[wd].Color,
				)

				angle += angleFactor
			}

			// Draw crons and their count
			rl.DrawText(
				cron,
				Tooltip.X+TextPad+4*4,
				Tooltip.Y+TextPad+FontSize*row-S_TooltipScroll.Val,
				FontSize,
				rl.Black,
			)

			jobs := data[time][cron]

			// Sort job names
			sort.Slice(jobs, func(i, j int) bool {
				return SortAlphabetically(jobs[i], jobs[j])
			})

			row++

			for i, job := range jobs {
				rl.DrawText(
					job,
					Tooltip.X+TextPad,
					Tooltip.Y+TextPad+FontSize*(int32(i)+row)-S_TooltipScroll.Val,
					FontSize,
					rl.Black,
				)
			}

			row += int32(len(jobs)) + 1
		}
	}
}
