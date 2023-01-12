package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"os"
	"regexp"
	"strconv"
)

func runInit(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s <interface>\n", os.Args[0], os.Args[1])
		fmt.Println("\t<interface> - wireguard interface name wg0, wg1, etc")
		fmt.Println("")
		fmt.Println(" The command will create directory <interface>-wg-dir-conf")
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

	dir := iface + "-wg-dir-conf"
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

	serverHost := discoverIP()

	ip4 := randomIP4()
	fmt.Println("- Random IP4 private address: " + ip4)
	fmt.Println("- Generating random unique local IPv6 unicast address")
	ip6 := randomIP6()
	fmt.Println("  " + ip6)

	server := newServer(iface, ip4, ip6, serverHost)
	err = server.save()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func validateIfaceArg(iface string) bool {
	// does anybody need interfaces not starting with wg? more than wg9999?
	return regexp.MustCompile(`^wg\d{1,4}$`).MatchString(iface)
}

func randomIP4() string {
	return "10." + strconv.Itoa(mrand.Intn(256)) + "." + strconv.Itoa(mrand.Intn(256)) + ".1/24"
}

func randomIP6() string {
	// I don't think that RFC 4193 with SHA1 and machine id is relevant now, let's just read
	// cryptographically random bytes
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic("failed to read 5 random bytes for IP6" + err.Error())
	}
	bhex := hex.EncodeToString(b)
	return "fd" + bhex[0:2] + ":" + bhex[2:6] + ":" + bhex[6:10] + "::1/64"
}
