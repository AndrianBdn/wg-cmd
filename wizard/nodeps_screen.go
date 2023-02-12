package wizard

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/sysinfo"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type noDepsScreen struct {
	app         *app.App
	sSize       tea.WindowSizeMsg
	hasWg       bool
	hasIPTables bool
}

func newNoDepsScreen(app *app.App, sSize tea.WindowSizeMsg) (ok bool, fail noDepsScreen) {
	n := noDepsScreen{
		app:         app,
		sSize:       sSize,
		hasWg:       sysinfo.HasWireguard(),
		hasIPTables: sysinfo.HasIPTables(),
	}
	return n.hasWg && n.hasIPTables, n
}

func (m noDepsScreen) Init() tea.Cmd {
	return nil
}

func (m noDepsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m noDepsScreen) View() string {
	s := newStyleSize(m.sSize)

	missing := ""

	if !m.hasWg {
		missing = s.xList.Render("•  Wireguard(R) tools (wg, wg-quick)\n")
	}

	if !m.hasIPTables {

		missing = tutils.JoinVerticalNonEmpty(lipgloss.Left,
			missing,
			s.xList.Render("•  iptables (ip6tables)\n"),
		)
	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Welcome to Setup.")),
		s.xText.Render(""),
		s.xText.Render("Unfortunately, some of the required packages are missing from your system:\n"),
		missing,
		s.xText.Render(""),
		s.xText.Render("Please install them and try again."),
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: WG Commander is not approved, sponsored, "+
			"or affiliated with WireGuard(R) or its community"),
		s.xText.Render(""),
		s.xTooltip.Render("ENTER=Continue"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
