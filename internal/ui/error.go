package ui

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

func DrawError() {
	if ErrorText == "" {
		return
	}

	red40 := rl.Fade(rl.Red, 0.4)
	red80 := rl.Fade(rl.Red, 0.8)

	black := rg.NewColorPropertyValue(rl.Black)
	red := rg.NewColorPropertyValue(rl.Red)

	red30p := rg.NewColorPropertyValue(rl.Fade(rl.Red, 0.3))
	red40p := rg.NewColorPropertyValue(red40)
	red50p := rg.NewColorPropertyValue(rl.Fade(rl.Red, 0.5))

	def_TEXT_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_NORMAL)
	def_TEXT_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED)
	def_TEXT_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED)

	def_BORDER_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL)
	def_BASE_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL)
	def_BORDER_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED)
	def_BASE_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED)
	def_BORDER_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED)
	def_BASE_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED)

	// Set styles based on status
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_NORMAL, black)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, black)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, black)

	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, red)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, red40p)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, red)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, red30p)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, red)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, red50p)

	if ShowHelp {
		rg.SetState(rg.STATE_DISABLED)
	}

	// Draw back button
	// ICON: CROSS
	if rg.Button(CloseErrorButton, "#113#") {
		ErrorText = ""
	}

	if ShowHelp {
		rg.SetState(rg.STATE_NORMAL)
	}

	// Reset style to defaults
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_NORMAL, def_TEXT_COLOR_NORMAL)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, def_TEXT_COLOR_PRESSED)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, def_TEXT_COLOR_FOCUSED)

	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, def_BORDER_COLOR_NORMAL)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, def_BASE_COLOR_NORMAL)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_FOCUSED)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, def_BASE_COLOR_FOCUSED)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_PRESSED)
	rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, def_BASE_COLOR_PRESSED)

	// Draw error background and text
	rg.DrawRectangle(ErrorBox, 2, red80, red40)
	rg.DrawText(ErrorText, ErrorMessageText, int32(rg.TEXT_ALIGN_LEFT), rl.Black)
}
