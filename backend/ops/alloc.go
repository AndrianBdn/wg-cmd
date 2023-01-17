package ops

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"log"
	"regexp"
)

func OpAlloc(peerName string, wlog *log.Logger) (*backend.Client, error) {
	r := regexp.MustCompile(`^` + backend.PeerNameRegExp + `$`)
	if !r.MatchString(peerName) {
		return nil, fmt.Errorf("<peer_name> must start with letter, contain only letters, numbers, underscore, dash")
	}

	state, err := backend.ReadState(wlog)
	if err != nil {
		return nil, err
	}

	for _, excl := range state.Clients {
		if excl.GetName() == peerName {
			return nil, fmt.Errorf("peer name %s is already used by %s", peerName, excl.GetFileName())
		}
	}

	foundIP := -1
	for i := 2; i < 4096; i++ {
		if _, ok := state.Clients[i]; !ok {
			foundIP = i
			break
		}
	}

	if foundIP == -1 {
		return nil, fmt.Errorf("subnet depleted, all addresses in use")
	}

	c := backend.NewClient(foundIP, peerName)
	err = c.WriteOnce()
	if err != nil {
		return nil, fmt.Errorf("can't write %s: %w", c.GetFileName(), err)
	}

	return c, nil
}
