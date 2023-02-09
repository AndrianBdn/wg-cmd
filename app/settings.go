package app

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"os"
)

type Settings struct {
	WireguardDir     string
	DatabaseDir      string
	DefaultInterface string
	ViewerQRMode     bool
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

	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
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
	dd := os.Getenv("DEBUG_WG_CMD_DIR")

	if dd == "" {
		dd = "/etc/wireguard"
	}
	return Settings{
		WireguardDir: dd,
		DatabaseDir:  dd,
	}
}
