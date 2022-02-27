package clusterresults

import "github.com/spf13/cobra"

var (
	namespace  string
	source     string
	output     string
	groupBy    string
	results    []string
	categories []string
)

func sharedFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output format. One of: yaml|json|wide|go-template|jsonpath")

	cmd.Flags().StringVarP(&source, "source", "s", "", "Filter PolicyReportResults by source")
	cmd.Flags().StringArrayVar(&results, "result", []string{}, "Filter PolicyReportResults by result")
	cmd.Flags().StringArrayVar(&categories, "category", []string{}, "Filter PolicyReportResults by category")
	cmd.Flags().StringVar(&groupBy, "group-by", "result", "Group PolicyReportResults by result, category, resource, none")

	return cmd
}
