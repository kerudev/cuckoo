package app

import (
	"math"
	"path/filepath"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	ui "github.com/kerudev/cuckoo/internal/ui"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawLoop(path string) {
	// Init cuckoo internals and state
	isDir, _ := IsDir(path)
	if path == "" {
		ErrorText = "Select a valid file to parse"
	} else {
		handleNewFile(path)
	}

	// If the current path doesn't exist, use the current directory
	if path == "" || ErrorText != "" {
		path = "."
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		ErrorText = err.Error()
	}

	S_FilePath.Set(absPath)

	if !isDir {
		S_FileName.Set(absPath)

		// Force update before loop so handleNewFile is not called again on the first frame
		S_FileName.Update()
	}

	// Init raylib
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun | rl.FlagMsaa4xHint)
	rl.InitWindow(800, 800, "cuckoo")
	rl.SetWindowMinSize(800, 800)
	rl.SetExitKey(rl.KeyNull)

	Font = rl.GetFontDefault()
	ListViewItemH = int32(rg.GetStyle(rg.LISTVIEW, rg.LIST_ITEMS_HEIGHT) + rg.GetStyle(rg.LISTVIEW, rg.BORDER_WIDTH)*2 + 1)

	for !rl.WindowShouldClose() {
		S_Screen.Val.W = int32(rl.GetScreenWidth())
		S_Screen.Val.H = int32(rl.GetScreenHeight())

		S_Mouse.Set(rl.GetMousePosition())
		S_IsOnWindow.Set(rl.IsCursorOnScreen())
		S_IsOverGrid.Set(rl.CheckCollisionPointRec(S_Mouse.Val, Grid.ToFloat32()))

		if !S_IsMouseLocked.Val {
			S_MouseWithLock.Set(S_Mouse.Val)
		}

		// Recalculate UI, Grid and coordinates only when Screen changes size
		if S_Screen.HasChanged() {
			Grid.Width = S_Screen.Val.W - Offset.X*2
			Grid.Height = S_Screen.Val.H - Grid.Y - Offset.Y - 200

			Footer.Y = Grid.Height + Grid.Y

			FileButton.Width = float32(S_Screen.Val.W-Offset.X*2) - BackButton.Width - float32(BoxPad) - float32(BoxSize)*2
			FileButtonText.Width = FileButton.Width

			ErrorBox.Width = float32(S_Screen.Val.W - Offset.X*2 - BoxPad)

			ErrorMessageText.Width = ErrorBox.Width

			LockButton.X = FileButton.X + FileButton.Width + float32(BoxSize)*0.5
			HelpButton.X = LockButton.X + LockButton.Width + float32(BoxSize)*0.5

			S_IsMouseLocked.Set(false)

			GridCoords = CoordToGrid(Coords)
		}

		// Drag&Drop: check if a file was dropped and reload coords
		if rl.IsFileDropped() {
			handleNewFile(rl.LoadDroppedFiles()[0])
		}

		// Show or hide help window
		if rl.IsKeyPressed(rl.KeyH) {
			ShowHelp = !ShowHelp
		}

		BlockUI = ShowHelp || S_FilePicker.Val || S_FileName.Val == ""

		if !BlockUI {
			handleKeyEvents()
			handleMouseEvents()
			handleMixedEvents()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.DrawGrid()
		ui.DrawFooter()
		ui.DrawTooltip()

		ui.DrawError()
		ui.DrawFilePicker()
		ui.DrawLockButton()
		ui.DrawHelpButton()

		if ShowHelp {
			ui.DrawHelp()
		}

		rl.EndDrawing()

		// Recalculate coordinates when the file picker gets closed and a new file is chosen
		if S_FileName.HasChanged() {
			handleNewFile(S_FileName.Val)
		}

		// Recalculate coordinates based on bucket
		if S_StepMin.HasChanged() {
			Coords = CoordsFromCrons(Crons)
			GridCoords = CoordToGrid(Coords)
		}

		// Recalculate coordinates based on group by
		if S_GroupBy.HasChanged() {
			Coords = CoordsFromCrons(Crons)
			GridCoords = CoordToGrid(Coords)
		}

		// Reset zoom and coordinates
		if S_Zoom.HasChanged() || S_ZoomSlider.HasChanged() && S_Zoom.Val > 1 {
			// NOTE: unlock Mouse as MouseOver as coordinates are recalculated
			// when Zoom changes. Might be good to change this at some point.
			S_IsMouseLocked.Set(false)

			C_Zoom.Offset = S_ZoomSlider.Val * (C_Zoom.Scale - 1)

			Coords = CoordsFromCrons(Crons)
			GridCoords = CoordToGrid(Coords)
		}

		// Reset tooltip scroll
		if S_IsMouseLocked.HasChanged() {
			S_TooltipScroll.Set(0)
			S_Mouse.Set(rl.Vector2{})
		}

		// Reset MouseOver when mouse goes out of the grid
		if !BlockUI && !S_IsMouseLocked.Val && S_IsOverGrid.HasChanged() && !S_IsOverGrid.Val {
			MouseOver = [WEEKDAYS][]GridCoord{}
			TotalOver = 0
		}

		// Save each state for next frame
		for _, state := range AllStates {
			state.Update()
		}
	}

	rl.CloseWindow()
}

func handleNewFile(path string) {
	isDir, err := IsDir(path)
	absPath, _ := filepath.Abs(path)

	if err != nil {
		ErrorText = "Path " + absPath + " doesn't exist"
		return
	}

	if isDir {
		ErrorText = "Path " + absPath + " is not a file"
		return
	}

	sample := map[string]string{}
	if err := ReadPath(absPath, &sample); err != nil {
		ErrorText = err.Error()
		return
	}

	Crons = CronsFromStrings(sample)
	Coords = CoordsFromCrons(Crons)
	GridCoords = CoordToGrid(Coords)

	ErrorText = ""

	for wd, dayCoords := range GridCoords {
		if len(dayCoords) <= 0 {
			S_Weekdays.Val[wd].Status = StatusDisabled
		} else {
			S_Weekdays.Val[wd].Status = StatusOn
		}
	}

	S_FilePath.Set(path)
	S_FileName.Set(path)
}

func handleKeyEvents() {
	// Lock or unlock coordinates
	if TotalOver > 0 && S_IsOverGrid.Val && rl.IsKeyPressed(rl.KeyL) {
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

	if key >= rl.KeyKp1 && key <= rl.KeyKp7 {
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
	isOverTooltip := rl.CheckCollisionPointRec(S_Mouse.Val, Tooltip.ToFloat32())

	if isOverTooltip && S_IsMouseLocked.Val {
		return
	}

	// Lock mouse position when clicking coordinates
	if TotalOver > 0 && S_IsOverGrid.Val && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		S_IsMouseLocked.Set(!S_IsMouseLocked.Val)
	}

	// Move zoom slider by dragging over grid
	if S_Zoom.Val > 1 && S_IsOverGrid.Val && rl.IsMouseButtonDown(rl.MouseButtonRight) {
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
	if !S_IsOverGrid.Val && !S_IsMouseLocked.Val && TotalOver == 0 {
		return
	}

	// Move zoom slider with mouse and key events
	scroll := rl.GetMouseWheelMove()

	// Move zoom slider (horizontal scroll)
	if rl.IsKeyDown(rl.KeyLeftShift) {
		calc := Cell.W / (C_Zoom.Scale * C_Zoom.Factor * 2)

		if scroll > 0 {
			S_ZoomSlider.Val += calc
		} else if scroll < 0 {
			S_ZoomSlider.Val -= calc
		}

		return
	}

	// Zoom in (vertical scroll)
	if scroll > 0 && S_Zoom.Val < MAX_ZOOM || scroll < 0 && S_Zoom.Val > MIN_ZOOM {
		S_Zoom.Set(Clamp(S_Zoom.Val+scroll, MIN_ZOOM, MAX_ZOOM))
		C_Zoom.Base = float32(Grid.Width) / float32(C_Grid.Cols)

		C_Zoom.Factor = (S_Zoom.Val - 1) / (MAX_ZOOM - 1)
		C_Zoom.Scale = float32(math.Pow(float64(Grid.Width)/float64(C_Zoom.Base), float64(C_Zoom.Factor)))

		C_Zoom.Offset = S_ZoomSlider.Val * (C_Zoom.Scale - 1)
	}
}
