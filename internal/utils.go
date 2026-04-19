package cuckoo

import (
	"strconv"
	"strings"
	"unicode"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func coordToVec2(coord GridCoord) rl.Vector2 {
	return rl.Vector2{X: coord.X, Y: coord.Y}
}

// sortAlphabetically sorts alphabetically, including numbers.
// It is meant to be used inside functions like `sort.Sort`.
//
// Regular sort: "1", "10", "2" (see https://stackoverflow.com/a/35087122).
// This sort: 	 "1", "2", "10"
func sortAlphabetically(a, b string) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		iChar, bChar := a[i], b[j]

		// If both characters are a digit, sort by number
		if unicode.IsDigit(rune(iChar)) && unicode.IsDigit(rune(bChar)) {
			iStart, iEnd := extractNumber(a, i)
			jStart, jEnd := extractNumber(b, j)

			// Compare numbers as integers
			if iStart != jStart {
				return iStart < jStart
			}
			i, j = iEnd, jEnd
			continue
		}

		// Regular character comparison
		if iChar != bChar {
			return iChar < bChar
		}
		i++
		j++
	}

	// Shorter strings come first
	return len(a) < len(b)
}

func extractNumber(s string, start int) (int, int) {
	end := start
	for end < len(s) && unicode.IsDigit(rune(s[end])) {
		end++
	}
	n, _ := strconv.Atoi(s[start:end])
	return n, end
}

func calcBucket(value int, segment int) int {
	if value == 0 {
		return 0
	}
	return ((value-1)/segment)*segment + (segment - 1)
}

func all[T comparable](arr []T, v T) bool {
	for _, el := range arr {
		if el == v {
			continue
		}
		return false
	}
	return true
}

func parseCronField(field string, min int, max int) []int {
	list := []int{}

	for value := range strings.SplitSeq(field, ",") {
		// wildcard ("*")
		if value == "*" {
			for i := min; i <= max; i++ {
				list = append(list, i)
			}
			continue
		}

		// step (x/y)
		if parts := strings.Split(value, "/"); len(parts) == 2 {
			start := min
			if parts[0] != "*" {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					start = v
				} else {
					continue
				}
			}

			end, err1 := strconv.Atoi(parts[1])
			if err1 != nil {
				continue
			}

			for i := start; i <= max; i += end {
				list = append(list, i)
			}
			continue
		}

		// range (x-y)
		if parts := strings.Split(value, "-"); len(parts) == 2 {
			start, err0 := strconv.Atoi(parts[0])
			end, err1 := strconv.Atoi(parts[1])

			if err0 != nil || err1 != nil {
				continue
			}

			for i := start; i <= end; i++ {
				list = append(list, i)
			}
			continue
		}

		// literal (x)
		v, err := strconv.Atoi(value)
		if err == nil {
			list = append(list, v)
		}
	}
	return list
}

func stringsToCrons(crons map[string]string) []Cron {
	result := []Cron{}

	for name, cron := range crons {
		// "A B C D E" => "Min Hour Day Month Weekday"
		split := strings.Split(cron, " ")

		weekdays := parseCronField(split[4], 0, 6)
		hours := parseCronField(split[1], 0, 23)
		mins := parseCronField(split[0], 0, 59)

		for _, wd := range weekdays {
			for _, h := range hours {
				for _, m := range mins {
					result = append(result, Cron{
						Name:    name,
						Min:     m,
						Hour:    h,
						Weekday: wd,
					})
				}
			}
		}
	}

	return result
}

func cronsToCoords(crons []Cron) [][]Coord {
	result := make([][]Coord, 7)

	minuteSegment := float32(0)

	switch stepMin {
	case StepMin1:
		minuteSegment = 1
	case StepMin5:
		minuteSegment = 5
	case StepMin10:
		minuteSegment = 10
	case StepMin15:
		minuteSegment = 15
	case StepMin20:
		minuteSegment = 20
	case StepMin30:
		minuteSegment = 30
	}

	for _, cron := range crons {
		x := float32(cron.Hour)

		if groupBy == GroupByWdHourMin {
			bucket := calcBucket(cron.Min, int(minuteSegment))
			x += float32(bucket) / 60
		}

		result[cron.Weekday] = append(result[cron.Weekday], Coord{Name: cron.Name, X: x, Y: 1})
	}

	for wd, weekdays := range result {
		if len(weekdays) > 0 {
			weekdaysToggle[wd] = StatusOn
		} else {
			weekdaysToggle[wd] = StatusDisabled
		}
	}

	return result
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

// f <  0.5  -> darken color
// f == 0.5 -> same color
// f >  0.5  -> brighten color
func lerpRGB(r uint8, g uint8, b uint8, a uint8, f float32) (uint8, uint8, uint8, uint8) {
	f = max(0.0, min(1.0, f))

	r2 := uint8(0)
	g2 := uint8(0)
	b2 := uint8(0)
	a2 := uint8(0)

	if f < 0.5 {
		// Darken color
		factor := f / 0.5
		r2 = uint8(float32(r) * factor)
		g2 = uint8(float32(g) * factor)
		b2 = uint8(float32(b) * factor)
		a2 = uint8(float32(a) * factor)
	} else {
		// Brighten color
		factor := (f - 0.5) / 0.5
		r2 = uint8(float32(r) + (255-float32(r))*factor)
		g2 = uint8(float32(g) + (255-float32(g))*factor)
		b2 = uint8(float32(b) + (255-float32(b))*factor)
		a2 = uint8(float32(a) + (255-float32(a))*factor)
	}

	return r2, g2, b2, a2
}

func lerpColor(color rl.Color, f float32) rl.Color {
	r, g, b, a := lerpRGB(color.R, color.G, color.B, color.A, f)
	return rl.NewColor(r, g, b, a)
}

func lerpColorToHex(color rl.Color, f float32) rg.PropertyValue {
	return rg.NewColorPropertyValue(lerpColor(color, f))
}
