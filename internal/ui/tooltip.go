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

	maxCronW := int32(0)
	maxNameW := int32(0)

	nRows := 0

	// TODO the code below this point to the end of the function can be heavily
	// optimized and reduced

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

	result := map[string]map[string][]string{}
	for time, cronKeys := range crons {
		for cron, names := range cronKeys {
			if w := rl.MeasureText(cron, FontSize) + TextPad + FontRadius; w > maxCronW {
				maxCronW = w
			}

			counts := CountDuplicates(names)

			for name, count := range counts {
				s := fmt.Sprintf("%s (%d)", name, count)

				if _, ok := result[time]; !ok {
					result[time] = make(map[string][]string)
				}

				result[time][cron] = append(result[time][cron], s)

				if w := rl.MeasureText(s, FontSize); w > maxNameW {
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

	drawTooltipRec(tooltip.ToFloat32())

	row := int32(0)

	// TODO optimize so this doesn't have to run every time
	times := slices.Collect(maps.Keys(result))
	sort.Slice(times, func(i, j int) bool {
		return SortAlphabetically(times[i], times[j])
	})

	for _, time := range times {
		// Draw text on tooltip
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

		// TODO optimize so this doesn't have to run every time
		crons := slices.Collect(maps.Keys(result[time]))
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

			// TODO optimize so this doesn't have to run every time
			sort.Slice(result[time][cron], func(i, j int) bool {
				return SortAlphabetically(result[time][cron][i], result[time][cron][j])
			})

			row++

			for i, name := range result[time][cron] {
				rl.DrawText(
					name,
					tooltip.X+TextPad,
					tooltip.Y+TextPad+FontSize*(int32(i)+row),
					FontSize,
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
	boxRoundness := 2 * BoxRadius / MinF32(rec.Height, rec.Width)

	rl.DrawRectangleRounded(rec, boxRoundness, BoxSegments, rl.White)
	rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, BoxSegments, 2, rl.Black)
}
