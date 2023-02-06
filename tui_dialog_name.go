package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDialogCancel struct{}
type TuiDialogValue string

type TuiDialogName struct {
	Title           string
	Question        string
	ValidationFunc  func(string) string
	field           textinput.Model
	validationError string
	selectedButton  int
}

func NewTuiDialogName() TuiDialogName {
	m := TuiDialogName{}
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 30
	ti.Width = 50
	ti.TextStyle = lipgloss.NewStyle().Background(lipgloss.Color("6")).Foreground(lipgloss.Color("0"))
	ti.CursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("15")).Foreground(lipgloss.Color("7"))
	ti.Prompt = ""
	ti.Focus()
	ti.SetCursorMode(textinput.CursorBlink)

	m.field = ti

	return m
}

func (m TuiDialogName) Init() tea.Cmd {
	return textinput.Blink
}

func (m TuiDialogName) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case tea.KeyEscape, tea.KeyF3:
			return m, func() tea.Msg {
				return TuiDialogCancel{}
			}

		case tea.KeyEnter:
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

func (m TuiDialogName) View() string {
	if m.validationError != "" {
		return ErrorView("Error", m.validationError)
	}

	frameStyle := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Width(54).Padding(1, 2)
	titleStyle := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("4")).Width(50).Align(0.5)
	questionStyle := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Width(50)
	buttonStyle := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Align(0.5).Width(50)

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(m.Title),
			questionStyle.Render(m.Question),
			m.field.View(),
			buttonStyle.Render("[<  Enter  >] [   ESC   ]"),
		),
	)

}
