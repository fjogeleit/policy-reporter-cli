package results

import "github.com/spf13/cobra"

var (
	allNamespaces bool

	namespace  string
	source     string
	output     string
	groupBy    string
	results    []string
	categories []string
	kinds      []string
)

func sharedFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request")
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, search results across all namespaces.")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output format. One of: yaml|json|wide|go-template|jsonpath")

	cmd.Flags().StringVarP(&source, "source", "s", "", "Filter PolicyReportResults by source")
	cmd.Flags().StringArrayVar(&results, "result", []string{}, "Filter PolicyReportResults by result")
	cmd.Flags().StringArrayVarP(&kinds, "kind", "k", []string{}, "Filter PolicyReportResults by kinds (only fullqualified singular kind names are supported)")
	cmd.Flags().StringArrayVar(&categories, "category", []string{}, "Filter PolicyReportResults by category")
	cmd.Flags().StringVar(&groupBy, "group-by", "result", "Group PolicyReportResults by result, category, resource, none")

	return cmd
}
