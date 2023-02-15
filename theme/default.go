package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	HelpLineBackground lipgloss.Style
	HelpLineKey        lipgloss.Style
	HelpLineHelp       lipgloss.Style

	ViewerTopBar lipgloss.Style
	ViewerMain   lipgloss.Style

	MainTableHeader        lipgloss.Style
	MainTableBody          lipgloss.Style
	MainTableSelected      lipgloss.Style
	MainTableFirst         lipgloss.Style
	MainTableSelectedFirst lipgloss.Style
	MainTableDimmed        lipgloss.Style

	DialogBackground             lipgloss.Style
	DialogTitle                  lipgloss.Style
	DialogInput                  lipgloss.Style
	DialogInputCursor            lipgloss.Style
	DialogButtonFirstLetter      lipgloss.Style
	DialogButtonFocus            lipgloss.Style
	DialogButtonFocusFirstLetter lipgloss.Style

	DialogErrorBackground lipgloss.Style
	DialogErrorTitle      lipgloss.Style
	DialogErrorMessage    lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		HelpLineBackground: style(252, 239, false),
		HelpLineKey:        style(39, 239, true),
		HelpLineHelp:       style(252, 241, false),

		ViewerTopBar: style(39, 240, false),
		ViewerMain:   style(254, 235, false),

		MainTableHeader:        style(33, 235, false),
		MainTableBody:          style(254, 235, false),
		MainTableSelected:      style(254, 240, false),
		MainTableSelectedFirst: style(66, 240, false),
		MainTableFirst:         style(66, 235, false),
		MainTableDimmed:        style(238, 235, true),

		DialogBackground:             style(254, 238, false),
		DialogTitle:                  style(33, 238, false),
		DialogInput:                  style(39, 235, false),
		DialogInputCursor:            style(39, 238, false),
		DialogButtonFirstLetter:      style(33, 238, false),
		DialogButtonFocus:            style(39, 240, false),
		DialogButtonFocusFirstLetter: style(33, 240, false),

		DialogErrorBackground: style(0, 1, false),
		DialogErrorTitle:      style(11, 1, false),
		DialogErrorMessage:    style(15, 1, false),
	}
}
