package main

import (
	"log"

	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/theme"
	"github.com/andrianbdn/wg-cmd/tutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainScreen struct {
	app              *app.App
	sSize            tea.WindowSizeMsg
	dialog           tea.Model
	dialogFullScreen bool
	table            DynamicTableList
	helpKeys         []helpKey
	extEditor        extEditorState
	exitBanner       int
}

const (
	exitBannerNone = iota
	exitBannerShouldShow
	exitBannerShow
)

func NewMainScreen(app *app.App, sSize tea.WindowSizeMsg) MainScreen {
	helpKeys := []helpKey{
		{key: "F1", help: "About"},
		{key: "F4", help: "Edit"},
		{key: "F7", help: "Add Peer"},
		{key: "F8", help: "Delete Peer"},
		{key: "F10", help: "Quit"},
	}

	m := MainScreen{
		app:      app,
		sSize:    sSize,
		helpKeys: helpKeys,
		table:    newAppDynamicTableList(app, nil),
	}
	return m.SetSize(sSize)
}

func (m MainScreen) Init() tea.Cmd {
	return nil
}

func (m MainScreen) SetSize(sSize tea.WindowSizeMsg) MainScreen {
	m.sSize = sSize
	m.table.SetTableSize(sSize, 0, -1)
	return m
}

func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m = m.SetSize(msg)
	}

	if edRet, m, cmd := m.editorEvents(msg); edRet {
		return m, cmd
	}

	switch msg := msg.(type) {

	case TuiDialogMsgResult:
		m.dialogFullScreen = false
		m.dialog = nil
		if m.extEditor.reopen {
			return m, m.extEditor.launchEditor()
		}

	case TuiDialogYes:
		return m.ReallyDeletePeer(), nil

	case TuiDialogNo:
		m.dialog = nil

	case TuiDialogValue:
		peer := string(msg)
		return m.ReallyAddPeer(peer), nil

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

		case tea.KeyF1:
			return m.Help()

		case tea.KeyF4:
			return m.EditCurrentItem()

		case tea.KeyF10:
			if m.exitBanner == exitBannerShouldShow {
				// this is a temporary workaround for screen not being cleared / restored
				// after launching an external editor and quitting
				m.exitBanner = exitBannerShow
				return m, tea.Sequence(tea.ExitAltScreen,
					tea.Println("Thank you for using WG Commander!"),
					tea.Quit)
			}
			return m, tea.Quit

		case tea.KeyF7:
			return m.CreatePeer()

		case tea.KeyF8:
			return m.DeletePeer()

		case tea.KeyEnter:
			return m.ViewPeer()

		case tea.KeyDown:
			m.table.Down()

		case tea.KeyPgUp:
			m.table.PageUp()

		case tea.KeyUp:
			m.table.Up()

		case tea.KeyPgDown:
			m.table.PageDown()
		}
	}

	return m, cmd
}

func (m MainScreen) View() string {
	if m.exitBanner == exitBannerShow {
		return ""
	}
	mainScreen := lipgloss.JoinVertical(lipgloss.Left,
		m.table.Render(),
		RenderHelpLine(m.sSize.Width, m.helpKeys...),
	)

	if m.dialog != nil {
		if m.dialogFullScreen {
			return m.dialog.View()
		}
		return tutils.PlaceDialog(m.dialog.View(), mainScreen, m.sSize, theme.Current.MainTableDimmed)
	}

	return mainScreen
}

func (m MainScreen) ViewPeer() (tea.Model, tea.Cmd) {
	row := m.table.GetSelectedIndex()
	var p *backend.Client
	if row != 0 {
		p = clientFromRow(m.app, m.table.GetSelected())
		if p == nil {
			return m, nil
		}
	}

	m.dialogFullScreen = true
	m.dialog = NewViewPeer(m.sSize, m.app, p)
	return m, m.dialog.Init()
}

func (m MainScreen) CreatePeer() (tea.Model, tea.Cmd) {
	d := NewTuiDialogName()
	d.Title = "Create a new Peer"
	d.Question = "Enter a new peer name"
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
}

func (m MainScreen) DeletePeer() (tea.Model, tea.Cmd) {
	row := m.table.GetSelectedIndex()
	if row == 0 {
		m.dialog = NewTuiDialogMsg("Error", "Deleting the server is currently not supported", true)
		return m, m.dialog.Init()
	}

	peer := clientFromRow(m.app, m.table.GetSelected())
	if peer == nil {
		return m, nil
	}

	m.dialog = NewTuiDialogYesNo("Delete", "Delete peer \""+peer.GetName()+"\"?")
	return m, m.dialog.Init()
}

func (m MainScreen) ReallyAddPeer(name string) MainScreen {
	err := m.app.State.AddPeer(name)
	if err != nil {
		log.Println("Error adding peer", err)
	}
	if err == nil {
		m.app.GenerateWireguardConfigLog()
	}
	m.dialog = nil
	m.table = newAppDynamicTableList(m.app, &m.table)
	return m
}

func (m MainScreen) ReallyDeletePeer() MainScreen {
	row := m.table.GetSelectedIndex()
	if row == 0 {
		panic("we don't delete server")
	}

	peer := clientFromRow(m.app, m.table.GetSelected())
	if peer != nil {
		err := m.app.State.DeletePeer(peer.GetIPNumber())
		if err != nil {
			log.Println("Error deleting peer", err)
		}
		if err == nil {
			m.app.GenerateWireguardConfigLog()
		}

		m.table.DeleteSelectedRow()
	}
	m.dialog = nil
	return m
}

func (m MainScreen) Help() (tea.Model, tea.Cmd) {
	ver := Version
	if ver == "" {
		ver = "<not set>"
	}

	if Build != "" {
		ver += " (" + Build + ")"
	}

	m.dialog = NewTuiDialogMsg(
		"WG Commander",
		"version "+ver+"\n\n"+
			"(c) 2023-2024 by Andrian Budantsov\n"+
			"              and GitHub contributors\n\n"+
			theme.Current.DialogButtonFocus.Render("https://github.com/andrianbdn/wg-cmd")+
			"\n\n"+
			"Comes with ABSOLUTELY NO WARRANTY, distributed under the terms of\nthe MIT license.",
		false)
	return m, m.dialog.Init()
}
