package cmd

import (
	"github.com/kyverno/policy-reporter-cli/cmd/results"
	"github.com/spf13/cobra"
)

// NewResultsCLI creates a new instance of the results CLI
func newResultsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "results",
		Aliases: []string{"res"},
		Short:   "Interact with the namespace scoped Policy Reporter APIs",
	}

	cmd.AddCommand(results.NewListCMD())
	cmd.AddCommand(results.NewSearchCMD())

	return cmd
}
