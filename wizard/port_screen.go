package wizard

import (
	"fmt"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"
	"strings"
)

type portStepResult uint16

type portScreen struct {
	sSize tea.WindowSizeMsg
	port  uint16
}

func newPortScreen(sSize tea.WindowSizeMsg) portScreen {
	return portScreen{
		sSize: sSize,
		port:  51820,
	}
}

func (m portScreen) Init() tea.Cmd {
	return nil
}

func (m portScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "4":
				m.port = 4500
			case "5":
				m.port = 51820
			case "d", "D":
				m.port = 53
			}

		case tea.KeyEnter:
			r := portStepResult(m.port)
			return m, func() tea.Msg {
				return r
			}

		case tea.KeyLeft:
			if m.port > 1 {
				m.port = m.port - 1
			}

		case tea.KeyShiftLeft:
			if m.port > 100 {
				m.port = m.port - 100
			}

		case tea.KeyCtrlShiftLeft:
			if m.port > 1000 {
				m.port = m.port - 1000
			}

		case tea.KeyRight:
			if m.port < math.MaxUint16 {
				m.port = m.port + 1
			}

		case tea.KeyShiftRight:
			if m.port < math.MaxUint16-100 {
				m.port = m.port + 100
			}

		case tea.KeyCtrlShiftRight:
			if m.port < math.MaxUint16-1000 {
				m.port = m.port + 1000
			}

		case tea.KeyF3:
			return m, tea.Quit
		}

	}

	return m, cmd
}

func (m portScreen) View() string {
	s := newStyleSize(m.sSize)

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render("Choose a UDP port for a WireGuard(R) VPN endpoint."),
		s.xText.Render("\nUse the LEFT and RIGHT ARROW keys to increment or decrement port number. Hold SHIFT key"+
			" while pressing the ARROW keys to change value faster. Hold CTRL+SHIFT to increase change speed further.\n"),
	)

	tw := m.sSize.Width - 3 - 3 - 2
	lr := tw * int(m.port) / math.MaxUint16
	rr := tw - lr - 1
	if rr < 0 {
		lr = tw - 1
		rr = 0
	}

	bar := "[" + strings.Repeat("-", lr) + "|" + strings.Repeat("-", rr) + "]"
	bar2 := " " + strings.Repeat(" ", lr) + "^" + strings.Repeat(" ", rr) + " "

	top = lipgloss.JoinVertical(0,
		top,
		s.xText.Render(bar),
		s.xText.Render(bar2),
		s.xText.Render(
			s.xText.Copy().Width(tw+2).Align(0.5).Inline(true).Render(fmt.Sprintf("PORT %d", m.port)),
		),
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: make sure that chosen UDP port is allowed by your firewall settings\n"),

		s.xTooltip.Render("ENTER=Continue  4=Set to 4500 5=Set to 51820"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
