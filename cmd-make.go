package main

import (
	"bytes"
	"fmt"
	"os"
)

func runMake(args []string) {
	state := readState()
	buf := bytes.NewBuffer(nil)
	err := generateServerConfig(state, buf)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if _, err := os.Stat(state.server.ServerConfigPath); err != nil {
		fmt.Println("- Found existing file", state.server.ServerConfigPath, "making .bak")
		backupFile := state.server.ServerConfigPath + ".bak"
		if _, err := os.Stat(backupFile); err != nil {
			fmt.Println("- Removing previous .bak file")
			err = os.Remove(backupFile)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Check that it is possible to delete ")
				os.Exit(1)
			}
		}
		err = os.Rename(state.server.ServerConfigPath, backupFile)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Check that directory is writable")
			os.Exit(1)
		}
	}

	err = os.WriteFile(state.server.ServerConfigPath, buf.Bytes(), 0600)
	if err != nil {
		fmt.Printf("Error while writing to %s: %s\n", state.server.ServerConfigPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfuly wrote Wireguard config to %s\n", state.server.ServerConfigPath)
	fmt.Printf("- Tip: run `systemctl restart wg-quick@%s.service` to restart Wireguard\n", state.server.Interface)
}
