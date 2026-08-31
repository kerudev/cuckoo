package ui

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/shared"
)

var HelpLines = []string{
	" ", // Reserved to LabelButton: "Open issue"

	"-",

	"Drag&Drop file > Changes grid coordinates",
	"#138# > Locked state indicator. Click to unlock",
	"#193#/H > Toggle help window",
	"D > Toggle debug information (just FPS count for now)",

	"-File Picker",
	"#118# > Go back by one directory",
	"Click on any file to change the contents on the grid",

	"-Grid",
	"LClick/L > Lock coordinates where mouse is over",
	"Wheel[U/D] > Zoom in/out",
	"LShift + Wheel[U/D] > Scroll right/left",
	"(Zoomed) Hold RMouse > Scroll right/left",
	"(Zoomed) LClick & Drag Scrollbar > Scroll right/left",

	"-Tooltip",
	"Wheel[U/D] > Scroll up/down",

	"-UI options",
	"Group by > Groups coordinates by hour or hour+minute (default)",
	"1-7 (numbers/keypad) > Toggle weekdays (0-6)",
	"Minute group > Group coordinates by N minutes",

	"-User options",
	"Tooltip position > Draw on the grid (default) or over coordinates",
	"Draw options:",
	"#212# Draw coordinates",
	"#127# Draw lines that join coordinates",
	"#94# Draw fade under coordinates",
	"#97# Draw grid lines",
}

func DrawHelpButton() {
	if ShowHelp {
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
	}

	// Icon: HELP
	if rg.Button(HelpButton, "#193#") {
		// Set to true as the button can only be pressed when the state is STATE_NORMAL
		ShowHelp = true
	}
}

func DrawHelp() {
	// Inspired from: https://github.com/raysan5/rguistyler/blob/3be4d40/src/gui_window_help.h

	rl.DrawRectangle(0, 0, S_Screen.Val.W, S_Screen.Val.H, rl.Fade(rl.White, 0.85))

	if HelpWindow.Height == 0 || S_Screen.HasChanged() {
		// Line + margin
		HelpWindow.Height = ModalLineH + 8

		for _, line := range HelpLines {
			if line == "-" {
				HelpWindow.Height += ModalLineEmptyH
			} else {
				HelpWindow.Height += ModalLineH
			}
		}

		HelpWindow.Width = 360
		HelpWindow.X = (S_Screen.Val.W - HelpWindow.Width) / 2
		HelpWindow.Y = (S_Screen.Val.H - HelpWindow.Height) / 2
	}

	lineY := ModalLineH + 4

	// Icon: HELP
	ShowHelp = !rg.WindowBox(HelpWindow.ToFloat32(), "#193# Help and user guide")

	// Icon: WARNING
	if rg.LabelButton(
		NewRectangleFromInt32(HelpWindow.X+12, HelpWindow.Y+lineY, HelpWindow.Width, ModalLineH),
		"#220# Found a bug or have an idea? Click to open an issue!",
	) {
		rl.OpenURL("https://github.com/kerudev/cuckoo/issues/new")
	}

	for _, line := range HelpLines {
		if line == "-" {
			rg.Line(NewRectangleFromInt32(HelpWindow.X, HelpWindow.Y+lineY, HelpWindow.Width, ModalLineEmptyH), "")
		} else if line[0] == '-' {
			rg.Line(NewRectangleFromInt32(HelpWindow.X, HelpWindow.Y+lineY, HelpWindow.Width, ModalLineH), line[1:])
		} else {
			rg.Label(NewRectangleFromInt32(HelpWindow.X+12, HelpWindow.Y+lineY, HelpWindow.Width, ModalLineH), line)
		}

		if line == "-" {
			lineY += ModalLineEmptyH
		} else {
			lineY += ModalLineH
		}
	}
}
