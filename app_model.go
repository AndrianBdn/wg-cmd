package main

import (
	"github.com/andrianbdn/wg-dir-conf/app"
	"github.com/andrianbdn/wg-dir-conf/wizard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppModel struct {
	app    *app.App
	wizard tea.Model
	sSize  tea.WindowSizeMsg
	iface  string
}

func NewAppModel(app *app.App) AppModel {
	a := AppModel{app: app}
	a.wizard = wizard.NewRootModel(app)
	return a
}

func (a AppModel) Init() tea.Cmd {
	if a.wizard != nil {
		return a.wizard.Init()
	}
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		a.sSize = msg
	}

	if msg, ok := msg.(wizard.Done); ok {
		a.wizard = nil
		a.iface = msg.InterfaceName
		return a, nil
	}

	if a.wizard != nil {
		w, c := a.wizard.Update(msg)
		a.wizard = w
		return a, c
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyF3:
			return a, tea.Quit
		}

	}
	return a, nil
}

func (a AppModel) View() string {
	if a.wizard != nil {
		return a.wizard.View()
	}

	style := lipgloss.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("15"))
	return style.Render("Hello World " + a.iface)
}
