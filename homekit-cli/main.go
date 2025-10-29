package main

import (
	"log"

	"github.com/homekit/homekit-cli/cmd/homekit"
)

func main() {
	if err := homekit.Execute(); err != nil {
		log.Fatalf("homekit-cli: %v", err)
	}
}
