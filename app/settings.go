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
