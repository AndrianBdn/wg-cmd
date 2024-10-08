package wizard

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/tutils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type interfaceScreenResult string

type interfaceScreen struct {
	app         *app.App
	ifName      textinput.Model
	ifNameError string
	sSize       tea.WindowSizeMsg
}

func newInterfaceScreen(app *app.App, sSize tea.WindowSizeMsg) interfaceScreen {
	ti := textinput.New()
	ti.Placeholder = "wg0  "
	ti.Focus()
	ti.CharLimit = 4
	ti.Width = 4
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Background(lipgloss.Color("7"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
	ti.Prompt = ""

	return interfaceScreen{
		app:    app,
		ifName: ti,
		sSize:  sSize,
	}
}

func (m interfaceScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m interfaceScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			ifName := m.ifName.Value()
			m.ifNameError = m.app.ValidateIfaceArg(ifName)
			if m.ifNameError != "" {
				return m, nil
			}
			return m, func() tea.Msg {
				return interfaceScreenResult(m.ifName.Value())
			}
		}
	}

	m.ifName, cmd = m.ifName.Update(msg)
	return m, cmd
}

func (m interfaceScreen) View() string {
	s := newStyleSize(m.sSize)

	p := lipgloss.JoinHorizontal(0,
		s.xColor.Render("   Enter name of a new WireGuard(R) network interface: "),
		m.ifName.View(),
	)

	p = tutils.WPad(p, m.sSize.Width, s.xColor)

	errorBlock := s.xText.Render("")
	if m.ifNameError != "" {
		errorBlock = lipgloss.JoinVertical(0,
			s.xText.Render(""),
			s.xText.Render("Error: "+lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(m.ifNameError)),
		)
	}

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render("This portion of the Setup helps you configure a new WireGuard(R) network interface."),
		s.xText.Render(""),
		s.xList.Render("•  To proceed, enter an interface name below and press ENTER"),
		s.xText.Render(""),
		p,
		errorBlock,
	)

	bottom := lipgloss.JoinVertical(0,
		s.xText.Render("Note: WG Commander is not approved, sponsored, or affiliated with WireGuard(R) or its community"),
		s.xText.Render(""),
		s.xTooltip.Render("ENTER=Continue  F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
