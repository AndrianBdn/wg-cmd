package main

import (
	"log"
	"os/exec"

	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/sysinfo"
	tea "github.com/charmbracelet/bubbletea"
)

type editingFinishedMsg struct {
	err error
}

type extEditorState struct {
	editor     string
	file       string
	reopen     bool
	editServer bool
}

func newExtEditorState(file string, editServer bool) extEditorState {
	return extEditorState{editor: sysinfo.GetSystemEditorPath(), file: file, editServer: editServer}
}

func (e *extEditorState) launchEditor() tea.Cmd {
	cmd := exec.Command(e.editor, e.file)
	log.Println("Launching editor:", cmd.String())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		log.Println("Editor finished with error:", err)
		return editingFinishedMsg{err}
	})
}

func (m MainScreen) editorEvents(msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editingFinishedMsg:
		if msg.err != nil {
			m.dialog = NewTuiDialogMsg("Error", msg.err.Error(), true)
			m.extEditor.reopen = false
			return true, m, m.dialog.Init()
		}

		err := m.app.LoadInterface(m.app.State.Server.Interface)
		if err != nil {
			m.dialog = NewTuiDialogMsg(
				"Error",
				err.Error()+". Edit file to fix the problem.",
				true)
			m.extEditor.reopen = true
			return true, m, m.dialog.Init()
		}
		m.extEditor.reopen = false
		if m.extEditor.editServer {
			m.table = newAppDynamicTableList(m.app, &m.table)
		}

		// Even if we edit peer, we may need to regenerate the server config
		// because of AddServerRoute client configuration option mainly
		m.app.GenerateWireguardConfigLog()

		return true, m, nil
	}

	return false, m, nil
}

func (m MainScreen) EditCurrentItem() (tea.Model, tea.Cmd) {
	file := ""
	editServer := false
	if m.table.GetSelectedIndex() == 0 {
		file = backend.ServerFileName
		editServer = true
	} else {
		peer := clientFromRow(m.app, m.table.GetSelected())
		file = peer.GetFileName()
	}

	file = m.app.State.JoinPath(file)
	m.extEditor = newExtEditorState(file, editServer)

	if m.extEditor.editor == "" {
		m.dialog = NewTuiDialogMsg("Error", "Cannot find any editor", true)
		return m, nil
	}

	m.exitBanner = exitBannerShouldShow
	return m, m.extEditor.launchEditor()
}
