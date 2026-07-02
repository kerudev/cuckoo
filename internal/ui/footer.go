package ui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

func DrawFooter() {
	footerX := S_Screen.Val.W - Offset.X - FooterW
	footerY := Grid.Height + Offset.Y*2 + FontSize*2

	// text := "Drop file to change sample"
	// textW := rl.MeasureText(text, FooterFontSize)

	text := "Count of crons & jobs"
	textW := rl.MeasureText(text, FontSize)

	rl.DrawText(text, S_Screen.Val.W-textW-Offset.X, Grid.Height+Offset.Y*2, FontSize, rl.Black)

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

		rl.DrawCircle(footerX+TextPad, footerY+TextPad*2*int32(wd)+FontRadius, float32(FontRadius), S_Weekdays.Val[wd].Color)
		rl.DrawText(s, footerX+TextPad+FontRadius*2, footerY+TextPad*2*int32(wd), FontSize, rl.Black)
	}

	rl.DrawText(fmt.Sprintf("%d (%d)", totalCrons, totalJobs), footerX+TextPad+FontRadius*2, footerY+TextPad*2*WEEKDAYS, FontSize, rl.Black)

	// texts := []string{
	// 	fmt.Sprintf("Scale: x%.2f", C_Zoom.Scale),
	// 	fmt.Sprintf("Cell.W: %.2f", Cell.W),
	// 	fmt.Sprintf("Cell.H: %.2f", Cell.H),
	// 	fmt.Sprint("[L]ocked: ", S_IsMouseLocked.Val),
	// }

	// rl.DrawText(strings.Join(texts, "\n"), footerX+TextPad, footerY+TextPad, FooterFontSize, rl.Black)
	// rl.DrawRectangleLines(footerX, footerY, FooterW, int32(len(texts))*FooterFontSize+TextPad*2, rl.Black)
}
