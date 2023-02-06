package tutils

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/leaanthony/go-ansi-parser"
	"strings"
)

func PlaceDialog(dialog, background string, bgStyle lipgloss.Style) string {
	result := ""

	cleanBg, _ := ansi.Cleanse(background)
	bw, bh := lipgloss.Size(cleanBg)
	w, h := lipgloss.Size(dialog)

	posx := (bw - w) / 2
	posy := (bh-h)/2 - 10

	bLines := strings.Split(cleanBg, "\n")
	dLines := strings.Split(dialog, "\n")

	for i := 0; i < bh; i++ {

		if i >= posy && i < posy+h {
			result += bgStyle.Render(bLines[i][0:posx])
			result += dLines[i-posy]
			st := posx + w
			result += bgStyle.Render(bLines[i][st:])
		} else {
			result += bgStyle.Render(bLines[i])
		}

		if i != bh-1 {
			result += "\n"
		}
	}
	//for i := 0; i < h; i++ {
	//	result += dLines[i] + "\n"
	_ = posx
	//}
	//
	//leftH := bh - posy - h
	//for i := 0; i < leftH; i++ {
	//	result += bgStyle.Render(bLines[i]) + "\n"
	//}

	return result
}
