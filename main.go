package main

import (
	"fmt"
	"os"

	"github.com/0xJacky/Nginx-UI/cmd"
)

// Version information set by build flags
var (
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "unknown"
)

func main() {
	// Inject version info into cmd package
	cmd.Version = Version
	cmd.BuildDate = BuildDate
	cmd.Commit = Commit

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
