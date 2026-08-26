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
	rl.DrawText("Group by", Offset.X, Footer.Y+Offset.Y+TextPad, FontSize, rl.Black)
	groupByRec := NewRectangleFromInt32(Offset.X, Footer.Y+Offset.Y*2, 100, ListViewItemH*2+2)
	groupByIdx := int32(S_GroupBy.Val)
	rg.ListView(groupByRec, "Wd+Hour;Wd+Hour+Min", nil, &groupByIdx)

	// Prevent ListView from having nothing selected
	if groupByIdx >= 0 {
		S_GroupBy.Set(GroupBy(groupByIdx))
	}

	// Draw option - Weekdays

	// Check the implementation of GuiLoadStyleDefault for additional keys
	// https://github.com/raysan5/raygui/blob/master/src/raygui.h

	rl.DrawText("Weekdays", 120+Offset.X, Footer.Y+Offset.Y+TextPad, FontSize, rl.Black)

	rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, 1)

	black := Style["BLACK"]
	base := Style["RAYWHITE"]

	for wd := range WEEKDAYS {
		status := S_Weekdays.Val[wd].Status
		hex := rg.NewColorPropertyValue(S_Weekdays.Val[wd].Color)

		// Set styles based on status
		switch status {
		case StatusDisabled:
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, Style["DEF_BORDER_COLOR_NORMAL"])
			rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, Style["DEF_BORDER_COLOR_NORMAL"])

			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, Style["DEF_BORDER_COLOR_NORMAL"])
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, Style["DEF_BORDER_COLOR_NORMAL"])
			rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, base)
			rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, Style["DEF_BORDER_COLOR_NORMAL"])
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

		active := status.Bool()

		button := NewRectangleFromInt32(120+Offset.X+BoxPad*int32(wd), Footer.Y+Offset.Y*2, BoxSize, BoxSize)
		rg.Toggle(button, strconv.Itoa(wd), &active)

		if status != StatusDisabled {
			S_Weekdays.Val.SetStatus(wd, active)
		}

		// Reset style to defaults
		rg.SetStyle(rg.BUTTON, rg.BORDER_WIDTH, Style["BUTTON_BORDER_WIDTH"])

		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, Style["DEF_TEXT_COLOR_PRESSED"])
		rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, Style["DEF_TEXT_COLOR_FOCUSED"])

		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, Style["DEF_BORDER_COLOR_NORMAL"])
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_NORMAL, Style["DEF_BASE_COLOR_NORMAL"])
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, Style["DEF_BORDER_COLOR_FOCUSED"])
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_FOCUSED, Style["DEF_BASE_COLOR_FOCUSED"])
		rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, Style["DEF_BORDER_COLOR_PRESSED"])
		rg.SetStyle(rg.DEFAULT, rg.BASE_COLOR_PRESSED, Style["DEF_BASE_COLOR_PRESSED"])
	}

	// Draw option - StepMin
	if S_GroupBy.Val == GroupByWdHour {
		rg.SetState(rg.STATE_DISABLED)
	}

	rl.DrawText("Minutes step", 120+Offset.X, Footer.Y+Offset.Y*3+TextPad, FontSize, rl.Black)
	stepMinRec := NewRectangleFromInt32(120+Offset.X, Footer.Y+Offset.Y*4, BoxSize, BoxSize)

	stepMinIdx := int32(S_StepMin.Val)
	rg.ToggleGroup(stepMinRec, "1;5;10;15;20;30", &stepMinIdx)

	S_StepMin.Set(StepMin(stepMinIdx))

	if !BlockUI && S_GroupBy.Val == GroupByWdHour {
		rg.SetState(rg.STATE_NORMAL)
	}
}

func DrawUserOptions() {
	// User option - TooltipPosition
	rl.DrawText("Tooltip position", Offset.X, Footer.Y+Offset.Y*6+TextPad, FontSize, rl.Black)
	positionRec := NewRectangleFromInt32(Offset.X, Footer.Y+Offset.Y*7, 100, ListViewItemH*2+2)

	positionIdx := int32(S_Position.Val)
	rg.ListView(positionRec, "Grid;Coordinate", nil, &positionIdx)

	// Prevent ListView from having nothing selected
	if positionIdx >= 0 {
		S_Position.Set(TooltipPosition(positionIdx))
	}

	// User option - Draw options
	rl.DrawText("Draw options", 120+Offset.X, Footer.Y+Offset.Y*6+TextPad, FontSize, rl.Black)

	var drawCoordsIcon string
	if UserOpt.DrawCoords {
		// Icon: BREAKPOINT_ON
		drawCoordsIcon = "#212#"
	} else {
		// Icon: BREAKPOINT_OFF
		drawCoordsIcon = "#213#"
	}

	options := []ToggleParams{
		{Icon: drawCoordsIcon, Ptr: &UserOpt.DrawCoords},
		// Icon: WAVE_TRIANGULAR
		{Icon: "#127#", Ptr: &UserOpt.DrawLines},
		// Icon: DITHERING
		{Icon: "#94#", Ptr: &UserOpt.DrawFade},
		// Icon: GRID
		{Icon: "#97#", Ptr: &UserOpt.DrawGrid},
	}

	toggleRec := NewRectangleFromInt32(120+Offset.X, Footer.Y+Offset.Y*7, BoxSize, BoxSize)

	for _, params := range options {
		rg.Toggle(toggleRec, params.Icon, params.Ptr)
		toggleRec.X += float32(BoxPad)
	}

	if !UserOpt.DrawCoords && !UserOpt.DrawLines && !UserOpt.DrawFade {
		UserOpt.DrawCoords = true
	}
}
