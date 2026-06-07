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

func all[T comparable](arr []T, pred func(T) bool) bool {
	for _, el := range arr {
		if pred(el) {
			continue
		}
		return false
	}
	return true
}

// minF32 returns the smaller of x or y.
func minF32(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
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

	minuteSegment := stepMin.Factor()

	for _, job := range jobs {
		x := float32(job.Hour)

		if groupBy == GroupByWdHourMin {
			bucket := calcBucket(job.Min, minuteSegment)
			x += float32(bucket) / 60
		}

		result[job.Weekday] = append(result[job.Weekday], Coord{Job: job, X: x, Y: 1})
	}

	for wd, coords := range result {
		if len(coords) > 0 {
			weekdays[wd].status = StatusOn
		} else {
			weekdays[wd].status = StatusDisabled
		}
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

	gridHighestY = 0

	for wd, coordDay := range coords {
		for _, coord := range coordDay {
			found := false

			if coord.Y > gridHighestY {
				gridHighestY = coord.Y
			}

			for i := range result[wd] {
				if coord.X == result[wd][i].X {
					found = true
					result[wd][i].Jobs = append(result[wd][i].Jobs, coord.Job)
				}

				if len(result[wd][i].Jobs) >= grid.Rows {
					grid.Rows = len(result[wd][i].Jobs) + 2
				}
			}

			if !found {
				cg := coord.GridCoord()
				// cg.OrigX = coord.X

				result[wd] = append(result[wd], cg)
			}
		}
	}

	grid.HighestRow = grid.Rows

	// Remove the last column, as it makes no sense when grouping by hour
	if groupBy == GroupByWdHour {
		grid.Cols -= 1
	}

	if grid.HighestRow > ROWS_CAP {
		grid.Rows = INITIAL_ROWS
	}

	cell.W = float32(grid.W) / float32(grid.Cols)
	cell.H = float32(grid.H) / float32(grid.Rows)

	scaledW := float32(grid.W) * zoomScale
	highestRowY := float32(grid.H) / float32(grid.HighestRow)

	for wd := range 7 {
		for i := range result[wd] {
			result[wd][i].OrigY = float32(len(result[wd][i].Jobs))

			if result[wd][i].OrigY > gridHighestY {
				gridHighestY = result[wd][i].OrigY
			}

			result[wd][i].X = (result[wd][i].X/float32(grid.Cols))*scaledW + float32(offset.X) - zoomOffset
			result[wd][i].Y = float32(grid.H+offset.Y) - highestRowY*result[wd][i].OrigY
		}
	}

	return result
}

// f <  0.5  -> darken color
// f == 0.5 -> same color
// f >  0.5  -> brighten color
func lerpRGB(r uint8, g uint8, b uint8, a uint8, f float32) (uint8, uint8, uint8, uint8) {
	f = max(0.0, min(1.0, f))

	rf := float32(r)
	gf := float32(g)
	bf := float32(b)
	af := float32(a)

	r2 := uint8(0)
	g2 := uint8(0)
	b2 := uint8(0)
	a2 := uint8(0)

	if f < 0.5 {
		// Darken color
		factor := f / 0.5
		r2 = uint8(rf * factor)
		g2 = uint8(gf * factor)
		b2 = uint8(bf * factor)
		a2 = uint8(af * factor)
	} else {
		// Brighten color
		factor := (f - 0.5) / 0.5
		r2 = uint8(rf + (255-rf)*factor)
		g2 = uint8(gf + (255-gf)*factor)
		b2 = uint8(bf + (255-bf)*factor)
		a2 = uint8(af + (255-af)*factor)
	}

	return r2, g2, b2, a2
}

func lerpColor(color rl.Color, f float32) rl.Color {
	return rl.NewColor(lerpRGB(color.R, color.G, color.B, color.A, f))
}

func lerpColorToHex(color rl.Color, f float32) rg.PropertyValue {
	return rg.NewColorPropertyValue(lerpColor(color, f))
}
