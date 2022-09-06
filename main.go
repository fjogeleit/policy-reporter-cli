package main

import (
	"fmt"
	"os"

	"github.com/kyverno/policy-reporter-cli/cmd"
)

var Version = "development"

func main() {
	if err := cmd.NewCLI(Version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
