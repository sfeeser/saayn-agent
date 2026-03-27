package main

import (
	"os"

	"github.com/saayn-agent/cmd/saayn"
)

func main() {
	if err := saayn.Execute(); err != nil {
		os.Exit(1)
	}
}
