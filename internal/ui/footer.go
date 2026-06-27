package ui

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

func DrawFooter() {
	footerX := S_Screen.Val.W - Offset.X - FooterW
	footerY := Grid.H + Offset.Y*2 + FontSize*2

	text := "Drop file to change sample"
	textW := rl.MeasureText(text, FooterFontSize)

	rl.DrawText(text, S_Screen.Val.W-textW-Offset.X, Grid.H+Offset.Y*2, FooterFontSize, rl.Black)

	texts := []string{
		fmt.Sprintf("Scale: x%.2f", ZoomScale),
		fmt.Sprintf("Cell.W: %.2f", Cell.W),
		fmt.Sprintf("Cell.H: %.2f", Cell.H),
		fmt.Sprint("[L]ocked: ", S_IsMouseLocked.Val),
	}

	rl.DrawText(strings.Join(texts, "\n"), footerX+TextPad, footerY+TextPad, FooterFontSize, rl.Black)
	rl.DrawRectangleLines(footerX, footerY, FooterW, int32(len(texts)+1)*FooterFontSize+TextPad*2, rl.Black)
}
