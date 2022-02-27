package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kyverno/policy-reporter-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var output string

func newTargetsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "targets",
		Short:   "List configured Policy Reporter Targets",
		Aliases: []string{"tar"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			resolver := config.NewResolver(config.LoadConfig())

			conn, err := resolver.ForwardPolicyReporter(ctx)
			if err != nil {
				return err
			}
			defer conn.Close()

			api := resolver.API(conn.Port)
			targets, err := api.Targets(ctx)
			if err != nil {
				return err
			}

			if len(targets) == 0 {
				fmt.Println("No targets are configured")
				return nil
			}

			prn, err := klo.PrinterFromFlag(output, &klo.Specs{
				DefaultColumnSpec: "TARGET:{.Name},MINIMUM PRIORITY:{.MinimumPriority},SKIP EXISTING ON STARTUP:{.SkipExistingOnStartup},SOURCE:{.Source}",
			})
			if err != nil {
				fmt.Println(err)
			}

			return prn.Fprint(os.Stdout, targets)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output Format")

	return cmd
}
