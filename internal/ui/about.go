package ui

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/shared"
)

var AboutLines = []string{
	" ", // Reserved to LabelButton: "Open issue"
	" ", // Reserved to LabelButton: "Home"

	"-",

	"cuckoo helps you visualize jobs and their crons as coordinates to:",
	"- Get insights about your crons with just a glance.",
	"- Know the next free spot where you can place a new cron.",
	"- Identify periods where there is a work overload.",
}

func DrawAboutButton() {
	if ShowAbout {
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
	}

	// Icon: INFO
	if rg.Button(AboutButton, "#191#") {
		// Set to true as the button can only be pressed when the state is STATE_NORMAL
		ShowAbout = true
	}
}

func DrawAbout() {
	// Inspired from: https://github.com/raysan5/rguistyler/blob/3be4d40/src/gui_window_about.h

	rl.DrawRectangle(0, 0, S_Screen.Val.W, S_Screen.Val.H, rl.Fade(rl.White, 0.85))

	if AboutWindow.Height == 0 || S_Screen.HasChanged() {
		// Line + margin
		AboutWindow.Height = ModalLineH + 8

		for _, line := range AboutLines {
			if line == "-" {
				AboutWindow.Height += ModalLineEmptyH
			} else {
				AboutWindow.Height += ModalLineH
			}
		}

		AboutWindow.Width = 370
		AboutWindow.X = (S_Screen.Val.W - AboutWindow.Width) / 2
		AboutWindow.Y = (S_Screen.Val.H - AboutWindow.Height) / 2
	}

	lineY := ModalLineH + 4

	// Icon: INFO
	ShowAbout = !rg.WindowBox(AboutWindow.ToFloat32(), fmt.Sprintf("#191# About cuckoo (v%s)", VERSION))

	// Icon: WARNING
	if rg.LabelButton(
		NewRectangleFromInt32(AboutWindow.X+12, AboutWindow.Y+lineY, AboutWindow.Width, ModalLineH),
		"#220# Found a bug or have an idea? Click to open an issue!",
	) {
		rl.OpenURL("https://github.com/kerudev/cuckoo/issues/new")
	}

	// Icon: HOUSE
	if rg.LabelButton(
		NewRectangleFromInt32(AboutWindow.X+12, AboutWindow.Y+lineY+ModalLineH, AboutWindow.Width, ModalLineH),
		"#185# Homepage (repository)",
	) {
		rl.OpenURL("https://github.com/kerudev/cuckoo")
	}

	for _, line := range AboutLines {
		if line == "-" {
			rg.Line(NewRectangleFromInt32(AboutWindow.X, AboutWindow.Y+lineY, AboutWindow.Width, ModalLineEmptyH), "")
		} else {
			rg.Label(NewRectangleFromInt32(AboutWindow.X+12, AboutWindow.Y+lineY, AboutWindow.Width, ModalLineH), line)
		}

		if line == "-" {
			lineY += ModalLineEmptyH
		} else {
			lineY += ModalLineH
		}
	}
}
