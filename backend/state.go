package backend

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const PeerNameRegExp = `([A-Za-z][0-9A-Za-z-_]*)`

type State struct {
	Server  *Server
	Clients map[int]*Client
}

func (s *State) CanAddPeer(peerName string) (int, error) {
	r := regexp.MustCompile(`^` + PeerNameRegExp + `$`)
	if !r.MatchString(peerName) {
		return -1, fmt.Errorf("<peer_name> must start with letter, contain only letters, numbers, underscore, dash")
	}

	for _, excl := range s.Clients {
		if excl.GetName() == peerName {
			return -1, fmt.Errorf("peer name %s is already used by %s", peerName, excl.GetFileName())
		}
	}

	foundIP := -1
	for i := 2; i < 4096; i++ {
		if _, ok := s.Clients[i]; !ok {
			foundIP = i
			break
		}
	}

	if foundIP == -1 {
		return -1, fmt.Errorf("subnet depleted, all addresses in use")
	}

	return foundIP, nil
}

func (s *State) AddPeer(peerName string) error {

	foundIP, err := s.CanAddPeer(peerName)
	if err != nil {
		return err
	}

	c := NewClient(foundIP, peerName)
	err = c.WriteOnce()
	if err != nil {
		return fmt.Errorf("can't write %s: %w", c.GetFileName(), err)
	}

	s.Clients[foundIP] = c

	return nil
}

func (s *State) DeletePeer(idx int) error {
	cl := s.Clients[idx]
	if cl == nil {
		log.Println("peer not found", idx)
		log.Printf("%+v", s.Clients)

		return fmt.Errorf("peer not found")
	}

	log.Printf("%+v", cl)
	log.Println(cl.GetFileName())

	err := os.Remove(cl.GetFileName())
	if err != nil {
		return fmt.Errorf("can't remove %w", err)
	}
	delete(s.Clients, idx)
	return nil
}

func (s *State) GenerateWireguardFile(wgConfigPath string, backup bool) error {
	buf := bytes.NewBuffer(nil)
	err := generateWireguardConfig(s, buf)
	if err != nil {
		return fmt.Errorf("server config generation: %w", err)
	}

	if _, err := os.Stat(wgConfigPath); err == nil {
		if backup {
			backupFile := wgConfigPath + ".bak"
			if _, err := os.Stat(backupFile); err == nil {
				err = os.Remove(backupFile)
				if err != nil {
					return fmt.Errorf("removing .bak file: %w", err)
				}
			}
			err = os.Rename(wgConfigPath, backupFile)
			if err != nil {
				return fmt.Errorf("creating .bak file: %w", err)
			}
		}
	}

	err = os.WriteFile(wgConfigPath, buf.Bytes(), 0600)
	if err != nil {
		return fmt.Errorf("writing to %s file: %w", wgConfigPath, err)
	}

	return nil
}

func ReadState(dir string, wlog *log.Logger) (*State, error) {
	if _, err := os.Stat(filepath.Join(dir, ServerFileName)); os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot find %s in wgcmd interface directory %s", ServerFileName, dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to list files in %s directory: %w", dir, err)
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

		client, err := ReadClient(dir, f.Name(), ip, name)
		if err != nil {
			return nil, fmt.Errorf("cannot read file %s error %w", f.Name(), err)
		}
		cls[ip] = client
	}

	return &State{Server: srv, Clients: cls}, nil
}
