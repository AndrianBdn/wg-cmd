package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/app"
	"github.com/andrianbdn/wg-dir-conf/tutils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultStep2 string

type Step2 struct {
	app     *app.App
	sSize   tea.WindowSizeMsg
	spinner spinner.Model
}

func NewStep2(app *app.App, sSize tea.WindowSizeMsg) Step2 {
	s := spinner.New()
	s.Spinner = spinner.Meter
	return Step2{
		app:     app,
		sSize:   sSize,
		spinner: s,
	}
}

func (m Step2) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Step2) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			//case tea.KeyEnter:
			//	return m, tea.Quit

		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m Step2) View() string {

	s := newStyleSize(m.sSize)

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(""),
		s.xText.Render("Step 2"),
		s.xText.Render(m.spinner.View()),
	)

	bottom := lipgloss.JoinVertical(0,
		s.xTooltip.Render("ENTER=Continue  F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Copy().Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
