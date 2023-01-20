package main

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/app"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	a := app.NewApp()
	m := NewAppModel(a)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
