package wizard

import (
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type optScreenResult struct {
	id     string
	result opt
}

type opt struct {
	id    string
	value string
}

type optionScreen struct {
	app      *app.App
	sSize    tea.WindowSizeMsg
	id       string
	opts     []opt
	prompt   string
	selIndex int
}

func newOptionScreen(app *app.App, sSize tea.WindowSizeMsg, id string, options []opt) optionScreen {
	return optionScreen{
		app:   app,
		sSize: sSize,
		id:    id,
		opts:  options,
	}
}

func (m *optionScreen) setPrompt(prompt string) {
	m.prompt = prompt
}

func (m optionScreen) Init() tea.Cmd {
	return nil
}

func (m optionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		return m, nil
	}

	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyEnter:
			r := optScreenResult{id: m.id, result: m.opts[m.selIndex]}
			return m, func() tea.Msg {
				return r
			}

		case tea.KeyUp:
			m.selIndex--
			if m.selIndex < 0 {
				m.selIndex = 0
			}

		case tea.KeyDown:
			m.selIndex++
			if m.selIndex > len(m.opts)-1 {
				m.selIndex = len(m.opts) - 1
			}

		case tea.KeyF3:
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m optionScreen) View() string {
	s := newStyleSize(m.sSize)

	top := lipgloss.JoinVertical(0,
		s.header(),
		s.xText.Render(m.prompt),
		s.xText.Render("\nUse the UP and DOWN ARROW keys to select the option you want, "+
			"and then press ENTER to continue.\n"),
	)

	for idx, opt := range m.opts {
		sel := lipgloss.NewStyle()

		if idx == m.selIndex {
			sel = sel.Inline(true).Background(lipgloss.Color("7")).Foreground(lipgloss.Color("4"))
		}

		top = lipgloss.JoinVertical(0, top, s.xList.Render(sel.Render(opt.value)))
	}

	bottom := lipgloss.JoinVertical(0,
		s.xTooltip.Render("ENTER=Continue  F3=Quit"),
	)

	top = tutils.HPad(top, m.sSize.Height-lipgloss.Height(bottom), s.xColor.Width(m.sSize.Width))
	return lipgloss.JoinVertical(0, top, bottom)
}
