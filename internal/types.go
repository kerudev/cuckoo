package cuckoo

import (
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	Name    string
	Min     string
	Hour    string
	Day     string
	Month   string
	Weekday string
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

func NewCron(name string, cron string) Cron {
	// "A B C D E" => "Min Hour Day Month Weekday"
	split := strings.Split(cron, " ")

	return Cron{
		Name:    name,
		Min:     split[0],
		Hour:    split[1],
		Day:     split[2],
		Month:   split[3],
		Weekday: split[4],
	}
}

func (c Cron) String() string {
	// "A B C D E" => "Min Hour Day Month Weekday"
	return strings.Join([]string{
		c.Min,
		c.Hour,
		c.Day,
		c.Month,
		c.Weekday,
	}, " ")
}

func (c Cron) Jobs() []Job {
	jobs := []Job{}

	weekdays := parseCronField(c.Weekday, 0, 6)
	hours := parseCronField(c.Hour, 0, 23)
	mins := parseCronField(c.Min, 0, 59)

	s := c.String()

	for _, wd := range weekdays {
		for _, h := range hours {
			for _, m := range mins {
				jobs = append(jobs, Job{
					Name:    c.Name,
					Cron:    s,
					Min:     m,
					Hour:    h,
					Weekday: wd,
				})
			}
		}
	}

	return jobs
}

type Job struct {
	Name    string
	Cron    string
	Min     int // 0-59
	Hour    int // 0-23
	Day     int // 1-31
	Month   int // 1-12
	Weekday int // 0-6
}

type Coord struct {
	Job Job
	X   float32
	Y   float32
}

func (c Coord) GridCoord() GridCoord {
	return GridCoord{Jobs: []Job{c.Job}, X: c.X, Y: c.Y}
}

type GridCoord struct {
	Jobs []Job
	X    float32
	Y    float32
}

func (c GridCoord) Vec2() rl.Vector2 {
	return rl.Vector2{X: c.X, Y: c.Y}
}

type Cell struct {
	W float32
	H float32
}

type Grid struct {
	W        float32
	H        float32
	Rows     int
	Cols     int
	HighestY int
}

type DrawMode int32

const (
	DrawNone DrawMode = iota
	DrawLines
	DrawBezier
)

type GroupBy int32

const (
	GroupByWdHour GroupBy = iota
	GroupByWdHourMin
	// GroupByMin
)

type StepMin int32

const (
	StepMin1 StepMin = iota
	StepMin5
	StepMin10
	StepMin15
	StepMin20
	StepMin30
)

func (s StepMin) Int() int {
	switch s {
	case StepMin1:
		return 1
	case StepMin5:
		return 5
	case StepMin10:
		return 10
	case StepMin15:
		return 15
	case StepMin20:
		return 20
	case StepMin30:
		return 30
	}

	return 1
}

type Status int32

const (
	StatusDisabled Status = iota - 1
	StatusOff
	StatusOn
)

func (s Status) Bool() bool {
	return s == StatusOn
}

func StatusFromBool(b bool) Status {
	if b {
		return StatusOn
	}
	return StatusOff
}
