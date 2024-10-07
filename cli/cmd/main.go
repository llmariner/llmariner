package main

import (
	"log"
	"os"

	"github.com/llmariner/llmariner/cli/internal/root"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
