package results

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

func grouingResults(ctx context.Context, results policyreporter.ResultList, api policyreporter.API, apiFilter policyreporter.Filter) []*model.Group {
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
		groups = utils.GroupResultsByCategory(results.Items, categories)
	case cli.PolicyGrouping:
		policies := apiFilter.Policies
		if len(policies) == 0 {
			policies, _ = api.Policies(ctx, apiFilter)
		}
		groups = utils.GroupResultsByPolicy(results.Items, policies)
	case cli.ResourceGrouping:
		groups = utils.GroupResultsByResource(results.Items)
	case cli.NoneGroup:
		groups = utils.NoneGrouping(results.Items)
	default:
		groups = utils.GroupResultsByResult(results.Items, result)
	}

	return groups
}

func buildTable(groups []*model.Group) {
	if len(groups) == 0 {
		fmt.Println("No results found")

		return
	}

	for _, group := range groups {
		if group.Label != "" && len(groups) > 1 {
			fmt.Println("")
			fmt.Println(chalk.Bold.TextStyle(group.Label))
			fmt.Println("")
		}

		prn, err := klo.PrinterFromFlag(output, &klo.Specs{
			WideColumnSpec:    "NAMESPACE:{.Namespace},KIND:{.Kind},NAME:{.Name},POLICY:{.Policy},RULE:{.Rule},SEVERITY:{.Severity},RESULT:{.Status}",
			DefaultColumnSpec: "NAMESPACE:{.Namespace},KIND:{.Kind},NAME:{.Name},POLICY:{.Policy},RULE:{.Rule},RESULT:{.Status}",
		})
		if err != nil {
			fmt.Println(err)
		}

		prn.Fprint(os.Stdout, group.List)
	}
}

func generateFilterFromFlags(currentNamespace string) policyreporter.Filter {
	filter := policyreporter.Filter{}

	if source != "" {
		filter.Sources = []string{source}
	}

	if namespace != "" {
		filter.Namespaces = []string{namespace}
	} else if allNamespaces {
		filter.Namespaces = []string{}
	} else if currentNamespace != "" {
		filter.Namespaces = []string{currentNamespace}
	}

	if len(results) != 0 {
		filter.Status = results
	}
	if len(categories) != 0 {
		filter.Categories = categories
	}
	if len(kinds) != 0 {
		filter.Kinds = kinds
	}
	if len(policies) != 0 {
		filter.Policies = policies
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

	if namespace == "" && !allNamespaces {
		options = append(options, "Namespace")
	}

	if len(policies) == 0 {
		options = append(options, "Policy")
	}

	if len(kinds) == 0 {
		options = append(options, "Kind")
	}

	options = append(options, "Resource", "Severity")

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
