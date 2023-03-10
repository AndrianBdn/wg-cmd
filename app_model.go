package main

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/wizard"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	app        *app.App
	wizard     tea.Model
	mainScreen tea.Model
	sSize      tea.WindowSizeMsg
}

func NewAppModel(app *app.App) AppModel {
	a := AppModel{app: app}

	if app.State == nil {
		a.wizard = wizard.NewRootModel(app)
	} else {
		a.mainScreen = NewMainScreen(app, a.sSize)
	}

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

	if _, ok := msg.(wizard.Done); ok {
		a.wizard = nil
		a.mainScreen = NewMainScreen(a.app, a.sSize)
		return a, nil
	}

	if a.wizard != nil {
		w, c := a.wizard.Update(msg)
		a.wizard = w
		return a, c
	}

	if a.mainScreen != nil {
		m, c := a.mainScreen.Update(msg)
		a.mainScreen = m
		return a, c
	}

	return a, nil
}

func (a AppModel) View() string {
	if a.wizard != nil {
		return a.wizard.View()
	}
	if a.mainScreen != nil {
		return a.mainScreen.View()
	}
	return "empty"
}
