package cmd

import (
	"github.com/kyverno/policy-reporter-cli/cmd/clusterresults"
	"github.com/spf13/cobra"
)

// newClusterResultsCMD creates a new instance of the results CLI
func newClusterResultsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster-results",
		Aliases: []string{"cres"},
		Short:   "Interact with the cluster scoped Policy Reporter APIs",
	}

	cmd.AddCommand(clusterresults.NewListCMD())
	cmd.AddCommand(clusterresults.NewSearchCMD())

	return cmd
}
