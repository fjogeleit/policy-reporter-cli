package clusterresults

import (
	"context"
	"fmt"
	"os"

	"github.com/kyverno/policy-reporter-cli/pkg/cli"
	"github.com/kyverno/policy-reporter-cli/pkg/model"
	"github.com/kyverno/policy-reporter-cli/pkg/policyreporter"
	"github.com/kyverno/policy-reporter-cli/pkg/utils"
	"github.com/thediveo/klo"
	"github.com/ttacon/chalk"
)

func grouingResults(ctx context.Context, results []policyreporter.PolicyReportResult, api policyreporter.API, apiFilter policyreporter.Filter) []*model.Group {
	result := apiFilter.Status
	if len(result) == 0 {
		result = policyreporter.AllResults
	}

	var groups []*model.Group

	switch groupBy {
	case cli.CategoryGrouping:
		categories := apiFilter.Categories
		if len(categories) == 0 {
			categories, _ = api.Categories(ctx)
		}
		groups = utils.GroupResultsByCategory(results, categories)
	case cli.PolicyGrouping:
		policies := apiFilter.Policies
		if len(policies) == 0 {
			policies, _ = api.ClusterPolicies(ctx, apiFilter)
		}
		groups = utils.GroupResultsByPolicy(results, policies)
	case cli.ResourceGrouping:
		groups = utils.GroupResultsByResource(results)
	case cli.NoneGroup:
		groups = utils.NoneGrouping(results)
	default:
		groups = utils.GroupResultsByResult(results, result)
	}

	return groups
}

func buildTable(groups []*model.Group) {
	if len(groups) == 0 {
		fmt.Println("No results found")

		return
	}

	for _, group := range groups {
		if group.Label != "" {
			fmt.Println("")
			fmt.Println(chalk.Bold.TextStyle(group.Label))
			fmt.Println("")
		}

		prn, err := klo.PrinterFromFlag(output, &klo.Specs{
			WideColumnSpec:    "KIND:{.Kind},NAME:{.Name},POLICY:{.Policy},RULE:{.Rule},SEVERITY:{.Severity},RESULT:{.Status}",
			DefaultColumnSpec: "KIND:{.Kind},NAME:{.Name},POLICY:{.Policy},RULE:{.Rule},RESULT:{.Status}",
		})
		if err != nil {
			fmt.Println(err)
		}

		prn.Fprint(os.Stdout, group.List)
	}
}

func generateFilterFromFlags() policyreporter.Filter {
	filter := policyreporter.Filter{}

	if source != "" {
		filter.Sources = []string{source}
	}
	if len(results) != 0 {
		filter.Status = results
	}
	if len(categories) != 0 {
		filter.Categories = categories
	}

	return filter
}

func generateSearchOptionsFromFlags() []string {
	options := []string{}

	if source == "" {
		options = append(options, "Source")
	}

	if len(categories) == 0 {
		options = append(options, "Category")
	}

	options = append(options, "Policy", "Kind", "Severity")

	if len(results) == 0 {
		options = append(options, "Result")
	}

	return options
}

func preselect(values []string) interface{} {
	if len(values) == 1 {
		return values
	}

	return nil
}
