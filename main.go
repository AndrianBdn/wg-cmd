package main

import (
	"fmt"
	"os"

	"github.com/andrianbdn/wg-cmd/app"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Version string
	Build   string
)

func main() {
	a := app.NewApp()
	a.RunCli() // will do nothing if cli is not used
	m := NewAppModel(a)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
