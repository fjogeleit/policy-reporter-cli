package cmd

import (
	"flag"

	"github.com/spf13/cobra"
)

// NewCLI creates a new instance of the root CLI
func NewCLI(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pr",
		Short: "CLI for Policy Reporter REST API",
		Long:  `Query information from the kyverno/policy-reporter REST API about (Cluster)PolicyReports`,
	}

	rootCmd.AddCommand(newTargetsCMD())
	rootCmd.AddCommand(newResultsCMD())
	rootCmd.AddCommand(newClusterResultsCMD())
	rootCmd.AddCommand(newVersionCMD(version))

	flag.Parse()

	return rootCmd
}
