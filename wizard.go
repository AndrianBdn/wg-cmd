package main

//
//import (
//	"github.com/charmbracelet/bubbles/textinput"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//)
//
//type Wizard struct {
//	app         *App
//	ifName      textinput.Model
//	ifNameError string
//	sSize       tea.WindowSizeMsg
//}
//
//func NewWizard(app *App) Wizard {
//	ti := textinput.New()
//	ti.Placeholder = "wg0  "
//	ti.Focus()
//	ti.CharLimit = 4
//	ti.Width = 4
//	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Background(lipgloss.Color("7"))
//	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7"))
//	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
//	ti.Prompt = ""
//
//	return Wizard{
//		app:    app,
//		ifName: ti,
//	}
//}
//
//func (m Wizard) Init() tea.Cmd {
//	return textinput.Blink
//}
//
//func (m Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	if msg, ok := msg.(tea.WindowSizeMsg); ok {
//		m.sSize = msg
//		return m, nil
//	}
//
//	var cmd tea.Cmd
//
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.Type {
//		case tea.KeyF3:
//			return m, tea.Quit
//		case tea.KeyEnter:
//			ifName := m.ifName.Value()
//			m.ifNameError = m.app.validateIfaceArg(ifName)
//			if m.ifNameError != "" {
//				return m, nil
//			}
//			return m, func() tea.Msg {
//				return WizardResult{ifName: m.ifName.Value(), nat: false}
//			}
//		}
//
//	}
//
//	m.ifName, cmd = m.ifName.Update(msg)
//	return m, cmd
//}
//
//func (m Wizard) View() string {
//
//	var xColor = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("7"))
//	var xStyleBase = xColor.Copy().Width(m.sSize.Width).PaddingRight(3)
//	var xTooltip = lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0")).Width(m.sSize.Width)
//	var xText = xStyleBase.Copy().PaddingLeft(3)
//	var xList = xText.Copy().PaddingLeft(6)
//
//	p := lipgloss.JoinHorizontal(0,
//		xColor.Render("   Enter name of a WireGuard(R) network interface: "),
//		m.ifName.View(),
//	)
//	p = wPad(p, m.sSize.Width, xColor)
//
//	errorBlock := xText.Render("")
//	if m.ifNameError != "" {
//		errorBlock = lipgloss.JoinVertical(0,
//			xText.Render(""),
//			xText.Render("Error: "+lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(m.ifNameError)),
//		)
//	}
//
//	top := lipgloss.JoinVertical(0,
//		xStyleBase.Render(" WG Commander Setup"),
//		xStyleBase.Render("====================="),
//		xText.Render(""),
//		xText.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("Welcome to Setup.")),
//		xText.Render(""),
//		xText.Render("This portion of the Setup helps you configure a new WireGuard(R) network interface."),
//		xText.Render(""),
//		xList.Render("•  To proceed, enter an interface name below and press ENTER"),
//		xText.Render(""),
//		xList.Render("•  To quit Setup without configuring the interface, press F3"),
//		xText.Render(""),
//		p,
//		errorBlock,
//	)
//
//	bottom := lipgloss.JoinVertical(0,
//		xText.Render("Note: WG Commander is not approved, sponsored, or affiliated with WireGuard(R) or its community"),
//		xText.Render(""),
//		xTooltip.Render("  ENTER=Continue  F3=Quit"),
//	)
//
//	top = hPad(top, m.sSize.Height-lipgloss.Height(bottom), xColor.Copy().Width(m.sSize.Width))
//	return lipgloss.JoinVertical(0, top, bottom)
//}
