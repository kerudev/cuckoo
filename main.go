package main

import (
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	Name 	string
	Hour 	int
	Min 	int
}

type Coord struct {
	Name 	string
	X 		int
	Y 		int
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

var offset = rl.Vector2{ X: 20, Y: 20 }
var grid = Grid{}

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

func main() {
	sample := []string{
		"25 1,16,20 * * *",
		"24 7 * * *",
		"46 1,12 * * *",
	}

	crons := stringsToCrons(sample)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(800, 600, "Cuckoo")
	rl.SetWindowMinSize(800, 600)

	cell := Cell{W: 24, H: 12}
	grid.Cell = cell

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		grid.W = screenW - 40
		grid.H = screenH - 140

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
			for i := range 24 {
				fontSize := int32(12)
				text := strconv.Itoa(i)
				
				textW := rl.MeasureText(text, fontSize)
				textX := grid.W / cell.W * float32(i) - float32(textW / 2) + offset.X

				rl.DrawText(text, int32(textX), int32(grid.H + offset.Y + 2), fontSize, rl.Black)
			}

			// Draw coordinates
			for _, cron := range crons {
				rl.DrawCircle(
					int32(float32(cron.Hour) / cell.W * grid.W + offset.X),
					int32(grid.H + offset.X - 2),
					4,
					rl.Red,
				)
			}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
