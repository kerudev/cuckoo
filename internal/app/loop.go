package app

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	ui "github.com/kerudev/cuckoo/internal/ui"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawLoop(sample map[string]string) {
	crons := CronsFromStrings(sample)
	coords := CoordsFromCrons(crons)

	gridCoords := [][]GridCoord{}

	groupByScroll := int32(0)
	positionScroll := int32(0)

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 700, "Cuckoo")
	rl.SetWindowMinSize(800, 700)

	Font = rl.GetFontDefault()

	for !rl.WindowShouldClose() {
		S_Screen.Val.W = int32(rl.GetScreenWidth())
		S_Screen.Val.H = int32(rl.GetScreenHeight())

		S_Mouse.Set(rl.GetMousePosition())

		// Recalculate Grid and coordinates only when Screen changes size
		if S_Screen.HasChanged() {
			Grid.W = S_Screen.Val.W - 40
			Grid.H = S_Screen.Val.H - 240

			gridCoords = CoordToGrid(coords, &Grid)
		}

		// Check if a file was dropped and reload coords

		// TODO this sometimes crashes when the file is being dragged over the
		// window. This may be a problem of Go's bindings
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			sample = map[string]string{}
			err := ReadPath(droppedFiles[0], &sample)

			if err != nil {
				fmt.Println(err)
			} else {
				crons = CronsFromStrings(sample)
				coords = CoordsFromCrons(crons)
				gridCoords = CoordToGrid(coords, &Grid)
			}
		}

		// Input event handling (keyboard)
		if rl.IsKeyPressed(rl.KeyL) {
			IsMouseLocked.Set(!IsMouseLocked.Val)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.DrawGrid(gridCoords)

		ui.DrawUIOptions(&groupByScroll)

		// Vertical line
		lineY := Grid.H + Offset.Y*7 - TextPad/2
		rl.DrawLine(Offset.X, lineY, 290, lineY, rl.Gray)

		ui.DrawUserOptions(&positionScroll)

		// Horizontal line
		lineX := 150 + Offset.X + BoxPad*6
		rl.DrawLine(lineX, Grid.H+Offset.Y*2+TextPad, lineX, S_Screen.Val.H-Offset.Y, rl.Gray)

		ui.DrawFooter()
		ui.DrawTooltip(gridCoords)

		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if S_StepMin.HasChanged() {
			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords, &Grid)
		}

		// Recalculate coordinates based on group by
		if S_GroupBy.HasChanged() {
			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords, &Grid)
		}

		if S_Zoom.HasChanged() || S_ZoomSlider.HasChanged() && S_Zoom.Val > 1 {
			// NOTE: unlock Mouse as MouseOver as coordinates are recalculated
			// when Zoom changes. Might be good to change this at some point.
			IsMouseLocked.Set(false)

			ZoomOffset = S_ZoomSlider.Val * (ZoomScale - 1)

			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords, &Grid)
		}

		// Save each state for next frame
		for _, state := range AllStates {
			state.Update()
		}
	}

	rl.CloseWindow()
}
