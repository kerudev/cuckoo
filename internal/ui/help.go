package ui

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

var HelpLines = []string{
	"Drag&Drop file: changes grid coordinates",
	"#138# - Locked state indicator. Click to unlock",
	"#193#/H - Toggle help window",

	"-",

	"-Grid",
	"LClick/L - Lock coordinates where mouse is over",
	"Wheel[Up/Down] - Zoom in/out",
	"LShift + Wheel[Up/Down] - Scroll right/left",
	"Hold RMouse - Scroll right/left",

	"-Tooltip",
	"Wheel[Up/Down] - Scroll up/down",

	"-",

	"-UI options",
	"Group by: groups coordinates by minute",
	"1-7 - Toggle weekdays (0-6)",
	"Minutes step: group coordinates by bucket",

	"-User options",
	"Tooltip position: change where the tooltip is",
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
		HelpWindow.Height = HelpLineH + 8

		for _, line := range HelpLines {
			if line == "-" {
				HelpWindow.Height += HelpLineEmptyH
			} else {
				HelpWindow.Height += HelpLineH
			}
		}

		HelpWindow.Width = 300
		HelpWindow.X = (S_Screen.Val.W - HelpWindow.Width) / 2
		HelpWindow.Y = (S_Screen.Val.H - HelpWindow.Height) / 2
	}

	lineY := HelpLineH + 4

	// Icon: HELP
	ShowHelp = !rg.WindowBox(HelpWindow.ToFloat32(), "#193# Help and user guide")

	for _, line := range HelpLines {
		if line == "-" {
			rg.Line(NewRectangleFromInt32(HelpWindow.X, HelpWindow.Y+lineY, HelpWindow.Width, HelpLineEmptyH), "")
		} else if line[0] == '-' {
			rg.Line(NewRectangleFromInt32(HelpWindow.X, HelpWindow.Y+lineY, HelpWindow.Width, HelpLineH), line[1:])
		} else {
			rg.Label(NewRectangleFromInt32(HelpWindow.X+12, HelpWindow.Y+lineY, HelpWindow.Width, HelpLineH), line)
		}

		if line == "-" {
			lineY += HelpLineEmptyH
		} else {
			lineY += HelpLineH
		}
	}
}
