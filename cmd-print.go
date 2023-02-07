package main

import (
	"flag"
	"fmt"
	"github.com/andrianbdn/wg-cmd/backend/ops"
	"github.com/lithammer/dedent"
	"log"
	"os"
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

	arg := args[0]

	config, err := ops.OpPrint(arg, *printQR, log.New(os.Stdout, "", 0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Print(config)
	fmt.Fprintln(os.Stderr, "-- END OF CONFIG; DO NOT COPY-PASTE THIS LINE --")
	fmt.Fprintf(os.Stderr, "- Tip: run `%s make` to re-generate real Wireguard config\n", os.Args[0])
}
