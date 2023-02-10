package main

import (
	"fmt"
	"github.com/andrianbdn/wg-cmd/app"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"os"
	"time"
)

var version string

func main() {
	rand.Seed(time.Now().UnixNano()) // TODO: Remove when upgrade to go 1.20

	a := app.NewApp()
	a.RunCli() // will do nothing if cli is not used
	m := NewAppModel(a)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
