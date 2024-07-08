package wizard

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type doneScreenResult struct{}

type doneScreen struct {
	app        *app.App
	sSize      tea.WindowSizeMsg
	blueprint  backend.ServerBlueprint
	lastError  string
	pressEnter bool
}

func newDoneScreen(app *app.App, sSize tea.WindowSizeMsg, blueprint backend.ServerBlueprint) doneScreen {
	return doneScreen{
		app:       app,
		sSize:     sSize,
		blueprint: blueprint,
	}
}

func (m doneScreen) Init() tea.Cmd {
	return nil
}

func (m doneScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.pressEnter {
				err := m.app.CreateNewServer(m.blueprint)
				m.lastError = ""
				if err != nil {
					m.lastError = err.Error()
				} else {
					err = m.app.LoadInterface(m.blueprint.InterfaceName)

					if err == nil {
						m.app.Settings.DefaultInterface = m.blueprint.InterfaceName
						_ = m.app.SaveSettings()
						_, _ = m.app.GenerateWireguardConfig()
					}
				}
				m.pressEnter = true
				return m, nil
			} else {
				if m.lastError != "" {
					return m, tea.Quit
				} else {
					return m, func() tea.Msg {
						return doneScreenResult{}
					}
				}
			}
		}
	}

	return m, cmd
}

func (m doneScreen) View() string {
	s := newStyleSize(m.sSize)

	dynamicBlock := ""
	if m.lastError != "" {
		dynamicBlock = lipgloss.JoinVertical(0,
			s.xText.Render("Error: "+lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(m.lastError)+"\n"),

			s.xText.Render("Press ENTER to quit. You can try reading documentation or asking other people "+
				"how to fix the problem."),
		)
	} else {
		if !m.pressEnter {
			dynamicBlock = lipgloss.JoinVertical(0,
				s.xList.Render("•  To proceed, press ENTER"),
				s.xText.Render(""),
			)
		} else {
			dynamicBlock = s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Success.") +
				s.xColor.Render(" Configuration written. Press ENTER to continue."))
		}
	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Ready to perform WireGuard(R)"+
			" configuration.")),
		s.xText.Render(""),
		s.xText.Render("The Setup can now generate necessary configuration files, "+
			"both for WG Commander and for WireGuard(R) itself.\n"),

		dynamicBlock,
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: WG Commander is not approved, sponsored, "+
			"or affiliated with WireGuard(R) or its community"),
		s.xText.Render(""),
		s.xTooltip.Render("ENTER=Continue  F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
