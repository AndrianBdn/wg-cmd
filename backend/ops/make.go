package ops

import (
	"bytes"
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"log"
	"os"
)

func OpMake(wlog *log.Logger) (*backend.Server, error) {
	state, err := backend.ReadState(".", wlog)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = backend.GenerateServerConfig(state, buf)
	if err != nil {
		return nil, fmt.Errorf("server config generation: %w", err)
	}

	if _, err := os.Stat(state.Server.ServerConfigPath); err == nil {
		wlog.Println("- Found existing file", state.Server.ServerConfigPath, "making .bak")
		backupFile := state.Server.ServerConfigPath + ".bak"
		if _, err := os.Stat(backupFile); err == nil {
			wlog.Println("- Removing previous .bak file")
			err = os.Remove(backupFile)
			if err != nil {
				return nil, fmt.Errorf("removing .bak file: %w", err)
			}
		}
		err = os.Rename(state.Server.ServerConfigPath, backupFile)
		if err != nil {
			return nil, fmt.Errorf("creating .bak file: %w", err)
		}
	}

	err = os.WriteFile(state.Server.ServerConfigPath, buf.Bytes(), 0600)
	if err != nil {
		return nil, fmt.Errorf("writing to %s file: %w", state.Server.ServerConfigPath, err)
	}

	return state.Server, nil
}
