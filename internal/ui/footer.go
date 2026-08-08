package ui

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

func DrawFooter() {
	if BlockUI {
		rg.SetState(rg.STATE_DISABLED)
	}

	DrawUIOptions()
	DrawUserOptions()

	if BlockUI {
		rg.SetState(rg.STATE_NORMAL)
	}

	// Horizontal line
	lineY := Footer.Y + Offset.Y*6 - TextPad/2
	rl.DrawLine(Offset.X, lineY, 290, lineY, rl.Gray)

	// Vertical line
	lineX := 150 + Offset.X + BoxPad*6
	rl.DrawLine(lineX, Footer.Y+Offset.Y+TextPad, lineX, S_Screen.Val.H-Offset.Y, rl.Gray)

	drawFooterData()
}

func drawFooterData() {
	footerX := S_Screen.Val.W - Offset.X - FooterW
	footerY := Footer.Y + Offset.Y + FontSize*2

	// text := "Drop file to change sample"
	// textW := rl.MeasureText(text, FooterFontSize)

	text := "Count of crons & jobs"
	textW := rl.MeasureText(text, FontSize)

	rl.DrawText(text, S_Screen.Val.W-textW-Offset.X, Footer.Y+Offset.Y+TextPad, FontSize, rl.Black)

	totalCrons := 0
	totalJobs := 0

	for wd, count := range WdCounts {
		var s string
		if S_Weekdays.Val[wd].Status == StatusOn {
			s = fmt.Sprintf("%d (%d)", count.Crons, count.Jobs)

			totalCrons += count.Crons
			totalJobs += count.Jobs
		} else {
			s = "0 (0)"
		}

		rl.DrawCircle(footerX-TextPad, footerY+TextPad*2*int32(wd)+FontRadius, float32(FontRadius), S_Weekdays.Val[wd].Color)
		rl.DrawText(s, footerX+TextPad, footerY+TextPad*2*int32(wd), FontSize, rl.Black)
	}

	rl.DrawText(fmt.Sprintf("%d (%d)", totalCrons, totalJobs), footerX+TextPad, footerY+TextPad*2*WEEKDAYS, FontSize, rl.Black)

	// texts := []string{
	// 	fmt.Sprintf("Scale: x%.2f", C_Zoom.Scale),
	// 	fmt.Sprintf("Cell.W: %.2f", Cell.W),
	// 	fmt.Sprintf("Cell.H: %.2f", Cell.H),
	// 	fmt.Sprint("[L]ocked: ", S_IsMouseLocked.Val),
	// }

	// rl.DrawText(strings.Join(texts, "\n"), footerX+TextPad, footerY+TextPad, FooterFontSize, rl.Black)
	// rl.DrawRectangleLines(footerX, footerY, FooterW, int32(len(texts))*FooterFontSize+TextPad*2, rl.Black)
}
