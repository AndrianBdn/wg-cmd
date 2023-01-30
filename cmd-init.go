package main

import (
	"fmt"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"github.com/andrianbdn/wg-dir-conf/sysinfo"
	"os"
	"regexp"
)

func runInit(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s <interface>\n", os.Args[0], os.Args[1])
		fmt.Println("\t<interface> - wireguard interface name wg0, wg1, etc")
		fmt.Println("")
		fmt.Println(" The command will create directory wdc-<interface>")
		fmt.Println(" inside current working directory.")
		fmt.Println("")
		fmt.Println(" Subsequent commands must be run from that directory.")
		os.Exit(1)
	}

	if !validateIfaceArg(args[0]) {
		fmt.Println("Error: interface name be: wg<number>")
		os.Exit(1)
	}
	iface := args[0]

	dir := "wdc-" + iface
	if _, err := os.Stat(dir); err == nil {
		fmt.Printf("Error: directory %s exists in cwd\n", dir)
		fmt.Printf("       remove it, or choose other interface name\n")
		os.Exit(1)
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: we are unable to get current working directory:", err)
		os.Exit(1)
	}
	fmt.Printf("- Creating %s in %s\n", dir, wd)

	err = os.Mkdir(dir, 0700)
	if err != nil {
		fmt.Println("Error: can't create", dir, ":", err)
		os.Exit(1)
	}

	err = os.Chdir(dir)
	if err != nil {
		fmt.Println("Error: can't chdir to", dir, ":", err)
		os.Exit(1)
	}

	serverHost := sysinfo.DiscoverIPOld()

	server := backend.NewServer(iface, serverHost)
	err = server.WriteOnce()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("- Config %s written successfuly to %s\n", backend.ServerFileName, dir)
	fmt.Printf("  Feel free to inspect and change this file\n")
	fmt.Printf("- Tip: cd to %s and run `%s alloc <peer-name>` to add peers\n", dir, os.Args[0])
}

func validateIfaceArg(iface string) bool {
	// does anybody need interfaces not starting with wg? - more than wg9999?
	return regexp.MustCompile(`^wg\d{1,4}$`).MatchString(iface)
}
