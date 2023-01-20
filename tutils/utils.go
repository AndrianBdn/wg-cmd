package tutils

import "github.com/charmbracelet/lipgloss"

func WPad(l string, sWidth int, s lipgloss.Style) string {
	addW := sWidth - lipgloss.Width(l)
	if addW > 0 {
		wb := s.Copy().Width(addW)
		l = lipgloss.JoinHorizontal(0, l, wb.Render(""))
	}
	return l
}

func HPad(b string, sHeight int, s lipgloss.Style) string {
	addH := sHeight - lipgloss.Height(b)
	if addH > 0 {
		wb := s.Copy().Height(addH)
		b = lipgloss.JoinVertical(0, b, wb.Render(" "))
	}
	return b
}
