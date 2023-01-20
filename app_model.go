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
	//if msg, ok := msg.(WizardResult); ok {
	//	//
	//	fmt.Println(msg)
	//	return a, nil
	//}

	if a.wizard != nil {
		w, c := a.wizard.Update(msg)
		a.wizard = w
		return a, c
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		//case "q", "ctrl+c", "f10":
		//	return a, tea.Quit

		//case "enter":
		//	w := NewWizard()
		//	cmd := w.Init()
		//	a.dialog = &w
		//	return a, cmd
		}

	}
	return a, nil
}

func (a AppModel) View() string {
	if a.wizard != nil {
		return a.wizard.View()
	}

	style := lipgloss.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("15"))
	return style.Render("Hello World")
}
