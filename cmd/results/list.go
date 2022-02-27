package results

import (
	"context"

	"github.com/kyverno/policy-reporter-cli/pkg/config"
	"github.com/spf13/cobra"
)

func NewListCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List PolicyReportResults",
		RunE: func(command *cobra.Command, args []string) error {
			ctx := context.Background()
			resolver := config.NewResolver(config.LoadConfig())

			conn, err := resolver.ForwardPolicyReporter(ctx)
			if err != nil {
				return nil
			}
			defer conn.Close()

			api := resolver.API(conn.Port)

			ns, err := resolver.CurrentNamespace()
			if err != nil {
				return err
			}

			filter := generateFilterFromFlags(ns)
			results, err := api.Results(ctx, filter)
			if err != nil {
				return err
			}

			buildTable(grouingResults(ctx, results, api, filter))

			return nil
		},
	}

	return sharedFlags(cmd)
}
