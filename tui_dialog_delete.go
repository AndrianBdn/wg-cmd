package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDialogYes struct{}
type TuiDialogNo struct{}

type TuiDialogYesNo struct {
	Title          string
	Message        string
	selectedButton int
}

func NewTuiDialogYesNo(title, message string) TuiDialogYesNo {
	return TuiDialogYesNo{Title: title, Message: message}
}

func (m TuiDialogYesNo) Init() tea.Cmd {
	return nil
}

func (m TuiDialogYesNo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "Y", "y":
				return m, func() tea.Msg { return TuiDialogYes{} }

			case "N", "n":
				return m, func() tea.Msg { return TuiDialogNo{} }
			}

		case tea.KeyRight:
			m.selectedButton = 1
			return m, nil

		case tea.KeyLeft:
			m.selectedButton = 0
			return m, nil

		case tea.KeyTab:
			if m.selectedButton == 0 {
				m.selectedButton = 1
			} else {
				m.selectedButton = 0
			}
			return m, nil

		case tea.KeyEscape:
			return m, func() tea.Msg {
				return TuiDialogNo{}
			}

		case tea.KeyEnter:
			if m.selectedButton == 0 {
				return m, func() tea.Msg {
					return TuiDialogYes{}
				}
			} else {
				return m, func() tea.Msg {
					return TuiDialogNo{}
				}
			}

		}
	}

	return m, cmd
}

func (m TuiDialogYesNo) View() string {

	redBg := lipgloss.NewStyle().Background(lipgloss.Color("1"))
	frameStyle := redBg.Copy().Foreground(lipgloss.Color("0")).Width(44).Padding(1, 2)
	titleStyle := redBg.Copy().Foreground(lipgloss.Color("11")).Width(40).Align(0.5)
	msgStyle := redBg.Copy().Foreground(lipgloss.Color("15")).Width(40).Align(0.5).Padding(0, 0, 1, 0)

	hStyle := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0"))
	hBright := lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("15"))

	redYellow := redBg.Copy().Foreground(lipgloss.Color("11"))

	btnStyle := redBg.Copy().Width(40).Align(.5)

	yes := redBg.Render("[ ") + redYellow.Render("Y") + redBg.Render("es ]")
	no := redBg.Render("[ ") + redYellow.Render("N") + redBg.Render("o ]")

	if m.selectedButton == 0 {
		yes = hStyle.Render("[ ") + hBright.Render("Y") + hStyle.Render("es ]")
	} else {
		no = hStyle.Render("[ ") + hBright.Render("N") + hStyle.Render("o ]")
	}

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(m.Title),
			msgStyle.Render(m.Message),
			btnStyle.Render(lipgloss.JoinHorizontal(0, yes, redBg.Render(" "), no)),
		),
	)

}
