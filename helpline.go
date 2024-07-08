package main

import (
	"github.com/andrianbdn/wg-cmd/theme"
	"github.com/charmbracelet/lipgloss"
)

type helpKey struct {
	key    string
	help   string
	hidden bool
}

func (h helpKey) Render() string {
	b := theme.Current.HelpLineHelp.Width(12)
	return theme.Current.HelpLineKey.Render(h.key) + b.Render(h.help)
}

func RenderHelpLine(w int, keys ...helpKey) string {
	helpline := ""
	for i, k := range keys {
		if k.hidden {
			continue
		}
		helpline += k.Render()
		if i != len(keys)-1 {
			helpline += theme.Current.HelpLineBackground.Render("  ")
		}
	}
	bw := w - lipgloss.Width(helpline)
	return theme.Current.HelpLineBackground.Width(bw).Render(" ") + helpline
}
