package cuckoo

import (
	"sort"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	Name string
	Hour int
	Min  int
}

type Coord struct {
	Name string
	X    float32
	Y    float32
}

type Cell struct {
	W float32
	H float32
}

type Grid struct {
	W    float32
	H    float32
	Cell Cell
}

type GridCoord struct {
	Names []string
	X     float32
	Y     float32
}

type DrawMode int32

const (
	DrawNone DrawMode = iota
	DrawLines
	DrawBezier
)

type GroupBy int32

const (
	GroupByHourMin GroupBy = iota
	GroupByHour
	// GroupByMin
)

type BucketMin int32

const (
	BucketMin1 BucketMin = iota
	BucketMin5
	BucketMin10
	BucketMin15
	BucketMin20
	BucketMin30
)

var offset = rl.Vector2{X: 20, Y: 20}
var step = rl.Vector2{X: 0, Y: 0}
var grid = Grid{Cell: Cell{W: float32(nCols), H: float32(nRows)}}

var drawCoords = true
var drawMode = DrawLines
var bucketMin = BucketMin1
var groupBy = GroupByHourMin

var nCols = 24
var nRows = 10

func coordToVec2(coord GridCoord) rl.Vector2 {
	return rl.Vector2{X: coord.X, Y: coord.Y}
}

func coordToGrid(coords []Coord, grid *Grid) []GridCoord {
	result := []GridCoord{}

	nRows = 10

	for _, coord := range coords {
		found := false

		for i := range result {
			if coord.X == result[i].X {
				found = true
				result[i].Names = append(result[i].Names, coord.Name)
			}

			if len(result[i].Names) >= nRows {
				nRows = len(result[i].Names) + 2
			}
		}

		if !found {
			result = append(result, GridCoord{
				Names: []string{coord.Name},
				X:     coord.X,
				Y:     coord.Y,
			})
		}
	}

	grid.Cell.H = float32(nRows)

	step.X = grid.W / grid.Cell.W
	step.Y = grid.H / grid.Cell.H

	cell := grid.Cell

	for i := range result {
		result[i].X = result[i].X/cell.W*grid.W + offset.X
		result[i].Y = grid.H + offset.Y - (grid.H / cell.H * float32(len(result[i].Names)))
	}

	return result
}

func DrawLoop(sample []string) {
	crons := stringsToCrons(sample)
	coords := cronsToCoords(crons)

	gridCoords := []GridCoord{}

	boxRoundness := float32(0.2)
	boxSegments := int32(8)
	boxPadX := float32(16)

	groupByScroll := int32(0)

	// Previous state
	prevScreenW := int32(0)
	prevScreenH := int32(0)

	prevGroupBy := groupBy
	prevBucketMin := bucketMin

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(800, 600, "Cuckoo")
	rl.SetWindowMinSize(800, 600)

	fontSize := float32(12)
	font := rl.GetFontDefault()
	textH := rl.MeasureTextEx(font, "0", fontSize, 1).Y

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		cell := grid.Cell

		// Recalculate grid and coordinates only when screen changes size
		if screenH != float32(prevScreenH) && screenW != float32(prevScreenW) {
			grid.W = screenW - 40
			grid.H = screenH - 140

			gridCoords = coordToGrid(coords, &grid)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw lines vertically
		for col := range int(cell.W) - 1 {
			x := step.X*float32(col+1) + offset.X

			rl.DrawLineEx(
				rl.Vector2{X: x, Y: offset.Y},
				rl.Vector2{X: x, Y: grid.H + offset.Y},
				2,
				rl.LightGray,
			)
		}

		// Draw lines horizontally
		for row := range int(nRows) - 1 {
			y := step.Y*float32(row+1) + offset.Y

			rl.DrawLineEx(
				rl.Vector2{X: offset.X, Y: y},
				rl.Vector2{X: grid.W + offset.Y, Y: y},
				2,
				rl.LightGray,
			)
		}

		// Draw grid container
		rl.DrawRectangleLinesEx(
			rl.Rectangle{X: offset.X, Y: offset.Y, Width: grid.W, Height: grid.H},
			2,
			rl.Black,
		)

		// Draw text on X axis
		for i := range nCols {
			text := strconv.Itoa(i)

			textW := rl.MeasureTextEx(font, text, fontSize, 1).X
			textX := step.X*float32(i) - textW/2 + offset.X

			rl.DrawText(text, int32(textX), int32(grid.H+offset.Y+2), int32(fontSize), rl.Black)
		}

		// Draw text on Y axis
		for i := range nRows + 1 {
			text := strconv.Itoa(i)
			textY := -step.Y*float32(i) - textH/2 + offset.Y + grid.H

			rl.DrawText(text, int32(offset.X/2), int32(textY), int32(fontSize), rl.Black)
		}

		if drawMode != DrawNone {
			// Sort coordinates to draw line in order
			sort.Slice(gridCoords, func(i, j int) bool {
				return gridCoords[i].X < gridCoords[j].X
			})

			// Draw lines that connect coordinates
			for i := 0; i < len(gridCoords)-1; i++ {
				start := coordToVec2(gridCoords[i])
				end := coordToVec2(gridCoords[i+1])

				switch drawMode {
				case DrawLines:
					rl.DrawLineEx(start, end, 2, rl.Red)
				case DrawBezier:
					rl.DrawLineBezier(start, end, 2, rl.Red)
				}
			}
		}

		// Draw coordinates
		if drawCoords {
			for _, coord := range gridCoords {
				rl.DrawCircle(int32(coord.X), int32(coord.Y), 4, rl.Red)
			}
		}

		// Draw option - GroupBy
		rl.DrawText("Group by", int32(offset.X), int32(grid.H+offset.Y*2+2), 12, rl.Black)
		groupByIdx := int32(groupBy)
		groupByIdx = rg.ListView(
			rl.Rectangle{X: offset.X, Y: grid.H + offset.Y*3, Width: 100, Height: 63},
			"Hour+Minute;Hour",
			&groupByScroll,
			groupByIdx,
		)
		groupBy = GroupBy(groupByIdx)

		// Draw option - DrawMode
		rl.DrawText("Draw mode", int32(120+offset.X), int32(grid.H+offset.Y*2+2), 12, rl.Black)
		drawModeIdx := int32(drawMode)
		drawModeIdx = rg.ToggleGroup(
			rl.Rectangle{X: 120 + offset.X, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
			"#113#;#127#;#125#",
			drawModeIdx,
		)
		drawMode = DrawMode(drawModeIdx)

		// Draw option - BucketMin
		if groupBy == GroupByHourMin {
			rl.DrawText("Minute bucket", int32(120+offset.X), int32(grid.H+offset.Y*4+2), 12, rl.Black)
			bucketMinIdx := int32(bucketMin)
			bucketMinIdx = rg.ToggleGroup(
				rl.Rectangle{X: 120 + offset.X, Y: grid.H + offset.Y*5, Width: 20, Height: 20},
				"#113#;5;10;15;20;30",
				bucketMinIdx,
			)
			bucketMin = BucketMin(bucketMinIdx)
		}

		// Draw option - DrawCoords
		drawCoords = rg.CheckBox(
			rl.Rectangle{X: 220 + offset.X, Y: grid.H + offset.Y*3, Width: 20, Height: 20},
			"Draw coordinates",
			drawCoords,
		)

		if drawMode == DrawNone && !drawCoords {
			drawCoords = true
		}

		// Draw coordinate information on mouse hover
		for _, coord := range gridCoords {
			mouseOverCoord := rl.CheckCollisionPointCircle(rl.GetMousePosition(), coordToVec2(coord), 4)

			if mouseOverCoord {
				maxW := float32(0)

				// Calculate max text size
				for _, name := range coord.Names {
					textW := float32(rl.MeasureText(name, int32(fontSize)))

					if textW > maxW {
						maxW = textW
					}
				}

				rec := rl.Rectangle{
					X:      coord.X + 8,
					Y:      coord.Y,
					Width:  maxW + 16,
					Height: fontSize * float32(len(coord.Names)),
				}

				rl.DrawRectangleRounded(rec, boxRoundness, boxSegments, rl.White)
				rl.DrawRectangleRoundedLinesEx(rec, boxRoundness, boxSegments, 2, rl.Black)

				for i, name := range coord.Names {
					spacingY := float32(i) * fontSize

					rl.DrawText(
						name,
						int32(coord.X+boxPadX),
						int32(coord.Y+spacingY),
						int32(fontSize),
						rl.Black,
					)
				}
			}
		}
		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if prevBucketMin != bucketMin {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		// Recalculate coordinates based on group by
		if prevGroupBy != groupBy {
			coords = cronsToCoords(crons)
			gridCoords = coordToGrid(coords, &grid)
		}

		prevScreenW = int32(screenW)
		prevScreenH = int32(screenH)

		prevGroupBy = groupBy
		prevBucketMin = bucketMin
	}

	rl.CloseWindow()
}
