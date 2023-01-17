package main

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend/ops"
	"github.com/lithammer/dedent"
	"log"
	"os"
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

	c, err := ops.OpAlloc(peerName, log.New(os.Stdout, "", 0))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("- Config %s written successfuly\n", c.GetFileName())
	fmt.Println("- Warning: we saved client's PrivateKey in a config. This is useful.")
	fmt.Println("           You can still remove it from from", c.GetFileName(), "after printing client config.")
	fmt.Println("           (replace PrivateKey string value with \"redacted\")")
	fmt.Printf("- Tip: run `%s print %d` to print config to console\n", os.Args[0], c.GetIPNumber())
	fmt.Printf("       run `%s print --qr %d` to print QR code\n", os.Args[0], c.GetIPNumber())
}
