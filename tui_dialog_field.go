package main

import (
	"github.com/andrianbdn/wg-cmd/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDialogCancel struct{}
type TuiDialogValue string

const (
	tdfSelectionField = iota
	tdfSelectionOK
	tdfSelectionCancel
)

type TuiDialogField struct {
	Title           string
	Question        string
	buttonOK1       string
	buttonCancel2   string
	ValidationFunc  func(string) string
	field           textinput.Model
	validationError string
	selectedItem    int
}

func NewTuiDialogName() TuiDialogField {
	m := TuiDialogField{}

	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 30
	ti.Width = 50
	ti.TextStyle = theme.Current.DialogInput
	ti.CursorStyle = theme.Current.DialogInputCursor
	ti.Prompt = ""
	ti.Focus()
	ti.SetCursorMode(textinput.CursorBlink)

	m.field = ti
	m.buttonOK1 = " OK "
	m.buttonCancel2 = "Cancel"

	return m
}

func (m TuiDialogField) Init() tea.Cmd {
	return textinput.Blink
}

func (m TuiDialogField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.validationError != "" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {

			case tea.KeyEscape, tea.KeyEnter:
				m.validationError = ""
				return m, textinput.Blink
			}
		}

	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyTab, tea.KeyDown, tea.KeyUp:
			dir := 1
			if msg.Type == tea.KeyUp {
				dir = -1
			}
			m.selectedItem += dir

			if m.selectedItem > tdfSelectionCancel {
				m.selectedItem = tdfSelectionField
			}

			if m.selectedItem < tdfSelectionField {
				m.selectedItem = tdfSelectionCancel
			}

			if m.selectedItem == tdfSelectionField {
				m.field.Focus()
			} else {
				m.field.Blur()
			}

		case tea.KeyEscape, tea.KeyF3:
			return m, func() tea.Msg {
				return TuiDialogCancel{}
			}

		case tea.KeyEnter:
			if m.selectedItem == tdfSelectionCancel {
				return m, func() tea.Msg {
					return TuiDialogCancel{}
				}
			}

			m.validationError = m.ValidationFunc(m.field.Value())
			if m.validationError == "" {
				v := m.field.Value()
				return m, func() tea.Msg {
					return TuiDialogValue(v)
				}
			}
		}
	}

	m.field, cmd = m.field.Update(msg)
	return m, cmd
}

func (m TuiDialogField) View() string {
	if m.validationError != "" {
		return ErrorView("Error", m.validationError)
	}

	frameStyle := theme.Current.DialogBackground.Copy().Width(54).Padding(1, 2)
	titleStyle := theme.Current.DialogTitle.Copy().Width(50).Align(0.5)
	questionStyle := lipgloss.NewStyle().Width(50)
	//.Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Width(50)

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(m.Title),
			questionStyle.Render(m.Question),
			m.field.View(),
			m.RenderButtons(),
		),
	)

}

func (m TuiDialogField) RenderButtons() string {
	bg := theme.Current.DialogBackground.Copy().Align(lipgloss.Center).Width(50)
	renderButton := func(text string, hl bool) string {
		firstLetter := text[0:1]
		restText := text[1:]
		var l, r string
		var styleBtnText, styleBtnFirstLetter lipgloss.Style
		if hl {
			styleBtnText = theme.Current.DialogButtonFocus
			styleBtnFirstLetter = theme.Current.DialogButtonFocusFirstLetter
			l = "[<"
			r = ">]"
		} else {
			styleBtnText = theme.Current.DialogBackground
			styleBtnFirstLetter = theme.Current.DialogButtonFirstLetter
			l = "[ "
			r = " ]"
		}
		return styleBtnText.Render(l) + styleBtnFirstLetter.Render(firstLetter) + styleBtnText.Render(restText+r)
	}

	return bg.Render(renderButton(m.buttonOK1, m.selectedItem == tdfSelectionOK) +
		theme.Current.DialogBackground.Render(" ") +
		renderButton(m.buttonCancel2, m.selectedItem == tdfSelectionCancel))

}
