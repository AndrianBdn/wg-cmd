package app

type Settings struct {
	WireguardDir     string
	DatabaseDir      string
	DefaultInterface string
}

func ReadSettings() *Settings {
	s := Settings{
		WireguardDir:     "/etc/wireguard",
		DatabaseDir:      ".",
		DefaultInterface: "wg0",
	}

	//configPath, err := xdg.ConfigFile("wg-cmd/config.toml")
	//fmt.Println(configPath, err)
	//
	//configPath, err = xdg.SearchConfigFile("wg-cmd/config.toml")
	//fmt.Println(configPath, err)
	return &s
}
