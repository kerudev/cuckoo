package ui

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/shared"
	. "github.com/kerudev/cuckoo/internal/utils"
)

// TODO optimize or refactor the whole thing at some point. It's getting too
// complex and hard to follow. There are too many loops and custom types.

func DrawTooltip() {
	// If Mouse is not over any coordinate, return
	if TotalOver == 0 {
		return
	}

	// Resize tooltip when mouse is unlocked or weekdays changed
	// NOTE: check if UI is blocked, as the tooltip should not move when mouse moves around
	if !BlockUI && (!S_IsMouseLocked.Val && S_Mouse.HasChanged() || S_Weekdays.HasChanged() || S_Position.HasChanged()) {
		nRows := 0
		maxW := int32(0)

		// Extract data from MouseOver
		// Time (HH:MM) -> Cron string -> Job names & weekdays
		cronsByTime := map[string]map[string]JobsByWd{}
		for wd, coords := range MouseOver {
			if S_Weekdays.Val[wd].Status != StatusOn {
				continue
			}

			for _, coord := range coords {
				for _, job := range coord.Jobs {
					time := fmt.Sprintf("%s (%d)", job.AsTime(), int32(coord.OrigY))

					if _, ok := cronsByTime[time]; !ok {
						cronsByTime[time] = map[string]JobsByWd{}
					}

					j, ok := cronsByTime[time][job.Cron]
					if !ok {
						j = JobsByWd{
							Jobs: map[string]int{},
							Wds:  []int{},
						}
					}

					j.Jobs[job.Name]++
					j.Wds = append(cronsByTime[time][job.Cron].Wds, wd)

					cronsByTime[time][job.Cron] = j
				}
			}
		}

		// Transform data to strings
		TooltipLines = map[string]map[string]JobsCountsByWd{}
		for time, crons := range cronsByTime {
			for cron, jobsByWd := range crons {
				if w := rl.MeasureText(cron, FontSize) + TextPad + FontRadius; w > maxW {
					maxW = w
				}

				for job, count := range maps.All(jobsByWd.Jobs) {
					s := fmt.Sprintf("%s (%d)", job, count)

					if w := rl.MeasureText(s, FontSize); w > maxW {
						maxW = w
					}

					if _, ok := TooltipLines[time]; !ok {
						TooltipLines[time] = map[string]JobsCountsByWd{}
					}

					j, ok := TooltipLines[time][cron]
					if !ok {
						j = JobsCountsByWd{
							Jobs: []string{},
							Wds:  []int{},
						}
					}

					j.Jobs = append(j.Jobs, s)
					j.Wds = jobsByWd.Wds

					TooltipLines[time][cron] = j
				}

				// Add one line for the cron string and another for spacing
				nRows += len(jobsByWd.Jobs) + 1 + 1
			}

			// Spacing
			nRows += 2
		}

		TooltipHasOverflow = false

		// Prepare Tooltip
		Tooltip.Width = maxW + TextPad*2
		Tooltip.Height = FontSize * int32(nRows)

		S_TooltipScrollMax.Set(Tooltip.Height)

		padX := Offset.X + Offset.Y
		padY := Offset.Y + Grid.Y

		padYInner := Offset.Y * 2

		switch S_Position.Val {
		case PositionGrid:
			Tooltip.X = padX
			Tooltip.Y = padY

			// Clamp height when it's too large
			if Tooltip.Height > Grid.Height-padYInner {
				Tooltip.Width += TooltipScrollW
				Tooltip.Height = Grid.Height - padYInner

				S_TooltipScrollMax.Val -= Tooltip.Height
				TooltipHasOverflow = true
			}

			// Move to the right when coordinates are on the left side
			if Tooltip.Width > int32(S_MouseWithLock.Val.X)-padX-Offset.X {
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
				if Tooltip.Y < Grid.Y {
					Tooltip.Y = Grid.Y + Offset.Y
					Tooltip.Height = Clamp(Tooltip.Height, padY, Tooltip.Height)
				}

				// Clamp height when it's too large
				if Tooltip.Height > Grid.Height-padYInner {
					Tooltip.Width += TooltipScrollW
					Tooltip.Height = Grid.Height - padYInner

					S_TooltipScrollMax.Val -= Tooltip.Height
					TooltipHasOverflow = true
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
	drawTooltipText()
	rl.EndScissorMode()
}

func drawTooltipRec() {
	tooltip := Tooltip.ToFloat32()

	// Raylib computes the radius using the formula:
	// float radius = (rec.width > rec.height)? (rec.height*roundness)/2 : (rec.width*roundness)/2;
	//
	// The radius depends on the "roundness", which must be known beforehand so
	// the radius is always the same.
	boxRoundness := BoxDiameter / min(tooltip.Height, tooltip.Width)

	rl.DrawRectangleRounded(tooltip, boxRoundness, BoxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(tooltip, boxRoundness, BoxSegments, 2, rl.Black)

	if !TooltipHasOverflow {
		return
	}

	rg.SetStyle(rg.SCROLLBAR, rg.BORDER_WIDTH, rg.GetStyle(rg.SLIDER, rg.BORDER_WIDTH))

	rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_NORMAL, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_NORMAL))
	rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_FOCUSED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_FOCUSED))
	rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_PRESSED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_PRESSED))
	rg.SetStyle(rg.LISTVIEW, rg.BORDER_COLOR_DISABLED, rg.GetStyle(rg.SLIDER, rg.BORDER_COLOR_DISABLED))

	rg.SetStyle(rg.BUTTON, rg.BASE_COLOR_NORMAL, rg.GetStyle(rg.SLIDER, rg.BASE_COLOR_NORMAL))

	tooltipScrollRec := rl.NewRectangle(
		tooltip.X+tooltip.Width-float32(TooltipScrollW),
		tooltip.Y+BoxRadius,
		float32(TooltipScrollW),
		tooltip.Height-BoxDiameter,
	)

	// Allow scroll just when mouse is over tooltip
	if rl.CheckCollisionPointRec(S_Mouse.Val, tooltip) {
		scroll := int32(rl.GetMouseWheelMove()) * int32(16)

		if scroll != 0 {
			S_TooltipScroll.Val -= scroll
			S_TooltipScroll.Set(Clamp(S_TooltipScroll.Val, 0, S_TooltipScrollMax.Val))
		}
	}

	S_TooltipScroll.Set(rg.ScrollBar(tooltipScrollRec, S_TooltipScroll.Val, 0, S_TooltipScrollMax.Val))
}

func drawTooltipText() {
	row := int32(0)

	// Sort HH:MM keys
	times := slices.Collect(maps.Keys(TooltipLines))
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
		crons := slices.Collect(maps.Keys(TooltipLines[time]))
		sort.Slice(crons, func(i, j int) bool {
			return SortAlphabetically(crons[i], crons[j])
		})

		for _, cron := range crons {
			jobsCount := TooltipLines[time][cron]

			wds := jobsCount.Wds

			segments := float32(len(wds))
			for _, wd := range wds {
				if S_Weekdays.Val[wd].Status != StatusOn {
					segments--
				}
			}

			angleFactor := float32(360) / segments
			angle := float32(270)

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

			jobs := jobsCount.Jobs

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
