package main

import (
	"sort"
	"strconv"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	Name 	string
	Hour 	int
	Min 	int
}

type Coord struct {
	Name 	string
	X 		float32
	Y 		float32
}

type Cell struct {
	W		float32
	H		float32	
}

type Grid struct {
	W		float32
	H		float32
	Cell	Cell
}

type GridCoord struct {
	Names []string
	X 		float32
	Y 		float32
}

var offset = rl.Vector2{ X: 20, Y: 20 }
var grid = Grid{}

var drawMode = int32(0)

func stringsToCrons(crons []string) []Cron {
	result := []Cron{}

	for i, cron := range crons {
		split := strings.Split(cron, " ")

		mins := strings.Split(split[0], ",")
		hours := strings.Split(split[1], ",")

		for _, h := range hours {
			hour, _ := strconv.Atoi(h)

			for _, m := range mins {
				min, _ := strconv.Atoi(m)
				result = append(result, Cron{
					Hour: hour,
					Min: min,
					Name: "process" + strconv.Itoa(i),
				})
			}
		}
	}

	return result
}

func cronsToCoords(crons []Cron) []Coord {
	result := []Coord{}

	for _, cron := range crons {
		result = append(result, Coord{
			Name: cron.Name,
			X: float32(cron.Hour) + float32(cron.Min) / 60,
			Y: 1,
		})
	}

	return result
}

func coordToGrid(coords []Coord, grid Grid) []GridCoord {
	result := []GridCoord{}

	for _, coord := range coords {
		found := false

		for i := range result {
			if coord.X == result[i].X {
				found = true
				result[i].Names = append(result[i].Names, coord.Name)
			}
		}

		if !found {
			result = append(result, GridCoord{
				Names: []string{coord.Name},
				X: coord.X,
				Y: coord.Y,
			})
		}
	}

	cell := grid.Cell

	for i := range result {
		result[i].X = float32(result[i].X) / cell.W * grid.W + offset.X
		result[i].Y = grid.H + offset.Y - (grid.H / cell.H * float32(len(result[i].Names)))
	}

	return result
}

func main() {
	sample := []string{
		"0 1,16,20 * * *",
		"25 1,16,20 * * *",
		"25 1,16,20 * * *",
		"24 7 * * *",
		"46 1,12 * * *",
	}

	crons := stringsToCrons(sample)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(800, 600, "Cuckoo")
	rl.SetWindowMinSize(800, 600)

	nCols := 24
	nRows := 10

	cell := Cell{W: float32(nCols), H: float32(nRows)}
	grid.Cell = cell

	font := rl.GetFontDefault()
	fontSize := float32(12)

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		grid.W = screenW - 40
		grid.H = screenH - 140

		coords := cronsToCoords(crons)
		gridCoords := coordToGrid(coords, grid)

		rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)

			// Draw lines vertically
			for col := range int(cell.W) - 1 {
				x := grid.W / cell.W * float32(col + 1) + offset.X

				rl.DrawLineEx(
					rl.Vector2{ X: x, Y: offset.Y },
					rl.Vector2{ X: x, Y: grid.H + offset.Y },
					2,
					rl.LightGray,
				)
			}

			// Draw lines horizontally
			for row := range int(cell.H) - 1 {
				y := grid.H / cell.H * float32(row + 1) + offset.Y

				rl.DrawLineEx(
					rl.Vector2{ X: offset.X, Y: y },
					rl.Vector2{ X: grid.W + offset.Y, Y: y },
					2,
					rl.LightGray,
				)
			}

			// Draw grid container
			rl.DrawRectangleLinesEx(
				rl.Rectangle{ X: offset.X, Y: offset.Y, Width: grid.W, Height: grid.H },
				2,
				rl.Black,
			)

			// Draw text on X axis
			for i := range nCols {
				text := strconv.Itoa(i)

				textW := rl.MeasureTextEx(font, text, fontSize, 1).X
				textX := grid.W / cell.W * float32(i) - float32(textW / 2) + offset.X

				rl.DrawText(text, int32(textX), int32(grid.H + offset.Y + 2), int32(fontSize), rl.Black)
			}

			// Draw text on Y axis
			for i := range nRows + 1 {
				text := strconv.Itoa(i)

				textH := rl.MeasureTextEx(font, text, fontSize, 1).Y
				textY := -grid.H / cell.H * float32(i) - float32(textH / 2) + offset.Y + grid.H

				rl.DrawText(text, int32(offset.X / 2), int32(textY), int32(fontSize), rl.Black)
			}

			// Draw coordinates
			for _, coord := range gridCoords {
				rl.DrawCircle(int32(coord.X), int32(coord.Y), 4, rl.Red)
			}

			rl.DrawText("Draw mode", int32(offset.X), int32(grid.H + offset.Y * 2 + 2), 12, rl.Black)
			drawMode = rg.ToggleGroup(
				rl.Rectangle{ X: offset.X, Y: grid.H + offset.Y * 3, Width: 20, Height: 20 },
				"#113#;#127#;#125#",
				drawMode,
			)

			if drawMode != 0 {
				// Sort coordinates to draw line in order
				sort.Slice(gridCoords, func(i, j int) bool {
					return gridCoords[i].X < gridCoords[j].X
				})
	
				for i := 0; i < len(gridCoords) - 1; i++ {
					start := rl.Vector2{X: float32(gridCoords[i].X), Y: float32(gridCoords[i].Y)}
					end := rl.Vector2{X: float32(gridCoords[i + 1].X), Y: float32(gridCoords[i + 1].Y)}
					
					if drawMode == 1 {
						rl.DrawLineEx(start, end, 2, rl.Red)
					} else if drawMode == 2 {
						rl.DrawLineBezier(start, end, 2, rl.Red)
					}
				}
			}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
