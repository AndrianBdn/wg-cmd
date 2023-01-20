package main

import (
	"github.com/76creates/stickers"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	NewPeer    key.Binding
	DeletePeer key.Binding
	Quit       key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NewPeer, k.DeletePeer, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{}, // first column
		{}, // second column
	}
}

type model struct {
	screenW int
	screenH int

	keys  keyMap
	table *stickers.TableSingleType[string]
	help  help.Model
}

func newModel() model {

	keys := keyMap{
		NewPeer: key.NewBinding(
			key.WithKeys("F7"),
			key.WithHelp("F7", "NewPeer"),
		),
		DeletePeer: key.NewBinding(
			key.WithKeys("F8"),
			key.WithHelp("F8", "DelPeer"),
		),
		Quit: key.NewBinding(
			key.WithKeys("F10"),
			key.WithHelp("F10", "Quit"),
		),
	}

	helpModel := help.New()
	helpModel.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	helpModel.Styles.ShortDesc = lipgloss.NewStyle().Background(lipgloss.Color("14")).Foreground(lipgloss.Color("0")).Width(12)
	helpModel.ShortSeparator = " "

	return model{
		keys:  keys,
		help:  helpModel,
		table: newTable(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.screenW = msg.Width
		m.screenH = msg.Height
		m.help.Width = msg.Width
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 1)

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "f10":
			return m, tea.Quit
		case "enter":
			/*
				return m, tea.Batch(
					tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
				)

			*/

		case "down":
			m.table.CursorDown()
		case "up":
			m.table.CursorUp()
		}
	}
	//m.table, cmd = m.table.
	return m, cmd
}

func (m model) View() string {
	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(lipgloss.Left, m.table.Render(), helpView)
}
