package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/backend"
	"github.com/andrianbdn/wg-dir-conf/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type netStepResult struct {
	net4 string
	net6 string
}

type netScreenStep struct {
	sSize   tea.WindowSizeMsg
	net6    string
	net4    string
	prefix4 int
}

func newNetScreenStep(sSize tea.WindowSizeMsg) netScreenStep {
	return netScreenStep{
		sSize: sSize,
		net4:  backend.RandomIP4("10"),
		net6:  backend.RandomIP6(),
	}
}

func (m netScreenStep) Init() tea.Cmd {
	return nil
}

func (m netScreenStep) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	netPrefix := []string{"10", "192", "172"}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "$", "4":
				if string(msg.Runes) == "$" {
					m.prefix4++
				}

				m.net4 = backend.RandomIP4(netPrefix[m.prefix4%len(netPrefix)])
			case "6":
				m.net6 = backend.RandomIP6()
			case "^":
				if m.net6 == "" {
					m.net6 = backend.RandomIP6()
				} else {
					m.net6 = ""
				}
			}

		case tea.KeyEnter:
			r := netStepResult{net4: m.net4, net6: m.net6}
			return m, func() tea.Msg {
				return r
			}

		case tea.KeyF3:
			return m, tea.Quit
		}

	}

	return m, cmd
}

func (m netScreenStep) View() string {
	s := newStyleSize(m.sSize)

	uiNet6 := m.net6
	if uiNet6 == "" {
		uiNet6 = "DISABLED"
	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render("The Setup has randomly generated the following networks:\n"),
		s.xList.Render("•  IPv4 "+m.net4+"\n"),

		s.xList.Render("•  IPv6 "+uiNet6+"\n"),
		s.xText.Render("Press ENTER if you accept these values. Press 4 to change IPv4 network, "+
			"press 6 to change IPv6 network. Press SHIFT-6 (^) to toggle IPv6 network."),
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: only /20 and /64 networks are supported at the moment. Make sure there is no collisions. \n"),

		s.xTooltip.Render("ENTER=Continue"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
