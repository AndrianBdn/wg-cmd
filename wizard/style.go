package wizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type style struct {
	xColor     lipgloss.Style
	xStyleBase lipgloss.Style
	xTooltip   lipgloss.Style
	xText      lipgloss.Style
	xList      lipgloss.Style
}

func newStyleSize(sSize tea.WindowSizeMsg) *style {
	s := style{}

	s.xColor = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("7"))
	s.xStyleBase = s.xColor.Copy().Width(sSize.Width).PaddingRight(3)
	s.xTooltip = lipgloss.NewStyle().
		Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).
		Width(sSize.Width).PaddingLeft(2)
	s.xText = s.xStyleBase.Copy().PaddingLeft(3)
	s.xList = s.xText.Copy().PaddingLeft(6)

	return &s
}

func (s *style) header() string {
	return lipgloss.JoinVertical(0,
		s.xStyleBase.Render(" WG Commander Setup"),
		s.xStyleBase.Render("====================="),
		s.xText.Render(""),
	)
}
