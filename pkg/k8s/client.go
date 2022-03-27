package k8s

import (
	"context"
	"fmt"
	"strings"

	pr "github.com/kyverno/policy-reporter-cli/pkg/policyreporter"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Client interface {
	LabelFilter(ctx context.Context, results []pr.PolicyReportResult, labels string) []pr.PolicyReportResult
}

type k8sClient struct {
	client dynamic.Interface
}

func (k *k8sClient) LabelFilter(ctx context.Context, results []pr.PolicyReportResult, labels string) []pr.PolicyReportResult {
	groups := make(map[string][]pr.PolicyReportResult, 0)
	filtered := make([]pr.PolicyReportResult, 0, len(results))

	for _, res := range results {
		key := resKey(res)
		if list, ok := groups[key]; ok {
			groups[key] = append(list, res)
		} else {
			groups[key] = []pr.PolicyReportResult{res}
		}
	}

	for _, group := range groups {
		list, err := k.filter(ctx, group, group[0].Kind, group[0].APIVersion, labels)
		if err != nil {
			continue
		}
		filtered = append(filtered, list...)
	}

	return filtered
}

func (k *k8sClient) filter(ctx context.Context, results []pr.PolicyReportResult, kind, apiVersion, labels string) ([]pr.PolicyReportResult, error) {
	var group, version string

	filtered := make([]pr.PolicyReportResult, 0, len(results))

	parts := strings.Split(apiVersion, "/")
	if len(parts) == 2 {
		group = strings.TrimSpace(parts[0])
		version = strings.TrimSpace(parts[1])
	} else if len(parts) == 2 {
		version = strings.TrimSpace(parts[0])
	} else {
		return results, fmt.Errorf("Invalid apiVersion: %s", apiVersion)
	}

	// @TODO Check for more stable solutions
	resource, _ := meta.UnsafeGuessKindToResource(schema.GroupVersionKind{Group: group, Version: version, Kind: strings.ToLower(kind)})

	list, err := k.client.
		Resource(resource).
		List(ctx, v1.ListOptions{LabelSelector: labels})
	if err != nil {
		fmt.Println(err)
		return results, err
	}

	resources := make([]string, 0, len(list.Items))
	for _, resource := range list.Items {
		metadata, ok := resource.Object["metadata"].(map[string]interface{})
		if !ok {
			continue
		}
		name := toString(metadata["name"])

		if name != "" {
			resources = append(resources, name)
		}
	}

	for _, res := range results {
		if contains(res.Name, resources) {
			filtered = append(filtered, res)
		}
	}

	return filtered, nil
}

func toString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}

	return ""
}

func resKey(value pr.PolicyReportResult) string {
	return fmt.Sprintf("%s/%s", value.APIVersion, value.Kind)
}

func contains(name string, names []string) bool {
	for _, s := range names {
		if strings.EqualFold(s, name) {
			return true
		}
	}

	return false
}

func NewClient(client dynamic.Interface) Client {
	return &k8sClient{client}
}
