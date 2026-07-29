package app

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	ui "github.com/kerudev/cuckoo/internal/ui"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

var sample = map[string]string{}
var crons = []Cron{}
var coords = [][]Coord{}
var gridCoords = [][]GridCoord{}

func DrawLoop(path string) {
	// Init cuckoo internals and state
	handleNewFile(path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
	}

	S_FileName.Set(absPath)
	S_LastFile.Set(absPath)

	// Init raylib
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 800, "Cuckoo")
	rl.SetWindowMinSize(800, 800)
	rl.SetExitKey(rl.KeyNull)

	Font = rl.GetFontDefault()
	ListViewItemH = int32(rg.GetStyle(rg.LISTVIEW, rg.LIST_ITEMS_HEIGHT) + rg.GetStyle(rg.LISTVIEW, rg.BORDER_WIDTH)*2 + 1)

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
			Grid.Height = S_Screen.Val.H - Grid.Y - Offset.Y - 200

			Footer.Y = Grid.Height + Grid.Y

			S_IsMouseLocked.Set(false)

			gridCoords = CoordToGrid(coords)
		}

		// Check if a file was dropped and reload coords
		if rl.IsFileDropped() {
			handleNewFile(rl.LoadDroppedFiles()[0])
		}

		handleKeyEvents()
		handleMouseEvents()
		handleMixedEvents()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.DrawGrid(gridCoords)

		ui.DrawFooter()
		ui.DrawTooltip(gridCoords)

		ui.DrawFilePicker()

		if ShowHelp {
			ui.DrawHelp()
		}

		rl.EndDrawing()

		// Recalculate coordinates when the file picker gets closed and a new file is chosen
		// if S_FileName.HasChanged() && S_FilePicker.HasChanged() && S_FilePicker.Val {
		if S_FileName.HasChanged() {
			handleNewFile(S_FileName.Val)
		}

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
			S_Mouse.Set(rl.Vector2{})
		}

		// Save each state for next frame
		for _, state := range AllStates {
			state.Update()
		}
	}

	rl.CloseWindow()
}

func handleNewFile(path string) {
	stat, _ := os.Stat(path)
	if stat.IsDir() {
		return
	}

	sample = map[string]string{}
	err := ReadPath(path, &sample)

	if err != nil {
		fmt.Println(err)
	} else {
		crons = CronsFromStrings(sample)
		coords = CoordsFromCrons(crons)
		gridCoords = CoordToGrid(coords)
	}

	for wd, dayCoords := range gridCoords {
		if len(dayCoords) <= 0 {
			S_Weekdays.Val[wd].Status = StatusDisabled
		} else {
			S_Weekdays.Val[wd].Status = StatusOn
		}
	}
}

func handleKeyEvents() {
	// Show or hide help window
	if rl.IsKeyPressed(rl.KeyH) {
		ShowHelp = !ShowHelp

		if ShowHelp {
			rg.Lock()
		} else {
			rg.Unlock()
		}
	}

	if rg.IsLocked() {
		return
	}

	// Lock or unlock coordinates
	if rl.IsKeyPressed(rl.KeyL) {
		S_IsMouseLocked.Set(!S_IsMouseLocked.Val)
	}

	// Handle number keys to change Weekdays status
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

	wd := int(key % mod)

	if S_Weekdays.Val[wd].Status != StatusDisabled {
		S_Weekdays.Val.SetStatus(wd, !S_Weekdays.Val[wd].Status.Bool())
	}
}

func handleMouseEvents() {
	if rg.IsLocked() {
		return
	}

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
	if rg.IsLocked() {
		return
	}

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
