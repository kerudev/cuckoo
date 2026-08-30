package ui

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/shared"
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
	footerX := 150 + Offset.X + BoxPad*6
	footerY := Footer.Y + Offset.Y + FontSize*2

	text := "Count of crons & jobs"
	rl.DrawText(text, footerX+TextPad*2, Footer.Y+Offset.Y+TextPad, FontSize, rl.Black)

	totalCrons := 0
	totalJobs := 0

	wdCircle := rl.Vector2{
		X: float32(footerX + TextPad*2 + int32(CoordRadius)),
		Y: float32(footerY + FontRadius),
	}

	for wd, count := range WdCounts {
		var s string
		if S_Weekdays.Val[wd].Status == StatusOn {
			s = fmt.Sprintf("%d (%d)", count.Crons, count.Jobs)

			totalCrons += count.Crons
			totalJobs += count.Jobs
		} else {
			s = "0 (0)"
		}

		var color rl.Color
		if S_Weekdays.Val[wd].Status == StatusOn {
			color = S_Weekdays.Val[wd].Color
		} else {
			color = rl.White
		}

		rl.DrawCircle(int32(wdCircle.X), int32(wdCircle.Y), float32(FontRadius), color)
		rl.DrawRing(wdCircle, float32(FontRadius)-1, float32(FontRadius)+1, 0, 360, 16, rl.Black)
		rl.DrawText(s, footerX+TextPad*4+int32(CoordDiameter), footerY+TextPad*2*int32(wd), FontSize, rl.Black)

		wdCircle.Y += float32(TextPad * 2)
	}

	// Draw sum of all coords & jobs
	var wds []int
	for wd := range WEEKDAYS {
		if S_Weekdays.Val[wd].Status == StatusOn {
			wds = append(wds, wd)
		}
	}

	segments := float32(len(wds))

	angleFactor := float32(360) / segments
	angle := float32(270)

	for _, wd := range wds {
		if S_Weekdays.Val[wd].Status != StatusOn {
			continue
		}

		rl.DrawCircleSector(
			wdCircle,
			float32(FontRadius),
			angle,
			angle+angleFactor,
			8,
			S_Weekdays.Val[wd].Color,
		)

		angle += angleFactor
	}

	rl.DrawRing(wdCircle, float32(FontRadius)-1, float32(FontRadius)+1, 0, 360, 16, rl.Black)

	rl.DrawText(
		fmt.Sprintf("%d (%d)", totalCrons, totalJobs),
		footerX+TextPad*4+int32(CoordDiameter),
		footerY+TextPad*2*WEEKDAYS,
		FontSize,
		rl.Black,
	)
}
