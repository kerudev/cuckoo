package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(800, 600, "Go test")
	rl.SetWindowMinSize(800, 600)

	cells := 24

	for !rl.WindowShouldClose() {
		screenW := float32(rl.GetScreenWidth())
		screenH := float32(rl.GetScreenHeight())

		gridW := screenW - 40
		gridH := screenH - 140

		rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)

			// Draw lines vertically
			for col := range cells - 1 {
				x := gridW / float32(cells) * float32(col + 1) + 20

				rl.DrawLineEx(
					rl.Vector2{ X: x, Y: 20 },
					rl.Vector2{ X: x, Y: gridH + 20 },
					2,
					rl.LightGray,
				)
			}

			// Draw lines horizontally
			for row := range cells - 1 {
				y := gridH / float32(cells) * float32(row + 1) + 20

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
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
