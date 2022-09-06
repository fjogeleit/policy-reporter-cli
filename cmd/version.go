package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCMD(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Client version of Policy Reporter CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Client Version: %s\n", version)
		},
	}

	return cmd
}
