package clusterresults

import (
	"context"

	"github.com/kyverno/policy-reporter-cli/pkg/config"
	"github.com/spf13/cobra"
)

var labels string

func NewListCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ClusterPolicyReportResults",
		RunE: func(command *cobra.Command, args []string) error {
			ctx := context.Background()
			resolver := config.NewResolver(config.LoadConfig())

			conn, err := resolver.ForwardPolicyReporter(ctx)
			if err != nil {
				return nil
			}
			defer conn.Close()

			api := resolver.API(conn.Port)
			filter := generateFilterFromFlags()
			results, err := api.ClusterResults(ctx, filter)
			if err != nil {
				return err
			}

			if labels != "" {
				k8sClient, err := resolver.K8sClient()
				if err == nil {
					results = k8sClient.LabelFilter(ctx, results, labels)
				}
			}

			buildTable(grouingResults(ctx, results, api, filter))

			return nil
		},
	}

	cmd.Flags().StringVarP(&labels, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")

	return sharedFlags(cmd)
}
