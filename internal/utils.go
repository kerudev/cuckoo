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

func parseCronField(field string) []int {
	list := []int{}
	for value := range strings.SplitSeq(field, ",") {
		v, err := strconv.Atoi(value)

		// If v is not a literal (n), it will fail parsing
		if err == nil {
			list = append(list, v)
			continue
		}

		// Loop over range (n-m) values
		rng := strings.Split(value, "-")
		rng0, _ := strconv.Atoi(rng[0])
		rng1, _ := strconv.Atoi(rng[1])

		for i := rng0; i <= rng1; i++ {
			list = append(list, i)
		}
	}
	return list
}

func stringsToCrons(crons map[string]string) []Cron {
	result := []Cron{}

	for name, cron := range crons {
		// "A B C D E" => "Min Hour Day Month Weekday"
		split := strings.Split(cron, " ")

		weekdays := parseCronField(split[4])
		hours := parseCronField(split[1])
		mins := parseCronField(split[0])

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

		if groupBy == GroupByHourMin {
			bucket := calcBucket(cron.Min, int(minuteSegment))
			x += float32(bucket) / 60
		}

		result[cron.Weekday] = append(result[cron.Weekday], Coord{Name: cron.Name, X: x, Y: 1})
	}

	return result
}
