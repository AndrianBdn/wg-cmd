package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const PeerNameRegExp = `([A-Za-z][0-9A-Za-z-_]*)`

type state struct {
	server  *server
	clients map[int]*client
}

func readState() *state {
	if _, err := os.Stat(serverFileName); os.IsNotExist(err) {
		fmt.Println("Error: cannot find", serverFileName, "in current directory")
		fmt.Println("       cd wdc-wg<num> first")
		os.Exit(1)
	}

	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("Error: unable to list files in current directory:", err)
		os.Exit(1)
	}

	srv, err := readServer()
	if err != nil {
		fmt.Println("Error: can't read", serverFileName, "error:", err)
		os.Exit(1)
	}

	cls := make(map[int]*client)

	r := regexp.MustCompile(`^(\d+)-` + PeerNameRegExp + `\.toml$`)

	for _, f := range files {
		if f.Name() == serverFileName {
			continue
		}
		m := r.FindAllStringSubmatch(f.Name(), -1)
		if len(m) == 0 {
			fmt.Println("Warning: unknown file", f.Name(), "skipping")
		}
		if len(m[0]) != 3 {
			panic("regexp error in parsing file name" + f.Name())
		}
		ip, err := strconv.Atoi(m[0][1])
		if err != nil {
			panic("unlikely logical error - number from regexp cannot be parsed")
		}
		if ip > 254 {
			fmt.Println("Error: at the moment wg-dir-conf only supports 253 peers")
			os.Exit(1)
		}
		if _, ok := cls[ip]; ok {
			fmt.Println("Error: there are at least two conflicting files with the same IP number")
			fmt.Println("      ", cls[ip], "and", cls[ip].fileName)
			os.Exit(1)
		}

		name := m[0][2]

		client, err := readClient(f.Name(), ip, name)
		if err != nil {
			fmt.Println("Error: cannot read file", f.Name(), "error", err)
			os.Exit(1)
		}
		cls[ip] = client
	}

	return &state{server: srv, clients: cls}
}
