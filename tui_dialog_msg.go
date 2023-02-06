package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDialogMsgResult struct{}

type TuiDialogMsg struct {
	Title   string
	Message string
}

func NewTuiDialogMsg(title, message string) TuiDialogMsg {
	return TuiDialogMsg{Title: title, Message: message}
}

func (m TuiDialogMsg) Init() tea.Cmd {
	return nil
}

func (m TuiDialogMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyEnter:
			return m, func() tea.Msg {
				return TuiDialogMsgResult{}
			}
		}
	}
	return m, cmd
}

func (m TuiDialogMsg) View() string {
	return ErrorView(m.Title, m.Message)
}

func ErrorView(title, message string) string {

	frameStyle := lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("0")).Width(34).Padding(1, 2)
	titleStyle := lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("11")).Width(30).Align(0.5)
	msgStyle := lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("15")).Width(30).Align(0.5).Padding(1, 0)

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(title),
			msgStyle.Render(message),
		),
	)

}
