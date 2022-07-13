package k8s

import (
	"context"
	"fmt"
	"log"
	"strings"

	pr "github.com/kyverno/policy-reporter-cli/pkg/policyreporter"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Client interface {
	LabelFilter(ctx context.Context, results pr.ResultList, labels string) pr.ResultList
}

type k8sClient struct {
	client dynamic.Interface
}

func (k *k8sClient) LabelFilter(ctx context.Context, results pr.ResultList, labels string) pr.ResultList {
	groups := make(map[string][]pr.PolicyReportResult, 0)
	filtered := make([]pr.PolicyReportResult, 0, results.Count)

	for _, res := range results.Items {
		key := resKey(res)
		if list, ok := groups[key]; ok {
			groups[key] = append(list, res)
		} else {
			groups[key] = []pr.PolicyReportResult{res}
		}
	}

	for _, group := range groups {
		resources, err := k.fetchResources(ctx, group[0].Kind, group[0].APIVersion, labels)
		if err != nil {
			log.Printf("%s/%s: %s", group[0].APIVersion, group[0].Kind, err)
			continue
		}

		list := k.filter(group, resources)

		filtered = append(filtered, list...)
	}

	return pr.ResultList{Items: filtered, Count: len(filtered)}
}

func (k *k8sClient) fetchResources(ctx context.Context, kind, apiVersion, labels string) (*unstructured.UnstructuredList, error) {
	var group, version string

	parts := strings.Split(apiVersion, "/")
	if len(parts) == 2 {
		group = strings.TrimSpace(parts[0])
		version = strings.TrimSpace(parts[1])
	} else if len(parts) == 1 {
		version = strings.TrimSpace(parts[0])
	} else {
		return nil, fmt.Errorf("Invalid apiVersion: %s", apiVersion)
	}

	// @TODO Check for more stable solutions
	resource, _ := meta.UnsafeGuessKindToResource(schema.GroupVersionKind{Group: group, Version: version, Kind: strings.ToLower(kind)})

	return k.client.Resource(resource).List(ctx, v1.ListOptions{LabelSelector: labels})
}

func (k *k8sClient) filter(results []pr.PolicyReportResult, list *unstructured.UnstructuredList) []pr.PolicyReportResult {
	filtered := make([]pr.PolicyReportResult, 0, len(results))

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

	return filtered
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
