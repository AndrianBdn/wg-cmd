package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type Settings struct {
	WireguardDir     string
	DatabaseDir      string
	DefaultInterface string
	ViewerQRMode     bool

	cliCommand    string
	saveInterface bool
}

func readSettings() (*Settings, error) {
	configPath, err := xdg.ConfigFile("wg-cmd/config.toml")
	if err != nil {
		return nil, fmt.Errorf("xdg %w", err)
	}

	s := defaultSettings()

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		// not existing config is not an error
		return &s, nil
	}

	_, err = toml.DecodeFile(configPath, &s)
	if err != nil {
		return nil, fmt.Errorf("config error %w", err)
	}

	return &s, nil
}

func (a *App) SaveSettings() error {
	configPath, err := xdg.ConfigFile("wg-cmd/config.toml")
	if err != nil {
		return fmt.Errorf("xdg %w", err)
	}

	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("SaveSettings OpenFile %w", err)
	}

	if err := toml.NewEncoder(f).Encode(a.Settings); err != nil {
		return fmt.Errorf("SaveSettings err toml encode %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("SaveSettings close %w", err)
	}

	return nil
}

func defaultSettings() Settings {
	dd := os.Getenv("WG_CMD_DIR")

	if dd == "" {
		dd = "/etc/wireguard"
	}
	return Settings{
		WireguardDir: dd,
		DatabaseDir:  dd,
	}
}

func (s *Settings) applyCommandLine() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	arg := args[0]
	if arg == "make" {
		s.cliCommand = "make"
		return
	}

	if arg == "new" {
		arg = ""
	} else if strings.HasPrefix(arg, "wgc-") {
		arg = strings.Replace(arg, "wgc-", "", 1)
		s.saveInterface = true
	}

	s.DefaultInterface = arg

	if len(args) > 1 && args[1] == "make" {
		s.cliCommand = "make"
	}
}
