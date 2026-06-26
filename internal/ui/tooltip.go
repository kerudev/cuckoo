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
	if IsMouseLocked.HasChanged() || S_Zoom.HasChanged() || !IsMouseLocked.Val && S_Mouse.HasChanged() {
		MouseOver = make([][]GridCoord, 7)
		TotalOver = 0

		// Get coords where Mouse is over
		for wd, dayCoords := range gridCoords {
			// If a day is not on, there are no coordinates to check
			if Weekdays[wd].Status != StatusOn {
				continue
			}

			for _, coord := range dayCoords {
				// If the coordinate is not on the same Y range, skip it
				if !(S_Mouse.Val.Y >= coord.Y-CoordRadius && S_Mouse.Val.Y <= coord.Y+CoordRadius) {
					continue
				}

				// If the coordinate is behind the Mouse, don't check collisions
				if S_Mouse.Val.X > coord.X+CoordRadius {
					continue
				}

				// If the coordinate is ahead the Mouse, don't keep iterating
				if S_Mouse.Val.X+20 <= coord.X {
					break
				}

				if rl.CheckCollisionPointCircle(S_Mouse.Val, coord.Vector2(), CoordRadius) {
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

	nRows := 0
	maxW := int32(0)

	// Extract data from MouseOver
	// Time (HH:MM) -> Cron string -> Job names
	crons := map[string]map[string][]string{}
	for _, coords := range MouseOver {
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

		nRows += 2
	}

	// Prepare tooltip
	tooltip := rl.RectangleInt32{
		Width:  maxW + TextPad*2,
		Height: FontSize * int32(nRows),
	}

	switch Position {
	case PositionGrid:
		pad := Offset.X * 2

		tooltip.X = pad
		tooltip.Y = pad

		// Move tooltip to the right when coordinates are on the left side
		if !IsMouseLocked.Val && tooltip.Width > int32(S_Mouse.Val.X)-pad-Offset.X {
			tooltip.X = S_Screen.Val.W - pad - tooltip.Width
		}

	case PositionCoord:
		var base GridCoord

		for _, coords := range MouseOver {
			if len(coords) > 0 {
				base = coords[0]
				break
			}
		}

		tooltip.X = int32(base.X) + TextPad
		tooltip.Y = int32(base.Y) - TextPad

		// Move tooltip to the left when it renders out of the Grid
		if tooltip.X+tooltip.Width > Offset.X+Grid.W {
			tooltip.X = int32(base.X) - TextPad - tooltip.Width
		}
	}

	drawTooltipRec(tooltip)
	drawTooltipText(tooltip, schedule)
}

func drawTooltipRec(tooltip rl.RectangleInt32) {
	rec := tooltip.ToFloat32()

	// Raylib computes the radius using the formula:
	// float radius = (rec.width > rec.height)? (rec.height*roundness)/2 : (rec.width*roundness)/2;
	//
	// The radius depends on the "roundness", which must be known beforehand so
	// the radius is always the same.
	boxRoundness := 2 * BoxRadius / MinF32(rec.Height, rec.Width)

	rl.DrawRectangleRounded(rec, boxRoundness, BoxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, BoxSegments, 2, rl.Black)
}

func drawTooltipText(tooltip rl.RectangleInt32, data map[string]map[string][]string) {
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
			tooltip.X+TextPad,
			tooltip.Y+TextPad+2+FontSize*row,
			1,
			rl.Black,
		)

		rl.DrawText(
			time,
			tooltip.X+TextPad*4,
			tooltip.Y+TextPad+2+FontSize*row,
			16,
			rl.Black,
		)

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
				if Weekdays[wd].Status != StatusOn {
					segments--
				}
			}

			angleFactor := float32(360) / segments
			angle := float32(0)

			for _, wd := range wds {
				if Weekdays[wd].Status != StatusOn {
					continue
				}

				rl.DrawCircleSector(
					rl.Vector2{
						X: float32(tooltip.X + TextPad + FontRadius),
						Y: float32(tooltip.Y + TextPad + FontSize*row + FontRadius),
					},
					float32(FontRadius),
					angle,
					angle+angleFactor,
					8,
					Weekdays[wd].Color,
				)

				angle += angleFactor
			}

			// Draw crons and their count
			rl.DrawText(
				cron,
				tooltip.X+TextPad+4*4,
				tooltip.Y+TextPad+FontSize*row,
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
					tooltip.X+TextPad,
					tooltip.Y+TextPad+FontSize*(int32(i)+row),
					FontSize,
					rl.Black,
				)
			}

			row += int32(len(jobs)) + 1
		}
	}
}
