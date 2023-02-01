package app

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"github.com/andrianbdn/wg-dir-conf/sysinfo"
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
	settings, err := readSettings()
	if err != nil {
		fmt.Println("Fatal error when reading settings", err)
		os.Exit(1)
	}

	a := App{
		Settings: settings,
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
		return fmt.Sprintf("Found config for %s at %s. Try a different name.", ifName, app.Settings.WireguardDir)
	}

	p = app.interfaceDir(ifName)
	if _, err := os.Stat(p); err == nil {
		return fmt.Sprintf("Found directory %s at %s. Try a different name.",
			filepath.Base(p),
			app.Settings.WireguardDir)
	}

	if sysinfo.NetworkInterfaceExists(ifName) {
		return fmt.Sprintf("Network interface exists in routing tables. Try a different name.")
	}

	return ""
}

func (app *App) TestDirectories() string {
	dbTest := testDir(app.Settings.DatabaseDir)

	if app.Settings.DatabaseDir == app.Settings.WireguardDir || dbTest != "" {
		return dbTest
	}

	return testDir(app.Settings.WireguardDir)
}

func (app *App) interfaceDir(i string) string {
	d := "wgc-" + i
	return filepath.Join(app.Settings.DatabaseDir, d)
}

func testDir(dir string) string {
	if _, err := os.Stat(dir); err != nil {
		return fmt.Sprint("can't stat", dir, err.Error())
	}

	testFileName := randomFileName()
	testFile := filepath.Join(dir, testFileName)

	err := os.WriteFile(testFile, []byte(testFileName), 0600)
	if err != nil {
		return fmt.Sprint("can't write file in ", dir, err.Error())
	}

	rtest, err := os.ReadFile(testFile)
	if err != nil {
		return fmt.Sprint("can't read file ", testFileName, " in ", dir, err.Error())
	}

	if testFileName != string(rtest) {
		return fmt.Sprint("what we read from ", testFile, " is not what we wrote")
	}

	err = os.Remove(testFile)
	if err != nil {
		return fmt.Sprint("can't delete file ", testFileName, " in ", dir, err.Error())
	}

	return ""
}

func randomFileName() string {
	b := make([]byte, 15)
	if _, err := rand.Read(b); err != nil {
		panic("failed to read random bytes to test write-ability" + err.Error())
	}
	bhex := base32.StdEncoding.EncodeToString(b)
	return bhex + ".test"
}
