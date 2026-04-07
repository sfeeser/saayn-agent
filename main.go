package main

import (
	"os"

	"github.com/sfeeser/saayn-agent/cmd/saayn"
)

func main() {
	if err := saayn.Execute(); err != nil {
		os.Exit(1)
	}
}
