package results

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kyverno/policy-reporter-cli/pkg/config"
	"github.com/kyverno/policy-reporter-cli/pkg/policyreporter"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

type SurveyFunc = func(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error

var (
	surveys = map[string]SurveyFunc{
		"Source":    selectSources,
		"Category":  selectCategories,
		"Kind":      selectKinds,
		"Resource":  selectResources,
		"Policy":    selectPolicies,
		"Result":    selectResult,
		"Severity":  selectSeverity,
		"Namespace": selectNamespaces,
	}

	ErrNoResult = errors.New("No results")
)

func NewSearchCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search PolicyReportResults",
		RunE: func(command *cobra.Command, args []string) error {
			ctx := context.Background()

			resolver := config.NewResolver(config.LoadConfig())

			conn, err := resolver.ForwardPolicyReporter(ctx)
			if err != nil {
				return nil
			}
			defer conn.Close()

			api := resolver.API(conn.Port)

			apiFilter := generateFilterFromFlags("")
			filters := []string{}

			prompt := &survey.MultiSelect{
				Message:  "Search Results by:",
				Options:  generateSearchOptionsFromFlags(),
				PageSize: 10,
			}

			err = survey.AskOne(prompt, &filters)
			if err == terminal.InterruptErr {
				fmt.Println("")
				fmt.Println(chalk.Red, chalk.Bold.TextStyle("Search interrupted"))
				fmt.Println("")
				return nil
			}

			for _, filter := range filters {
				if surveyFunc, ok := surveys[filter]; ok {
					err = surveyFunc(ctx, api, &apiFilter)
					if err == terminal.InterruptErr {
						fmt.Println("")
						fmt.Println(chalk.Red, chalk.Bold.TextStyle("Search interrupted"))
						fmt.Println("")
						return nil
					} else if err == ErrNoResult {
						fmt.Println("No results found")
					} else if err != nil {
						return err
					}
				}
			}

			results, err := api.Results(ctx, apiFilter)
			if err != nil {
				return err
			}

			buildTable(grouingResults(ctx, results, api, apiFilter))

			return nil
		},
	}

	return sharedFlags(cmd)
}

func selectCategories(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Categories(ctx)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch categories from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Categories:",
		Options: values,
		Default: preselect(values),
	}
	err = survey.AskOne(prompt, &selected)

	filter.Categories = selected

	return err
}

func selectKinds(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Kinds(ctx, *filter)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch kinds from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Kinds:",
		Options: values,
		Default: preselect(values),
	}

	err = survey.AskOne(prompt, &selected)

	filter.Kinds = selected

	return err
}

func selectResources(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Resources(ctx, *filter)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch resources from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	options := make([]string, 0, len(values))
	for _, option := range values {
		options = append(options, fmt.Sprintf("%s/%s", strings.ToLower(option.Kind), strings.ToLower(option.Name)))
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Resources:",
		Options: options,
		Default: preselect(options),
	}

	err = survey.AskOne(prompt, &selected)

	query := make([]string, 0, len(selected))
	for _, s := range selected {
		parts := strings.Split(s, "/")
		query = append(query, strings.TrimSpace(parts[1]))
	}

	filter.Resources = query

	return err
}

func selectPolicies(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Policies(ctx, *filter)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch policies from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Policies:",
		Options: values,
		Default: preselect(values),
	}
	err = survey.AskOne(prompt, &selected)

	filter.Policies = selected
	return err
}

func selectNamespaces(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Namespaces(ctx, *filter)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch namespaces from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Namespaces:",
		Options: values,
		Default: preselect(values),
	}
	err = survey.AskOne(prompt, &selected)

	filter.Namespaces = selected

	return err
}

func selectSources(ctx context.Context, api policyreporter.API, filter *policyreporter.Filter) error {
	values, err := api.Sources(ctx)
	if err != nil {
		fmt.Println(chalk.Red, "[ERROR] Unable to fetch sources from API")
		return err
	}
	if len(values) == 0 {
		return ErrNoResult
	}

	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Sources:",
		Options: values,
		Default: preselect(values),
	}
	survey.AskOne(prompt, &selected)

	filter.Sources = selected

	return err
}

func selectResult(_ context.Context, _ policyreporter.API, filter *policyreporter.Filter) error {
	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Result:",
		Options: []string{policyreporter.Error, policyreporter.Fail, policyreporter.Warn, policyreporter.Pass, policyreporter.Skip},
	}
	err := survey.AskOne(prompt, &selected)

	filter.Status = selected

	return err
}

func selectSeverity(_ context.Context, _ policyreporter.API, filter *policyreporter.Filter) error {
	selected := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select Severity:",
		Options: []string{policyreporter.Low, policyreporter.Medium, policyreporter.High},
	}
	err := survey.AskOne(prompt, &selected)

	filter.Severities = selected

	return err
}
