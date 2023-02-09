package tutils

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leaanthony/go-ansi-parser"
	"strings"
)

func PlaceDialog(dialog, background string, screenSize tea.WindowSizeMsg, bgStyle lipgloss.Style) string {
	result := ""

	cleanBg, _ := ansi.Cleanse(background)
	bw, bh := screenSize.Width, screenSize.Height
	w, h := lipgloss.Size(dialog)

	posx := (bw - w) / 2
	posy := bh/2 - h/2 - h/4

	bLines := strings.Split(cleanBg, "\n")
	dLines := strings.Split(dialog, "\n")

	for i := 0; i < bh; i++ {

		if i >= posy && i < posy+h {
			if posx > 0 {
				result += bgStyle.Render(bLines[i][0:posx])
			}
			result += dLines[i-posy]
			st := posx + w
			if st < bw {
				result += bgStyle.Render(bLines[i][st:])
			}
		} else {
			result += bgStyle.Render(bLines[i])
		}

		if i != bh-1 {
			result += "\n"
		}
	}

	return result
}
