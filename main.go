package main

import (
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cron struct {
	name 	string
	hour 	int
	min 	int
}

var offset = rl.Vector2{ X: 20, Y: 20 }

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
					hour: hour,
					min: min,
					name: "process" + strconv.Itoa(i),
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

	rl.InitWindow(800, 600, "Go test")
	rl.SetWindowMinSize(800, 600)

	cellsX := float32(24)
	cellsY := float32(8)

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		gridW := screenW - 40
		gridH := screenH - 140

		rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)

			// Draw lines vertically
			for col := range int(cellsX) - 1 {
				x := gridW / cellsX * float32(col + 1) + offset.X

				rl.DrawLineEx(
					rl.Vector2{ X: x, Y: offset.Y },
					rl.Vector2{ X: x, Y: gridH + offset.Y },
					2,
					rl.LightGray,
				)
			}

			// Draw lines horizontally
			for row := range int(cellsY) - 1 {
				y := gridH / cellsY * float32(row + 1) + offset.Y

				rl.DrawLineEx(
					rl.Vector2{ X: offset.X, Y: y },
					rl.Vector2{ X: gridW + offset.Y, Y: y },
					2,
					rl.LightGray,
				)
			}

			// Draw grid container
			rl.DrawRectangleLinesEx(
				rl.Rectangle{ X: offset.X, Y: offset.Y, Width: gridW, Height: gridH },
				2,
				rl.Black,
			)

			// Draw text on X axis
			for i := range 24 {
				fontSize := int32(12)
				text := strconv.Itoa(i)
				
				textW := rl.MeasureText(text, fontSize)
				textX := gridW / cellsX * float32(i) - float32(textW / 2) + offset.X

				rl.DrawText(text, int32(textX), int32(gridH + 22), fontSize, rl.Black)
			}

			// Draw coordinates
			for _, cron := range crons {
				rl.DrawCircle(
					int32(float32(cron.hour) / cellsX * gridW + offset.X),
					int32(gridH + 18),
					4,
					rl.Red,
				)
			}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
