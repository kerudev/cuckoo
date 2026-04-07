package cuckoo

import (
	"strconv"
	"strings"
)

func calcBucket(value int, segment int) int {
	if value == 0 {
		return 0
	}
	return ((value - 1) / segment) * segment + (segment - 1)
}

func stringsToCrons(crons []string) []Cron {
	result := []Cron{}

	for i, cron := range crons {
		split := strings.Split(cron, " ")

		for h := range strings.SplitSeq(split[1], ",") {
			hour, _ := strconv.Atoi(h)

			for m := range strings.SplitSeq(split[0], ",") {
				min, _ := strconv.Atoi(m)

				result = append(result, Cron{
					Hour: hour,
					Min:  min,
					Name: "process_" + strconv.Itoa(i),
				})
			}
		}
	}

	return result
}

func cronsToCoords(crons []Cron) []Coord {
	result := []Coord{}

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
			x += float32(bucket)/60
		}

		result = append(result, Coord{Name: cron.Name, X: x, Y: 1})
	}

	return result
}
