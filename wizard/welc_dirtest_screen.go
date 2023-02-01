package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/app"
	"github.com/andrianbdn/wg-dir-conf/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type welcomeDirTestResult struct{}

type welcomeDirTestScreen struct {
	app      *app.App
	sSize    tea.WindowSizeMsg
	dirError string
	testDone bool
}

func newWelcomeDirTestScreen(app *app.App, sSize tea.WindowSizeMsg) welcomeDirTestScreen {
	return welcomeDirTestScreen{
		app:   app,
		sSize: sSize,
	}
}

func (m welcomeDirTestScreen) Init() tea.Cmd {
	return nil
}

func (m welcomeDirTestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF3:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.testDone {
				m.dirError = m.app.TestDirectories()
				m.testDone = true
				return m, nil
			} else {

				if m.dirError != "" {
					return m, tea.Quit
				} else {
					return m, func() tea.Msg {
						return welcomeDirTestResult{}
					}
				}

			}

		}

	}

	return m, cmd
}

func (m welcomeDirTestScreen) View() string {
	s := newStyleSize(m.sSize)

	dynamicBlock := s.xText.Render("")
	if m.dirError != "" {
		dynamicBlock = lipgloss.JoinVertical(0,
			s.xText.Render("Error: "+lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(m.dirError)+"\n"),
			s.xText.Render("Usually it happens because WireGuard(R) directories "+
				"are not accessible by current user.\n"),
			s.xText.Render("Press ENTER to quit. You can try resolving the problem by using sudo."),
		)
	} else {
		if !m.testDone {
			dynamicBlock = lipgloss.JoinVertical(0,
				s.xList.Render("•  To proceed, press ENTER"),
				s.xText.Render(""),
				s.xList.Render("•  To quit Setup without testing directories, press F3"),
				s.xText.Render(""),
			)
		} else {
			dynamicBlock = s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Success.") +
				s.xColor.Render(" Directory(s) are writable. Press ENTER to continue with configuration."))
		}

	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Welcome to Setup.")),
		s.xText.Render(""),
		s.xText.Render("Before configuring WireGuard(R) interface the Setup needs to check prerequisites.\n"),
		s.xText.Render("WG Commander uses "+m.app.Settings.DatabaseDir+" to store interface settings. "+
			"In addition generated "+
			"WireGuard(R) configuration files will be placed to "+m.app.Settings.WireguardDir),

		s.xText.Render("\nThe Setup will test if these directories are writable.\n "),

		dynamicBlock,
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: WG Commander is not approved, sponsored, "+
			"or affiliated with WireGuard(R) or its community"),
		s.xText.Render(""),
		s.xTooltip.Render("ENTER=Continue  F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
