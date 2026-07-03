package models

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"

	. "github.com/kerudev/cuckoo/internal/utils"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

//////////////////////////////
// Raylib-like types
//////////////////////////////

type Vector2Int32 struct {
	X int32
	Y int32
}

//////////////////////////////
// Data types
//////////////////////////////

type Cron struct {
	Name    string
	Min     string
	Hour    string
	Day     string
	Month   string
	Weekday string
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

	Weekdays := ParseCronField(c.Weekday, 0, 6)
	hours := ParseCronField(c.Hour, 0, 23)
	mins := ParseCronField(c.Min, 0, 59)

	s := c.String()

	for _, wd := range Weekdays {
		for _, h := range hours {
			for _, m := range mins {
				jobs = append(jobs, Job{
					Name:    c.Name,
					Cron:    s,
					Min:     m,
					Hour:    h,
					Weekday: wd,
				})

				WdCounts[wd].Jobs += 1
			}
		}

		WdCounts[wd].Crons += 1
	}

	return jobs
}

func ParseCronField(field string, min int, max int) []int {
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

			// Step is not something like */5
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

func CronsFromStrings(strings map[string]string) []Cron {
	result := []Cron{}

	for name, cron := range strings {
		result = append(result, NewCron(name, cron))
	}

	return result
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

func (j Job) AsTime() string {
	return fmt.Sprintf("%02d:%02d", j.Hour, j.Min)
}

func JobsFromCrons(crons []Cron) []Job {
	result := []Job{}

	WdCounts = [WEEKDAYS]CountsByWd{}

	for _, cron := range crons {
		result = append(result, cron.Jobs()...)
	}

	return result
}

//////////////////////////////
// GUI related
//////////////////////////////

type Coord struct {
	Job Job
	X   float32
	Y   float32
}

func (c Coord) GridCoord() GridCoord {
	return GridCoord{Jobs: []Job{c.Job}, X: c.X, Y: c.Y}
}

func CoordsFromCrons(crons []Cron) [][]Coord {
	jobs := JobsFromCrons(crons)
	return CoordsFromJobs(jobs)
}

func CoordsFromJobs(jobs []Job) [][]Coord {
	result := make([][]Coord, WEEKDAYS)

	minuteSegment := S_StepMin.Val.Factor()

	for _, job := range jobs {
		x := float32(job.Hour)

		if S_GroupBy.Eq(GroupByWdHourMin) {
			bucket := CalcBucket(job.Min, minuteSegment)
			x += float32(bucket) / 60
		}

		result[job.Weekday] = append(result[job.Weekday], Coord{Job: job, X: x, Y: 1})
	}

	for wd, coords := range result {
		if len(coords) <= 0 {
			S_Weekdays.Val[wd].Status = StatusDisabled
		}
	}

	return result
}

type GridCoord struct {
	Jobs []Job
	X    float32
	Y    float32
	// OrigX float32
	OrigY float32
}

func (c GridCoord) Vector2() rl.Vector2 {
	return rl.Vector2{X: c.X, Y: c.Y}
}

func CoordToGrid(coords [][]Coord) [][]GridCoord {
	result := make([][]GridCoord, WEEKDAYS)

	C_Grid.Rows = INITIAL_ROWS
	C_Grid.Cols = INITIAL_COLS

	C_Grid.HighestY = 0

	for wd, coordDay := range coords {
		for _, coord := range coordDay {
			found := false

			if coord.Y > C_Grid.HighestY {
				C_Grid.HighestY = coord.Y
			}

			for i := range result[wd] {
				if coord.X == result[wd][i].X {
					found = true
					result[wd][i].Jobs = append(result[wd][i].Jobs, coord.Job)
				}

				if len(result[wd][i].Jobs) >= C_Grid.Rows {
					C_Grid.Rows = len(result[wd][i].Jobs) + 2
				}
			}

			if !found {
				cg := coord.GridCoord()
				// cg.OrigX = coord.X

				result[wd] = append(result[wd], cg)
			}
		}
	}

	C_Grid.HighestRow = C_Grid.Rows

	// Remove the last column, as it makes no sense when grouping by hour
	if S_GroupBy.Eq(GroupByWdHour) {
		C_Grid.Cols -= 1
	}

	if C_Grid.HighestRow > ROWS_CAP {
		C_Grid.Rows = INITIAL_ROWS
	}

	Cell.W = float32(Grid.Width) / float32(C_Grid.Cols)
	Cell.H = float32(Grid.Height) / float32(C_Grid.Rows)

	C_Zoom.Base = Cell.W

	scaledW := float32(Grid.Width) * C_Zoom.Scale
	highestRowY := float32(Grid.Height) / float32(C_Grid.HighestRow)

	for wd := range WEEKDAYS {
		// Translate coordinates to Grid
		for i := range result[wd] {
			result[wd][i].OrigY = float32(len(result[wd][i].Jobs))

			if result[wd][i].OrigY > C_Grid.HighestY {
				C_Grid.HighestY = result[wd][i].OrigY
			}

			result[wd][i].X = (result[wd][i].X/float32(C_Grid.Cols))*scaledW + float32(Offset.X) - C_Zoom.Offset
			result[wd][i].Y = float32(Grid.Height+Offset.Y) - highestRowY*result[wd][i].OrigY
		}

		// Sort coordinates to draw them in order
		sort.Slice(result[wd], func(i, j int) bool {
			return result[wd][i].X < result[wd][j].X
		})
	}

	return result
}

type Rec[T Numeric] struct {
	W T
	H T
}

type GridRec struct {
	// Rectangle
	W int32
	H int32
	X int32
	Y int32
}

//////////////////////////////
// Enums
//////////////////////////////

type TooltipPosition int32

const (
	PositionGrid TooltipPosition = iota
	PositionCoord
	// PositionFooter
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

func (s StepMin) Factor() int {
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

// ////////////////////////////
// State & Context
// ////////////////////////////
type AnyState interface {
	Update()
}

var AllStates = []AnyState{}

// https://stackoverflow.com/a/71065353
type State[T comparable] struct {
	Val T
	Old T
}

func NewState[T comparable](initial T) *State[T] {
	s := &State[T]{Val: initial, Old: initial}
	AllStates = append(AllStates, s)
	return s
}

func (s *State[T]) Update() {
	s.Old = s.Val
}

func (s *State[T]) Set(val T) {
	s.Val = val
}

func (s *State[T]) HasChanged() bool {
	return s.Val != s.Old
}

func (s *State[T]) Eq(val T) bool {
	return s.Val == val
}

type GridContext struct {
	Rows       int
	Cols       int
	HighestRow int
	HighestY   float32
}

type ZoomContext struct {
	Offset float32
	Base   float32
	Factor float32
	Scale  float32
}

//////////////////////////////
// Utility data types
//////////////////////////////

type UserOptions struct {
	DrawCoords bool
	DrawLines  bool
	DrawGrid   bool
	DrawFade   bool
}

type Weekday struct {
	Status Status
	Color  rl.Color
	Faded  rl.Color
}

func NewWeekday(color rl.Color) Weekday {
	return Weekday{Status: StatusOn, Color: color, Faded: rl.Fade(color, 0)}
}

type CountsByWd struct {
	Crons int
	Jobs  int
}

type ToggleParams struct {
	Icon string
	Ptr  *bool
}
