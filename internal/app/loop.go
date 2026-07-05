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

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 700, "Cuckoo")
	rl.SetWindowMinSize(800, 700)

	Font = rl.GetFontDefault()

	for !rl.WindowShouldClose() {
		S_Screen.Val.W = int32(rl.GetScreenWidth())
		S_Screen.Val.H = int32(rl.GetScreenHeight())

		S_Mouse.Set(rl.GetMousePosition())

		if !S_IsMouseLocked.Val {
			S_MouseWithLock.Set(S_Mouse.Val)
		}

		// Recalculate Grid and coordinates only when Screen changes size
		if S_Screen.HasChanged() {
			Grid.Width = S_Screen.Val.W - Offset.X*2
			Grid.Height = S_Screen.Val.H - Offset.Y*2 - 200

			S_IsMouseLocked.Set(false)

			gridCoords = CoordToGrid(coords)
		}

		// Check if a file was dropped and reload coords
		if rl.IsFileDropped() {
			droppedFiles := rl.LoadDroppedFiles()

			sample = map[string]string{}
			err := ReadPath(droppedFiles[0], &sample)

			if err != nil {
				fmt.Println(err)
			} else {
				crons = CronsFromStrings(sample)
				coords = CoordsFromCrons(crons)
				gridCoords = CoordToGrid(coords)
			}
		}

		handleKeyEvents()
		handleMouseEvents()
		handleMixedEvents()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.DrawGrid(gridCoords)

		ui.DrawUIOptions()

		// Vertical line
		lineY := Grid.Height + Offset.Y*7 - TextPad/2
		rl.DrawLine(Offset.X, lineY, 290, lineY, rl.Gray)

		ui.DrawUserOptions()

		// Horizontal line
		lineX := 150 + Offset.X + BoxPad*6
		rl.DrawLine(lineX, Grid.Height+Offset.Y*2+TextPad, lineX, S_Screen.Val.H-Offset.Y, rl.Gray)

		ui.DrawFooter()
		ui.DrawTooltip(gridCoords)

		rl.EndDrawing()

		// Recalculate coordinates based on bucket
		if S_StepMin.HasChanged() {
			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords)
		}

		// Recalculate coordinates based on group by
		if S_GroupBy.HasChanged() {
			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords)
		}

		// Reset zoom and coordinates
		if S_Zoom.HasChanged() || S_ZoomSlider.HasChanged() && S_Zoom.Val > 1 {
			// NOTE: unlock Mouse as MouseOver as coordinates are recalculated
			// when Zoom changes. Might be good to change this at some point.
			S_IsMouseLocked.Set(false)

			C_Zoom.Offset = S_ZoomSlider.Val * (C_Zoom.Scale - 1)

			coords = CoordsFromCrons(crons)
			gridCoords = CoordToGrid(coords)
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
	isOverGrid := rl.CheckCollisionPointRec(S_Mouse.Val, Grid.ToFloat32())
	isOverTooltip := rl.CheckCollisionPointRec(S_Mouse.Val, Tooltip.ToFloat32())

	if isOverTooltip && S_IsMouseLocked.Val {
		return
	}

	// Lock mouse position when clicking coordinates
	if TotalOver > 0 && isOverGrid && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		S_IsMouseLocked.Set(!S_IsMouseLocked.Val)
	}

	// Move zoom slider by dragging over grid
	if S_Zoom.Val > 1 && isOverGrid && rl.IsMouseButtonDown(rl.MouseButtonRight) {
		mouseX := rl.GetMouseDelta().X

		if mouseX != 0 {
			S_ZoomSlider.Set(Clamp(
				S_ZoomSlider.Val-mouseX/(C_Zoom.Scale-1),
				0,
				float32(Grid.Width),
			))
		}
	}
}

func handleMixedEvents() {
	// Move zoom slider with mouse and key events
	isOverGrid := rl.CheckCollisionPointRec(S_Mouse.Val, Grid.ToFloat32())

	if !isOverGrid {
		return
	}

	if !(TotalOver > 0 && S_IsMouseLocked.Val) {
		scroll := rl.GetMouseWheelMove()

		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Move zoom slider (horizontal scroll)
			calc := Cell.W / (C_Zoom.Scale * C_Zoom.Factor * 2)

			if scroll > 0 {
				S_ZoomSlider.Val += calc
			} else if scroll < 0 {
				S_ZoomSlider.Val -= calc
			}
		} else {
			// Zoom in (vertical scroll)
			S_Zoom.Set(Clamp(S_Zoom.Val+scroll, 1, 9))
			C_Zoom.Base = float32(Grid.Width) / float32(C_Grid.Cols)

			C_Zoom.Factor = (S_Zoom.Val - 1) / 8.0
			C_Zoom.Scale = float32(math.Pow(float64(Grid.Width)/float64(C_Zoom.Base), float64(C_Zoom.Factor)))

			C_Zoom.Offset = S_ZoomSlider.Val * (C_Zoom.Scale - 1)
		}
	}
}
