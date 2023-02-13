package app

import (
	"fmt"
	"os"

	"github.com/andrianbdn/wg-cmd/backend"
)

func (a *App) CreateNewServer(b backend.ServerBlueprint) error {
	p := a.interfaceDir(b.InterfaceName)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("interface directory exists %s", p)
	}

	err := os.Mkdir(p, 0o700)
	if err != nil {
		return fmt.Errorf("can't create %s:%w", p, err)
	}

	server := backend.NewServerWithBlueprint(b)
	err = server.WriteOnce(p)
	if err != nil {
		return fmt.Errorf("server write %w", err)
	}

	return nil
}
