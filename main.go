package main

import (
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"
)

func printUsageAndQuit() {
	fmt.Printf("Usage: %s <command>\n", os.Args[0])
	fmt.Println("\nAvailable commands:")
	fmt.Println("\tinit\tinit an empty wg-dir-conf directory")
	fmt.Println("\tmake\tcreate a wireguard config out of wg-dir-conf directory")
	fmt.Println("\talloc\tallocate new client and print its config")
	fmt.Println("\tprint\tprint client config (only possible when private key is saved)")
	os.Exit(1)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	syscall.Umask(0000)

	if len(os.Args) < 2 {
		printUsageAndQuit()
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])

	case "make":
		runMake(os.Args[2:])

	default:
		fmt.Printf("Error: unknown command '%s'\n", os.Args[1])
	}
}
