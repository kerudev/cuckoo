package cuckoo

type Cron struct {
	Name    string
	Min     int // 0-59
	Hour    int // 0-23
	Day     int // 1-31
	Month   int // 1-12
	Weekday int // 0-6
}

type Coord struct {
	Name string
	X    float32
	Y    float32
}

type GridCoord struct {
	Names []string
	X     float32
	Y     float32
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
