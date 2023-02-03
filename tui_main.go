package main

import (
	"github.com/76creates/stickers"
	"github.com/andrianbdn/wg-dir-conf/app"
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

type mainScreenTable struct {
	app   *app.App
	sSize tea.WindowSizeMsg

	dialog tea.Model

	keys  keyMap
	table *stickers.TableSingleType[string]
	help  help.Model
}

func newMainScreenTable(app *app.App, sSize tea.WindowSizeMsg) mainScreenTable {

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

	return mainScreenTable{
		app:   app,
		sSize: sSize,
		keys:  keys,
		help:  helpModel,
		table: newTable(app),
	}
}

func (m mainScreenTable) Init() tea.Cmd {
	return nil
}

func (m mainScreenTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
		m.help.Width = msg.Width
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 1)
		return m, nil
	}

	switch msg := msg.(type) {
	case TuiDialogValue:
		peer := string(msg)
		_ = m.app.State.AddPeer(peer)
		m.dialog = nil
		m.table = newTable(m.app)
		m.table.SetWidth(m.sSize.Width)
		m.table.SetHeight(m.sSize.Height - 1)
		return m, nil

	case TuiDialogCancel:
		_ = msg
		m.dialog = nil
		return m, nil
	}

	if m.dialog != nil {
		w, c := m.dialog.Update(msg)
		m.dialog = w
		return m, c
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyF10:
			return m, tea.Quit

		case tea.KeyF7:
			d := NewTuiDialogName()
			d.Title = "Create a new Peer"
			d.Question = "Enter new peer name"
			d.ValidationFunc = func(s string) string {
				if s == "" {
					return "cannot be empty"
				}
				_, err := m.app.State.CanAddPeer(s)
				if err != nil {
					return err.Error()
				}
				return ""
			}
			m.dialog = d
			return m, d.Init()

		case tea.KeyEnter:
			/*
				return m, tea.Batch(
					tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
				)
			*/

		case tea.KeyDown:
			m.table.CursorDown()

		case tea.KeyUp:
			m.table.CursorUp()

		}
	}
	//m.table, cmd = m.table.
	return m, cmd
}

func (m mainScreenTable) View() string {
	if m.dialog != nil {
		return m.dialog.View()
	}

	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(lipgloss.Left, m.table.Render(), helpView)
}
