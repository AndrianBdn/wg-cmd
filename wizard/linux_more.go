package wizard

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/sysinfo"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type linuxMoreDone struct{}

type linuxMoreScreen struct {
	app           *app.App
	sSize         tea.WindowSizeMsg
	blueprint     backend.ServerBlueprint
	changeNat4    bool
	changeNat6    bool
	changeSystemd bool
	changeText    string
	lastError     string
	pressEnter    bool
}

func newLinuxMoreScreen(app *app.App, sSize tea.WindowSizeMsg, blueprint backend.ServerBlueprint) linuxMoreScreen {
	changeText := ""
	changeNat4 := false
	changeNat6 := false
	changeSystemd := false

	if blueprint.Nat4 && !sysinfo.NatEnabledIPv4() {
		changeText = "IPv4 NAT"
		changeNat4 = true
	}

	if blueprint.Nat6 && !sysinfo.NatEnabledIPv6() {
		if changeText != "" {
			changeText += " and "
		}
		changeText += "IPv6 NAT"
		changeNat6 = true
	}

	if changeText != "" {
		changeText = "enable " + changeText
	}
	if sysinfo.HasSystemd() {
		if changeText != "" {
			changeText += "; "
		}
		changeText += "create systemd services"
		changeSystemd = true
	}

	return linuxMoreScreen{
		app:       app,
		sSize:     sSize,
		blueprint: blueprint,

		changeText:    changeText,
		changeNat4:    changeNat4,
		changeNat6:    changeNat6,
		changeSystemd: changeSystemd,
	}
}

func (m linuxMoreScreen) Init() tea.Cmd {
	return nil
}

func (m linuxMoreScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEsc:
			return m, func() tea.Msg {
				return linuxMoreDone{}
			}

		case tea.KeyEnter:
			if !m.pressEnter {
				m.pressEnter = true

				err := sysinfo.EnableNat(m.changeNat4, m.changeNat6)
				if err != nil {
					m.lastError = err.Error()
				}

				if err == nil && m.changeSystemd {
					err = sysinfo.CreateSystemdStuff(m.blueprint.InterfaceName, m.app.Settings.WireguardDir)
					if err != nil {
						m.lastError = err.Error()
						return m, nil
					}
				}

				return m, nil
			} else {
				if m.lastError != "" {
					return m, tea.Quit
				} else {
					return m, func() tea.Msg {
						return linuxMoreDone{}
					}
				}
			}
		}

	}

	return m, cmd
}

func (m linuxMoreScreen) View() string {
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
				s.xList.Render("•  To skip these changes, press ESC"),
			)
		} else {
			dynamicBlock = s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Success.") +
				s.xColor.Render(" System changes were performed. Press ENTER to continue."))
		}

	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).
			Render("Additional System Changes")),
		s.xText.Render(""),
		s.xText.Render("The Setup can perform additional system configuration: "+m.changeText),
		s.xText.Render(""),
		dynamicBlock,
	)

	bottom := lipgloss.JoinVertical(0,
		s.xTooltip.Render("ENTER=Continue  ESC=Skip"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
