package app

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"os"
)

func (a *App) CreateNewServer(b backend.ServerBlueprint) error {

	p := a.interfaceDir(b.InterfaceName)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("interface directory exists %s", p)
	}

	err := os.Mkdir(p, 0700)
	if err != nil {
		return fmt.Errorf("can't create %s:%w", p, err)
	}

	err = os.Chdir(p)
	if err != nil {
		return fmt.Errorf("can't chdir %s:%w", p, err)
	}

	server := backend.NewServerWithBlueprint(b)
	err = server.WriteOnce()
	if err != nil {
		return fmt.Errorf("server write %w", err)
	}

	return nil
}
