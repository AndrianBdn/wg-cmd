package main

import (
	"fmt"
	"github.com/lithammer/dedent"
	"os"
	"regexp"
)

func runAlloc(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s <peer_name>\n", os.Args[0], os.Args[1])
		fmt.Println("\t<peer_name> - human readable peer name")
		fmt.Println("")
		fmt.Println(dedent.Dedent(`
			The command will find unused IP address, generate keys for client (peer)
			and writeOnce peers wg-dir-conf config in file <ip_num>-<peer_name>.toml
			`))
		os.Exit(1)
	}

	peerName := args[0]

	r := regexp.MustCompile(`^` + PeerNameRegExp + `$`)
	if !r.MatchString(peerName) {
		fmt.Println("Error: <peer_name> must contain only letters, numbers, _, -; it should start with letter")
		os.Exit(1)
	}

	state := readState()

	for _, excl := range state.clients {
		if excl.name == peerName {
			fmt.Println("Error: peer name", "'"+peerName+"'", "is already used by", excl.fileName)
			fmt.Println(excl.allowedIps(state.server))
			os.Exit(1)
		}
	}

	foundIP := -1
	for i := 2; i < 255; i++ {
		if _, ok := state.clients[i]; !ok {
			foundIP = i
			break
		}
	}

	if foundIP == -1 {
		fmt.Println("Error: subnet depleted, all addresses in use")
	}

	fmt.Printf("- Found IP address .%d for %s\n", foundIP, peerName)
	c := newClient(foundIP, peerName)
	err := c.writeOnce()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("- Config %s written successfuly\n", c.fileName)
	fmt.Println("- Warning: we saved client's PrivateKey in a config. This is useful.")
	fmt.Println("           You can still remove it from from", c.fileName, "after printing client config.\n"+
		"           (replace PrivateKey string value with \"redacted\")")
	fmt.Printf("- Tip: run `%s print %d` to print config to console\n", os.Args[0], c.ipNum)
	fmt.Printf("       run `%s print --qr %d` to print QR code\n", os.Args[0], c.ipNum)
}
