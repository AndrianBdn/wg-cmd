package app

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/sysinfo"
)

type App struct {
	Settings *Settings
	State    *backend.State
}

func NewApp() *App {
	configureLogger()

	settings, err := readSettings()
	if err != nil {
		fmt.Println("Fatal error when reading settings", err)
		os.Exit(1)
	}

	settings.applyCommandLine()

	a := App{
		Settings: settings,
	}

	if settings.cliCommand == "update-systemd-unit" {
		// needs no interface state; skip loading and wizard checks so their
		// error exits can't shadow the command
		return &a
	}

	log.Println("Loading interface", a.Settings.DefaultInterface)
	err = a.LoadInterface(a.Settings.DefaultInterface)
	log.Println("error", err, "state", a.State != nil)

	if err == nil && a.State != nil && settings.saveInterface && settings.cliCommand == "" {
		err := a.SaveSettings()
		if err != nil {
			log.Println("Error saving settings", err)
		}
	}

	if err != nil {
		fmt.Printf("Unable to load interface %s: %s\n", a.Settings.DefaultInterface, err)
		fmt.Println("To troubleshoot:")
		fmt.Printf("1. Check if %s exists\n", a.interfaceDir(a.Settings.DefaultInterface))
		fmt.Printf("2. Launch \"%s new\"  to create new interface\n", os.Args[0])
		os.Exit(1)
	}

	if a.State == nil {
		// No interface loaded: we're about to launch the setup wizard, which
		// needs a writable database directory. If it isn't (e.g. a normal user
		// pointed at the default /etc/wireguard), show a clear error instead of
		// a wizard the user can't complete.
		if msg := a.wizardDirError(); msg != "" {
			fmt.Println(msg)
			os.Exit(1)
		}
	}

	return &a
}

func (a *App) RunCli() {
	if a.Settings.cliCommand == "update-systemd-unit" {
		iface := a.Settings.cliArg
		if !regexp.MustCompile(`^wg\d{1,4}$`).MatchString(iface) {
			fmt.Printf("Usage: %s update-systemd-unit <wg-interface>\n", os.Args[0])
			os.Exit(1)
		}
		if err := sysinfo.UpdateSystemdUnit(iface, a.Settings.WireguardDir); err != nil {
			fmt.Println("Error updating systemd unit:", err)
			os.Exit(1)
		}
		fmt.Printf("Updated wgc-%s.service: config changes now apply via wg syncconf (no interface restart)\n", iface)
		os.Exit(0)
	}

	if a.Settings.cliCommand == "make" {
		if a.State == nil {
			fmt.Printf("No interface selected. Run \"%s <wg-interface> make\"\n", os.Args[0])
			os.Exit(1)
		}
		path, err := a.GenerateWireguardConfig()
		if err != nil {
			fmt.Println("Error making WireGuard config:", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully wrote WireGuard config to %s\n", path)
		os.Exit(0)
	}
}

func (a *App) LoadInterface(ifName string) error {
	if ifName == "" {
		return nil
	}
	p := a.interfaceDir(ifName)
	if _, err := os.Stat(p); err != nil {
		return err
	}
	state, err := backend.ReadState(p, log.New(io.Discard, "", 0))
	if err == nil {
		a.State = state
	}
	return err
}

func (a *App) GenerateWireguardConfigLog() {
	_, err := a.GenerateWireguardConfig()
	if err != nil {
		// TODO: probably not the best way to handle this error
		log.Println("Error generating config", err)
	}
}

func (a *App) GenerateWireguardConfig() (string, error) {
	configPath := filepath.Join(a.Settings.WireguardDir, a.State.Server.Interface) + ".conf"
	return configPath, a.State.GenerateWireguardFile(configPath, false)
}

func (a *App) ValidateIfaceArg(ifName string) string {
	if !regexp.MustCompile(`^wg\d{1,4}$`).MatchString(ifName) {
		return "Interface name should be in form wg<number>"
	}

	p := filepath.Join(a.Settings.WireguardDir, ifName+".conf")
	if _, err := os.Stat(p); err == nil {
		return fmt.Sprintf("Found config for %s at %s. Try a different name.", ifName, a.Settings.WireguardDir)
	}

	p = a.interfaceDir(ifName)
	if _, err := os.Stat(p); err == nil {
		return fmt.Sprintf("Found directory %s at %s. Try a different name.",
			filepath.Base(p),
			a.Settings.WireguardDir)
	}

	if sysinfo.NetworkInterfaceExists(ifName) {
		return "Network interface exists in routing tables. Try a different name."
	}

	return ""
}

// wizardDirError returns a user-facing message when the setup wizard cannot
// run because its database directory is missing or not writable by the current
// user. It returns "" when the directory is usable.
func (a *App) wizardDirError() string {
	dir := a.Settings.DatabaseDir

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Sprintf("Directory %s does not exist.\n"+
			"Install WireGuard or create the directory before running setup.", dir)
	}
	if err != nil {
		return fmt.Sprintf("Can't access %s: %s", dir, err)
	}
	if !info.IsDir() {
		return fmt.Sprintf("%s is not a directory.", dir)
	}

	if w := testIfDirWritable(dir); w != "" {
		msg := fmt.Sprintf("Directory %s is not writable by the current user.", dir)
		if os.Geteuid() != 0 {
			msg += fmt.Sprintf("\nThis usually needs root — try: sudo %s", os.Args[0])
		}
		return msg
	}

	return ""
}

func (a *App) TestDirectories() string {
	dbTest := testIfDirWritable(a.Settings.DatabaseDir)

	if a.Settings.DatabaseDir == a.Settings.WireguardDir || dbTest != "" {
		return dbTest
	}

	return testIfDirWritable(a.Settings.WireguardDir)
}

func (a *App) interfaceDir(i string) string {
	d := "wgc-" + i
	return filepath.Join(a.Settings.DatabaseDir, d)
}

func testIfDirWritable(dir string) string {
	if _, err := os.Stat(dir); err != nil {
		return fmt.Sprint("can't stat", dir, err.Error())
	}

	testFileName := randomFileName()
	testFile := filepath.Join(dir, testFileName)

	err := os.WriteFile(testFile, []byte(testFileName), 0o600)
	if err != nil {
		return fmt.Sprint("can't write file in ", dir, ":", err.Error())
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
