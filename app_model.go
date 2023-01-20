package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppModel struct {
	app    *App
	value  string
	dialog tea.Model
}

func NewAppModel(app *App) AppModel {
	a := AppModel{app: app}
	a.value = "Hi there"
	return a
}

func (a AppModel) Init() tea.Cmd {
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(WizardResult); ok {
		a.dialog = nil
		a.value = string(msg)
		return a, nil
	}

	if a.dialog != nil {
		dialog, cmd := a.dialog.Update(msg)
		a.dialog = dialog
		return a, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "f10":
			return a, tea.Quit

		case "enter":
			w := NewWizard()
			cmd := w.Init()
			a.dialog = &w
			return a, cmd
		}

	}
	return a, nil
}

func (a AppModel) View() string {
	if a.dialog != nil {
		return a.dialog.View()
	}

	style := lipgloss.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("15"))
	return style.Render(a.value)
}
