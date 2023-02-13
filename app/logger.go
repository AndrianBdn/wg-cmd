package app

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func configureLogger() {
	if os.Getenv("WG_CMD_DEBUG_LOG") == "" {
		log.SetOutput(io.Discard)
		return
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	f, err := os.OpenFile(filepath.Join(exPath, "wg-cmd.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
}
