package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.2.0"

func newVersionCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Client version of Policy Reporter CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Client Version: %s\n", version)
		},
	}

	return cmd
}
