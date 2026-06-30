package ui

import (
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawUIOptions() {
	// Draw option - GroupBy
	rl.DrawText("Group by", Offset.X, Grid.Height+Offset.Y*2+TextPad, FontSize, rl.Black)
	groupByRec := rl.RectangleInt32{X: Offset.X, Y: Grid.Height + Offset.Y*3, Width: 100, Height: 31*2 + 1}
	groupByIdx := int32(S_GroupBy.Val)
	rg.ListView(groupByRec.ToFloat32(), "Wd+Hour;Wd+Hour+Min", nil, &groupByIdx)

	// Prevent ListView from having nothing selected
	if groupByIdx >= 0 {
		S_GroupBy.Set(GroupBy(groupByIdx))
	}

	// Draw option - Weekdays

	// Check the implementation of GuiLoadStyleDefault for additional keys
	// https://github.com/raysan5/raygui/blob/master/src/raygui.h

	rl.DrawText("Weekdays", 120+Offset.X, Grid.Height+Offset.Y*2+TextPad, FontSize, rl.Black)

	def_BORDER_WIDTH := rg.GetStyle(rg.BUTTON, rg.BORDER_WIDTH)

	def_TEXT_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED)
	def_TEXT_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED)

	def_BORDER_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL)
	def_BASE_COLOR_NORMAL := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL)
	def_BORDER_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED)
	def_BASE_COLOR_FOCUSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED)
	def_BORDER_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED)
	def_BASE_COLOR_PRESSED := rg.GetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED)

	rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, 1)

	black := rg.NewColorPropertyValue(rl.Black)

	for wd := range S_Weekdays.Val {
		status := S_Weekdays.Val[wd].Status
		hex := rg.NewColorPropertyValue(S_Weekdays.Val[wd].Color)

		// Set styles based on status
		switch status {
		case StatusDisabled:
			base := rg.NewColorPropertyValue(rl.RayWhite)

			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, def_BORDER_COLOR_NORMAL)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_NORMAL)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, base)

		case StatusOff:
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, black)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, black)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, LerpHex(hex, 0.9))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, LerpHex(hex, 0.8))

		case StatusOn:
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, black)
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, black)

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, LerpHex(hex, 0.8))
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, hex)
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, LerpHex(hex, 0.7))
		}

		button := rl.RectangleInt32{
			X:      120 + Offset.X + BoxPad*int32(wd),
			Y:      Grid.Height + Offset.Y*3,
			Width:  BoxSize,
			Height: BoxSize,
		}

		active := status.Bool()
		rg.Toggle(button.ToFloat32(), strconv.Itoa(wd), &active)

		if status != StatusDisabled {
			S_Weekdays.Val[wd].Status = StatusFromBool(active)

			if !active && All(S_Weekdays.Val[:], func(wd Weekday) bool { return wd.Status != StatusOn }) {
				S_Weekdays.Val[wd].Status = StatusOn
			}
		}

		// Reset style to defaults
		rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, def_BORDER_WIDTH)

		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, def_TEXT_COLOR_PRESSED)
		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, def_TEXT_COLOR_FOCUSED)

		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, def_BORDER_COLOR_NORMAL)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, def_BASE_COLOR_NORMAL)
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, def_BORDER_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, def_BASE_COLOR_FOCUSED)
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, def_BORDER_COLOR_PRESSED)
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, def_BASE_COLOR_PRESSED)
	}

	// Draw option - StepMin
	if S_GroupBy.Eq(GroupByWdHourMin) {
		rl.DrawText("Step of x minutes", 120+Offset.X, Grid.Height+Offset.Y*4+TextPad, FontSize, rl.Black)
		stepMinRec := rl.RectangleInt32{X: 120 + Offset.X, Y: Grid.Height + Offset.Y*5, Width: BoxSize, Height: BoxSize}

		stepMinIdx := int32(S_StepMin.Val)
		rg.ToggleGroup(stepMinRec.ToFloat32(), "#113#;5;10;15;20;30", &stepMinIdx)

		S_StepMin.Set(StepMin(stepMinIdx))
	}
}

func DrawUserOptions() {
	// User option - TooltipPosition
	rl.DrawText("Tooltip position", Offset.X, Grid.Height+Offset.Y*7+TextPad, FontSize, rl.Black)
	positionRec := rl.RectangleInt32{X: Offset.X, Y: Grid.Height + Offset.Y*8, Width: 100, Height: 31*2 + 1}

	positionIdx := int32(Position)
	rg.ListView(positionRec.ToFloat32(), "Grid;Coordinate", nil, &positionIdx)

	// Prevent ListView from having nothing selected
	if positionIdx >= 0 {
		Position = TooltipPosition(positionIdx)
	}

	// User option - Draw options
	rl.DrawText("Draw options", 120+Offset.X, Grid.Height+Offset.Y*7+TextPad, FontSize, rl.Black)

	drawCoordsIcon := "#213#"
	if UserOpt.DrawCoords {
		drawCoordsIcon = "#212#"
	}

	options := []ToggleParams{
		{drawCoordsIcon, &UserOpt.DrawCoords},
		{"#127#", &UserOpt.DrawLines},
		{"#94#", &UserOpt.DrawFade},
		{"#97#", &UserOpt.DrawGrid},
	}

	toggleRec := rl.RectangleInt32{X: 120 + Offset.X, Y: Grid.Height + Offset.Y*8, Width: BoxSize, Height: BoxSize}

	for _, params := range options {
		rg.Toggle(toggleRec.ToFloat32(), params.Icon, params.Ptr)
		toggleRec.X += BoxPad
	}

	if !UserOpt.DrawCoords && !UserOpt.DrawLines && !UserOpt.DrawFade {
		UserOpt.DrawCoords = true
	}
}
