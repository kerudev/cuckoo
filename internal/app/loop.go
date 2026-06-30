package app

import (
	"fmt"
	"math"

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
			Grid.W = S_Screen.Val.W - Offset.X*2
			Grid.H = S_Screen.Val.H - Offset.Y*2 - 200

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

		handleKeyEvents()
		handleMouseEvents()
		handleMixedEvents()

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

		// Reset zoom and coordinates
		if S_Zoom.HasChanged() || S_ZoomSlider.HasChanged() && S_Zoom.Val > 1 {
			// NOTE: unlock Mouse as MouseOver as coordinates are recalculated
			// when Zoom changes. Might be good to change this at some point.
			S_IsMouseLocked.Set(false)

			ZoomOffset = S_ZoomSlider.Val * (ZoomScale - 1)

			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords, &Grid)
		}

		// Reset tooltip scroll
		if S_IsMouseLocked.HasChanged() {
			S_TooltipScroll.Set(0)
		}

		// Save each state for next frame
		for _, state := range AllStates {
			state.Update()
		}
	}

	rl.CloseWindow()
}

func handleKeyEvents() {
	// Input event handling (keyboard)
	if rl.IsKeyPressed(rl.KeyL) {
		S_IsMouseLocked.Set(!S_IsMouseLocked.Val)
	}

	key := rl.GetKeyPressed()

	// Return if no key was pressed
	if key == rl.KeyNull {
		return
	}

	mod := int32(rl.KeyNull)

	if key >= rl.KeyOne && key <= rl.KeySeven {
		mod = rl.KeyOne
	}

	if key >= rl.KeyKp1 && key <= rl.KeyKp8 {
		mod = rl.KeyKp1
	}

	// Return if the key pressed is not a number
	if mod == rl.KeyNull {
		return
	}

	idx := key % mod

	if S_Weekdays.Val[idx].Status != StatusDisabled {
		active := S_Weekdays.Val[idx].Status.Bool()
		S_Weekdays.Val[idx].Status = StatusFromBool(!active)

		if active && All(S_Weekdays.Val[:], func(wd Weekday) bool { return wd.Status != StatusOn }) {
			S_Weekdays.Val[idx].Status = StatusOn
		}
	}
}

func handleMouseEvents() {
	isOverTooltip := rl.CheckCollisionPointRec(S_Mouse.Val, Tooltip.ToFloat32())

	if isOverTooltip && S_IsMouseLocked.Val {
		return
	}

	// Lock mouse position when clicking coordinates
	if TotalOver > 0 && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		S_IsMouseLocked.Set(!S_IsMouseLocked.Val)
	}

	// Move zoom slider by dragging over grid
	if S_Zoom.Val > 1 && rl.IsMouseButtonDown(rl.MouseButtonRight) {
		mouseX := rl.GetMouseDelta().X

		if mouseX != 0 {
			S_ZoomSlider.Set(Clamp(
				S_ZoomSlider.Val-mouseX/(ZoomScale-1),
				0,
				float32(Grid.W),
			))
		}
	}
}

func handleMixedEvents() {
	// Move zoom slider with mouse and key events
	isOverTooltip := rl.CheckCollisionPointRec(S_Mouse.Val, Tooltip.ToFloat32())

	if !(TotalOver > 0 && S_IsMouseLocked.Val && isOverTooltip) {
		scroll := rl.GetMouseWheelMove()

		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Move zoom slider (horizontal scroll)
			calc := Cell.W / (ZoomScale * ZoomFactor * 2)

			if scroll > 0 {
				S_ZoomSlider.Val += calc
			} else if scroll < 0 {
				S_ZoomSlider.Val -= calc
			}
		} else {
			// Zoom in (vertical scroll)
			S_Zoom.Set(Clamp(S_Zoom.Val+scroll, 1, 9))
			ZoomBase = float32(Grid.W) / float32(Grid.Cols)

			ZoomFactor = (S_Zoom.Val - 1) / 8.0
			ZoomScale = float32(math.Pow(float64(Grid.W)/float64(ZoomBase), float64(ZoomFactor)))

			ZoomOffset = S_ZoomSlider.Val * (ZoomScale - 1)
		}
	}
}
