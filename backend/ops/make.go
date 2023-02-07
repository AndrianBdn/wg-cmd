package ops

import (
	"github.com/andrianbdn/wg-cmd/backend"
	"log"
)

func OpMake(wlog *log.Logger) (*backend.Server, error) {
	state, err := backend.ReadState(".", wlog)
	if err != nil {
		return nil, err
	}

	err = state.GenerateWireguardFile(state.Server.ServerConfigPath, true)
	if err != nil {
		return nil, err
	}

	return state.Server, nil
}
