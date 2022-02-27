package utils

import (
	"fmt"
	"strings"

	"github.com/kyverno/policy-reporter-cli/pkg/model"
	"github.com/kyverno/policy-reporter-cli/pkg/policyreporter"
)

func GroupResultsByResult(reportResults []policyreporter.PolicyReportResult, results []string) []*model.Group {
	groups := make(map[policyreporter.Result]*model.Group, 0)
	for _, s := range results {
		groups[s] = &model.Group{
			Label: fmt.Sprintf("%s Policy Results", strings.Title(s)),
			List:  make([]policyreporter.PolicyReportResult, 0),
		}
	}

	for _, reportResult := range reportResults {
		if group, ok := groups[reportResult.Status]; ok {
			group.List = append(group.List, reportResult)
		}
	}

	list := []*model.Group{}

	for _, s := range results {
		if len(groups[s].List) == 0 {
			continue
		}

		list = append(list, groups[s])
	}

	return list
}

func GroupResultsByCategory(results []policyreporter.PolicyReportResult, categories []string) []*model.Group {
	groups := make(map[string]*model.Group, 0)
	if len(categories) == 0 {
		return []*model.Group{
			{
				Label: "No Category",
				List:  results,
			},
		}
	}

	for _, s := range categories {
		var label = s
		if label == "" {
			label = "No Category"
		}

		groups[s] = &model.Group{
			Label: label,
			List:  make([]policyreporter.PolicyReportResult, 0),
		}
	}

	for _, result := range results {
		if group, ok := groups[result.Category]; ok {
			group.List = append(group.List, result)
		}
	}

	result := []*model.Group{}

	for _, s := range categories {
		if len(groups[s].List) == 0 {
			continue
		}

		result = append(result, groups[s])
	}

	return result
}

func GroupResultsByPolicy(results []policyreporter.PolicyReportResult, policies []string) []*model.Group {
	groups := make(map[string]*model.Group, 0)

	for _, s := range policies {
		groups[s] = &model.Group{
			Label: s,
			List:  make([]policyreporter.PolicyReportResult, 0),
		}
	}

	for _, result := range results {
		if group, ok := groups[result.Policy]; ok {
			group.List = append(group.List, result)
		}
	}

	result := []*model.Group{}

	for _, s := range policies {
		if len(groups[s].List) == 0 {
			continue
		}

		result = append(result, groups[s])
	}

	return result
}

func GroupResultsByResource(results []policyreporter.PolicyReportResult) []*model.Group {
	groups := make(map[string]*model.Group, 0)

	genKey := func(result policyreporter.PolicyReportResult) string {
		return fmt.Sprintf("%s/%s/%s", result.Namespace, result.Kind, result.Name)
	}

	for _, result := range results {
		key := genKey(result)
		if group, ok := groups[key]; ok {
			group.List = append(group.List, result)
		} else {
			groups[key] = &model.Group{
				Label: fmt.Sprintf("%s %s", result.Kind, result.Name),
				List:  []policyreporter.PolicyReportResult{result},
			}
		}
	}

	result := []*model.Group{}

	for s := range groups {
		result = append(result, groups[s])
	}

	return result
}

func NoneGrouping(results []policyreporter.PolicyReportResult) []*model.Group {
	return []*model.Group{
		{
			Label: "",
			List:  results,
		},
	}
}
