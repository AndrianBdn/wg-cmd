package main

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend/ops"
	"log"
	"os"
)

func runMake(args []string) {
	s, err := ops.OpMake(log.New(os.Stdout, "", 0))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfuly wrote Wireguard config to %s\n", s.ServerConfigPath)
	fmt.Printf("- Tip: run `systemctl restart wg-quick@%s.service` to restart Wireguard\n", s.Interface)
}
