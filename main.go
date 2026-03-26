package main

import (
	// "fmt"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	crons := []string{
		"25 1,16,20 * * *",
		"24 7 * * *",
		"46 1,12 * * *",
	}

	colors := []rl.Color{
		{R: 99, G: 52, B: 180, A: 255},
		{R: 65, G: 147, B: 138, A: 255},
		{R: 254, G: 3, B: 93, A: 255},
	}

	coords := [][]string{}

	for _, cron := range crons {
		split := strings.Split(cron, " ")
		hours := split[1]

		coords = append(coords, strings.Split(hours, ","))
	}

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
				x := gridW / cellsX * float32(col + 1) + 20

				rl.DrawLineEx(
					rl.Vector2{ X: x, Y: 20 },
					rl.Vector2{ X: x, Y: gridH + 20 },
					2,
					rl.LightGray,
				)
			}

			// Draw lines horizontally
			for row := range int(cellsY) - 1 {
				y := gridH / cellsY * float32(row + 1) + 20

				rl.DrawLineEx(
					rl.Vector2{ X: 20, Y: y },
					rl.Vector2{ X: gridW + 20, Y: y },
					2,
					rl.LightGray,
				)
			}

			// Draw grid container
			rl.DrawRectangleLinesEx(
				rl.Rectangle{ X: 20, Y: 20, Width: gridW, Height: gridH },
				2,
				rl.Black,
			)

			for i := range 24 {
				x := gridW / cellsX * float32(i) + 20
				rl.DrawText(strconv.Itoa(i), int32(x), int32(gridH + 22), 12, rl.Black)
			}

			for color, coord := range coords {
				for _, hour := range coord {
					
					hour_f, _ := strconv.ParseFloat(hour, 32)
					// fmt.Printf("Hour: %f\n", hour_f)

					rl.DrawCircle(
						int32(float32(hour_f) / cellsX * gridW + 20),
						int32(gridH + 20),
						4,
						colors[color],
					)
				}
			}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
