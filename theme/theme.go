package theme

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// Current is a current theme
var Current Theme = DefaultTheme()

func style(fg int, bg int, bold bool) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(strconv.Itoa(fg))).
		Background(lipgloss.Color(strconv.Itoa(bg))).
		Bold(bold)
}
