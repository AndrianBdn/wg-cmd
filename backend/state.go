package backend

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const PeerNameRegExp = `([A-Za-z][0-9A-Za-z-_]*)`

type State struct {
	Server  *Server
	Clients map[int]*Client
}

func ReadState(wlog *log.Logger) (*State, error) {
	if _, err := os.Stat(ServerFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot find %s in current directory", ServerFileName)
	}

	files, err := os.ReadDir("./")
	if err != nil {
		return nil, fmt.Errorf("unable to list files in current directory: %w", err)
	}

	srv, err := ReadServer()
	if err != nil {
		return nil, fmt.Errorf("can't read %s error: %w", ServerFileName, err)
	}

	cls := make(map[int]*Client)

	r := regexp.MustCompile(`^(\d+)-` + PeerNameRegExp + `\.toml$`)

	for _, f := range files {
		if f.Name() == ServerFileName {
			continue
		}
		m := r.FindAllStringSubmatch(f.Name(), -1)
		if len(m) == 0 || len(m[0]) != 3 {
			if wlog != nil {
				wlog.Println("Warning: unknown file", f.Name(), "skipping")
			}
			continue
		}
		ip, err := strconv.Atoi(m[0][1])
		if err != nil {
			panic("unlikely logical error - number from regexp cannot be parsed")
		}
		if ip > 4094 {
			return nil, fmt.Errorf("at the moment wg-dir-conf only supports 4094 peers (/20 subnet)")
		}
		if _, ok := cls[ip]; ok {
			return nil, fmt.Errorf("conflicting files %s and %s", f.Name(), cls[ip].fileName)
		}

		name := m[0][2]

		client, err := ReadClient(f.Name(), ip, name)
		if err != nil {
			return nil, fmt.Errorf("cannot read file %s error %w", f.Name(), err)
		}
		cls[ip] = client
	}

	return &State{Server: srv, Clients: cls}, nil
}
