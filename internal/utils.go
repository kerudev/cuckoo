package cuckoo

import (
	"strconv"
	"strings"
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
	return ((value-1)/segment)*segment + (segment - 1)
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

	switch bucketMin {
	case BucketMin1:
		minuteSegment = 1
	case BucketMin5:
		minuteSegment = 5
	case BucketMin10:
		minuteSegment = 10
	case BucketMin15:
		minuteSegment = 15
	case BucketMin20:
		minuteSegment = 20
	case BucketMin30:
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

	return result
}

// f <  0.5  -> darken color
// f == 0.5 -> same color
// f >  0.5  -> brighten color
func LerpRGB(r uint8, g uint8, b uint8, a uint8, f float32) (uint8, uint8, uint8, uint8) {
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

func LerpColor(color rl.Color, f float32) rl.Color {
	r, g, b, a := LerpRGB(color.R, color.G, color.B, color.A, f)
	return rl.NewColor(r, g, b, a)
}

func LerpColorToHex(color rl.Color, f float32) rg.PropertyValue {
	return rg.NewColorPropertyValue(LerpColor(color, f))
}
