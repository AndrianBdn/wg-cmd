package main

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"path/filepath"
)

type App struct {
	Settings *Settings
	logger   *log.Logger
	state    *backend.State

	dialog tea.Model
}

func NewApp() *App {
	a := App{
		Settings: readSettings(),
		logger:   log.New(os.Stderr, "", 0),
	}

	st, err := backend.ReadState(
		filepath.Join(a.Settings.DatabaseDir, a.Settings.DefaultInterface),
		a.logger,
	)

	if st != nil && err == nil {
		a.state = st
	}

	return &a
}

func main() {
	a := NewApp()
	m := NewAppModel(a)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
