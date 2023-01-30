package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/sysinfo"
	"github.com/andrianbdn/wg-dir-conf/tutils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type endpointStepResult string

const (
	stateInit = iota
	stateDetect
	stateSuccess
)

type endpointScreenStep struct {
	sSize   tea.WindowSizeMsg
	state   int
	spinner spinner.Model
	logs    string
	result  endpointStepResult
}

func newEndpointStep(sSize tea.WindowSizeMsg) endpointScreenStep {
	s := spinner.New()
	s.Spinner = spinner.Meter
	return endpointScreenStep{
		sSize:   sSize,
		spinner: s,
		state:   stateInit,
	}
}

func (m endpointScreenStep) Init() tea.Cmd {
	return nil
}

func (m endpointScreenStep) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case sysinfo.DiscoverStep:
		if msg.Result != "" {
			m.result = endpointStepResult(msg.Result)
			m.logs = m.appendLogs("", msg.Result)
			m.state = stateSuccess
			return m, nil
		}
		m.logs = m.appendLogs("checking "+msg.Service, msg.Log)

		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			return detectEndpoint(msg)
		})

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if m.state == stateInit {
				step := sysinfo.NewDiscoverIPStep()
				m.logs = m.appendLogs("checking "+step.Service, "")
				m.state = stateDetect
				return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
					return detectEndpoint(step)
				})
			}

			if m.state == stateSuccess {
				return m, func() tea.Msg {
					return m.result
				}
			}

			return m, nil

		case tea.KeyF3:
			return m, tea.Quit
		}

	default:
		if m.state == stateDetect {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}

func (m endpointScreenStep) View() string {
	s := newStyleSize(m.sSize)

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render("The Setup will attempt to discover this computer external IPv4 address.\n"),
	)

	if m.state == stateInit {
		top = lipgloss.JoinVertical(0,
			top,
			s.xText.Render("This will require the Internet connection. Press ENTER to continue."),
		)
	}

	if m.state == stateDetect || m.state == stateSuccess {
		logs := strings.ReplaceAll(m.logs, "...\n", m.spinner.View()+"\n")
		top = lipgloss.JoinVertical(0,
			top,
			s.xText.Render(logs),
		)
	}

	if m.state == stateSuccess {
		top = lipgloss.JoinVertical(0,
			top,
			s.xText.Render("Discovery finished. Press ENTER to continue."),
		)
	}

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: The Setup expects that WireGuard(R) UDP port is accessible from the Internet\n"),

		s.xTooltip.Render("ENTER=Continue F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}

func (m endpointScreenStep) appendLogs(status, result string) string {
	if result != "" {
		m.logs = strings.TrimRight(m.logs, " .\n")
		m.logs = m.logs + " ... " + result + "\n"
	}
	if status != "" {
		m.logs = m.logs + status + " ...\n"
	}
	return m.logs
}

func detectEndpoint(s sysinfo.DiscoverStep) tea.Msg {
	return sysinfo.DiscoverIP(s)
}
