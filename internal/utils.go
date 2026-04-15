package cuckoo

import (
	"strconv"
	"strings"
	"unicode"
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

func stringsToCrons(crons map[string]string) []Cron {
	result := []Cron{}

	for name, cron := range crons {
		// "A B C D E" => "Min Hour Day Month Weekday"
		split := strings.Split(cron, " ")

		for wd := range strings.SplitSeq(split[4], ",") {
			weekday, _ := strconv.Atoi(wd)

			for h := range strings.SplitSeq(split[1], ",") {
				hour, _ := strconv.Atoi(h)

				for m := range strings.SplitSeq(split[0], ",") {
					min, _ := strconv.Atoi(m)

					result = append(result, Cron{
						Name:       name,
						Min:        min,
						Hour:       hour,
						Weekday:    weekday,
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

		if groupBy == GroupByHourMin {
			bucket := calcBucket(cron.Min, int(minuteSegment))
			x += float32(bucket) / 60
		}

		result[cron.Weekday] = append(result[cron.Weekday], Coord{Name: cron.Name, X: x, Y: 1})
	}

	return result
}
