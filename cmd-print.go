package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/mdp/qrterminal/v3"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func runPrint(args []string) {

	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	printQR := printCmd.Bool("qr", false, "print QR code")
	err := printCmd.Parse(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	args = printCmd.Args()
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s [--qr] <ip number|file name|peer_name>\n", os.Args[0], os.Args[1])
		fmt.Println("\targument is either a file name, or an IP number, or a peer name")
		fmt.Println(dedent.Dedent(`
			This command will print peer config
			Use --qr flag to output QR code`))
		os.Exit(1)
	}

	st := readState()
	arg := args[0]

	c := clientTryNumber(arg, st)
	if c == nil {
		c = clientTryFileName(arg, st)
	}
	if c == nil {
		c = clientTryPeerName(arg, st)
	}

	if c == nil {
		fmt.Println("Error: could not find any client which looks like", arg)
		os.Exit(1)
	}

	buf := bytes.NewBuffer(nil)
	err = generateClientConfig(st.server, c, buf)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if *printQR == false {
		fmt.Println(clientCommentLong(st.server, c))
		fmt.Print(buf.String())
		fmt.Println("-- END OF CONFIG; DO NOT COPY-PASTE THIS LINE --")
	} else {
		printQRCode(c.fileName, buf.String())
	}

	fmt.Printf("- Tip: run `%s make` to re-generate real Wireguard config\n", os.Args[0])
}

func printQRCode(fileName, configText string) {
	fmt.Println("WireGuard QR Code for", fileName, ":")
	qrterminal.GenerateHalfBlock(configText, qrterminal.L, os.Stdout)
}

func clientTryNumber(arg string, st *state) *client {
	if !regexp.MustCompile(`^\d+$`).MatchString(arg) {
		return nil
	}
	ipNum, err := strconv.Atoi(arg)
	if err != nil {
		return nil
	}
	if ipNum < 2 {
		fmt.Println("Error: clients IP numbers start with 2")
		os.Exit(1)
	}
	return st.clients[ipNum]
}

func clientTryFileName(arg string, st *state) *client {
	if !regexp.MustCompile(`^(\d+)-`).MatchString(arg) {
		return nil
	}
	if !strings.HasSuffix(arg, ".toml") {
		arg = arg + ".toml"
	}
	for _, c := range st.clients {
		if c.fileName == arg {
			return c
		}
	}
	return nil
}

func clientTryPeerName(arg string, st *state) *client {
	if !regexp.MustCompile("^" + PeerNameRegExp + "$").MatchString(arg) {
		return nil
	}
	for _, c := range st.clients {
		if c.name == arg {
			return c
		}
	}
	return nil
}
