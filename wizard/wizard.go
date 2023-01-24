package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/app"
	tea "github.com/charmbracelet/bubbletea"
)

type Result struct {
	ifName string
	nat    bool
}

type RootModel struct {
	app           *app.App
	stepInterface StepInterfaceName
	currentModel  tea.Model
	result        Result
	sSize         tea.WindowSizeMsg
}

func NewRootModel(app *app.App) RootModel {
	m := RootModel{}
	m.app = app
	m.currentModel = NewStepInterfaceName(app)
	return m
}

func (m RootModel) Init() tea.Cmd {
	if m.currentModel != nil {
		return m.currentModel.Init()
	}
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
	}

	if msg, ok := msg.(ifNameMsg); ok {
		m.result.ifName = string(msg)
		//step2 := NewStep2(m.app, m.sSize)

		//opts := []opt{
		//	{"nat", "Use NAT to provide the Internet access to the VPN"},
		//	{"no", "Skip NAT, only setup basic networking"},
		//}
		//optStep := newGenericOption(m.app, m.sSize, "nat", opts)
		//optStep.setPrompt("Choose if the Setup should configure Network Address Translation (NAT) for the VPN network")
		//
		//m.currentModel = optStep
		//return m, optStep.Init()

		portStep := newPortScreen(m.sSize)
		m.currentModel = portStep
		return m, portStep.Init()
	}

	if m.currentModel != nil {
		w, c := m.currentModel.Update(msg)
		m.currentModel = w
		return m, c
	}
	return m, nil
}

func (m RootModel) View() string {
	if m.currentModel != nil {
		return m.currentModel.View()
	}
	return "HELLO WORLD"
}
