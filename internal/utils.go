package cuckoo

import (
	"strconv"
	"unicode"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	return (value/segment)*segment + (segment - 1)
}

func countDuplicates[T comparable](arr []T) map[T]int {
	res := map[T]int{}

	for _, item := range arr {
		_, ok := res[item]
		if !ok {
			res[item] = 1
		} else {
			res[item]++
		}
	}

	return res
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

func stringsToCrons(sample map[string]string) []Cron {
	result := []Cron{}

	for name, cron := range sample {
		result = append(result, NewCron(name, cron))
	}

	return result
}

func cronsToJobs(crons []Cron) []Job {
	result := []Job{}

	for _, cron := range crons {
		result = append(result, cron.Jobs()...)
	}

	return result
}

func jobsToCoords(jobs []Job) [][]Coord {
	result := make([][]Coord, 7)

	minuteSegment := stepMin.Int()

	for _, job := range jobs {
		x := float32(job.Hour)

		if groupBy == GroupByWdHourMin {
			bucket := calcBucket(job.Min, minuteSegment)
			x += float32(bucket) / 60
		}

		result[job.Weekday] = append(result[job.Weekday], Coord{Job: job, X: x, Y: 1})
	}

	return result
}

func cronsToCoords(crons []Cron) [][]Coord {
	jobs := cronsToJobs(crons)
	return jobsToCoords(jobs)
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
					result[day][i].Jobs = append(result[day][i].Jobs, coord.Job)
				}

				if len(result[day][i].Jobs) >= grid.Rows {
					grid.Rows = len(result[day][i].Jobs) + 2
				}
			}

			if !found {
				result[day] = append(result[day], coord.GridCoord())
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

	scaledW := grid.W * scale
	highestYPos := grid.H / float32(grid.HighestY)

	for day := range 7 {
		for i := range result[day] {
			result[day][i].X = (result[day][i].X/float32(grid.Cols))*scaledW + offset.X - zoomOffset
			result[day][i].Y = grid.H + offset.Y - (highestYPos * float32(len(result[day][i].Jobs)))
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
