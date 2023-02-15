package main

import (
	"github.com/andrianbdn/wg-cmd/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDialogMsgResult struct{}

type TuiDialogMsg struct {
	Title   string
	Message string
	IsError bool
}

func NewTuiDialogMsg(title, message string, isError bool) TuiDialogMsg {
	return TuiDialogMsg{Title: title, Message: message, IsError: isError}
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
	if m.IsError {
		return ErrorView(m.Title, m.Message)
	} else {
		return InfoView(m.Title, m.Message)
	}
}

func InfoView(title, message string) string {
	frameStyle := theme.Current.DialogBackground.Copy().Width(44).Padding(1, 2)
	titleStyle := theme.Current.DialogTitle.Copy().Width(40).Align(lipgloss.Left)
	msgStyle := theme.Current.DialogBackground.Copy().Width(40).Align(lipgloss.Left).Padding(1, 0, 0, 0)

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(title),
			msgStyle.Render(message),
		),
	)
}

func ErrorView(title, message string) string {
	padLR := 2
	w := 34
	frameStyle := theme.Current.DialogErrorBackground.Copy().Width(w+padLR*2).Padding(1, 2)
	titleStyle := theme.Current.DialogErrorTitle.Copy().Width(w).Align(lipgloss.Center)
	msgStyle := theme.Current.DialogErrorMessage.Copy().Width(w).Align(lipgloss.Center).Padding(1, 0)

	return frameStyle.Render(
		lipgloss.JoinVertical(0,
			titleStyle.Render(title),
			msgStyle.Render(message),
		),
	)
}
