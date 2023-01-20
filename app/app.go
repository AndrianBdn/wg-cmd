package app

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type App struct {
	Settings *Settings
	logger   *log.Logger
	state    *backend.State

	dialog tea.Model
}

func NewApp() *App {
	a := App{
		Settings: ReadSettings(),
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

func (app *App) ValidateIfaceArg(ifName string) string {
	if !regexp.MustCompile(`^wg\d{1,4}$`).MatchString(ifName) {
		return "Interface name should be in form wg<number>"
	}

	p := filepath.Join(app.Settings.WireguardDir, ifName+".conf")
	if _, err := os.Stat(p); err == nil {
		return fmt.Sprintf("Found existing config for %s at %s. Try a different name.", ifName, app.Settings.WireguardDir)
	}
	return ""
}
