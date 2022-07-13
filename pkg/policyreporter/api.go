package policyreporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type API interface {
	Categories(context.Context) ([]string, error)
	Kinds(context.Context, Filter) ([]string, error)
	ClusterKinds(context.Context, Filter) ([]string, error)
	Resources(context.Context, Filter) ([]Resource, error)
	ClusterResources(context.Context, Filter) ([]Resource, error)
	Namespaces(context.Context, Filter) ([]string, error)
	Policies(context.Context, Filter) ([]string, error)
	ClusterPolicies(context.Context, Filter) ([]string, error)
	Sources(context.Context) ([]string, error)
	ClusterSources(context.Context) ([]string, error)
	Targets(context.Context) ([]Target, error)
	Results(context.Context, Filter) (ResultList, error)
	ClusterResults(context.Context, Filter) (ResultList, error)
}

type Filter struct {
	Kinds      []string
	Categories []string
	Namespaces []string
	Sources    []string
	Policies   []string
	Severities []string
	Status     []string
	Resources  []string
}

type api struct {
	URL    string
	client *http.Client
}

func (a *api) Categories(ctx context.Context) ([]string, error) {
	var categories = make([]string, 0)

	req, err := http.NewRequestWithContext(ctx, "GET", a.fullPath("categories"), new(bytes.Buffer))
	if err != nil {
		return categories, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return categories, err
	}

	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (a *api) Kinds(ctx context.Context, filter Filter) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "namespaced-resources/kinds", &values, filter)

	return values, err
}

func (a *api) ClusterKinds(ctx context.Context, filter Filter) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "cluster-resources/kinds", &values, filter)

	return values, err
}

func (a *api) Resources(ctx context.Context, filter Filter) ([]Resource, error) {
	var values = make([]Resource, 0)

	err := a.Request(ctx, "namespaced-resources/resources", &values, filter)

	return values, err
}

func (a *api) ClusterResources(ctx context.Context, filter Filter) ([]Resource, error) {
	var values = make([]Resource, 0)

	err := a.Request(ctx, "cluster-resources/resources", &values, filter)

	return values, err
}

func (a *api) Policies(ctx context.Context, filter Filter) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "namespaced-resources/policies", &values, filter)

	return values, err
}

func (a *api) ClusterPolicies(ctx context.Context, filter Filter) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "cluster-resources/policies", &values, filter)

	return values, err
}

func (a *api) Sources(ctx context.Context) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "namespaced-resources/sources", &values, Filter{})

	return values, err
}

func (a *api) ClusterSources(ctx context.Context) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "cluster-resources/sources", &values, Filter{})

	return values, err
}

func (a *api) Namespaces(ctx context.Context, filter Filter) ([]string, error) {
	var values = make([]string, 0)

	err := a.Request(ctx, "namespaces", &values, filter)

	return values, err
}

func (a *api) Targets(ctx context.Context) ([]Target, error) {
	var targets = make([]Target, 0)

	err := a.Request(ctx, "targets", &targets, Filter{})

	return targets, err
}

func (a *api) Results(ctx context.Context, filter Filter) (ResultList, error) {
	var results = ResultList{}

	err := a.Request(ctx, "namespaced-resources/results", &results, filter)

	return results, err
}

func (a *api) ClusterResults(ctx context.Context, filter Filter) (ResultList, error) {
	var results = ResultList{}

	err := a.Request(ctx, "cluster-resources/results", &results, filter)

	return results, err
}

func (a *api) fullPath(path string) string {
	return fmt.Sprintf("%s/v1/%s", a.URL, path)
}

func (a *api) Request(ctx context.Context, url string, body interface{}, filter Filter) error {
	req, err := http.NewRequestWithContext(ctx, "GET", a.fullPath(url), new(bytes.Buffer))
	if err != nil {
		return err
	}

	req.URL.RawQuery = buildQuery(filter).Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return err
	}

	return nil
}

func NewV1API(port uint16) API {
	return &api{
		URL:    fmt.Sprintf("http://localhost:%d", port),
		client: http.DefaultClient,
	}
}

func buildQuery(filter Filter) url.Values {
	query := url.Values{}

	for _, value := range filter.Kinds {
		query.Add("kinds", value)
	}
	for _, value := range filter.Resources {
		query.Add("resources", value)
	}
	for _, value := range filter.Sources {
		query.Add("sources", value)
	}
	for _, value := range filter.Categories {
		query.Add("categories", value)
	}
	for _, value := range filter.Severities {
		query.Add("severities", value)
	}
	for _, value := range filter.Policies {
		query.Add("policies", value)
	}
	for _, value := range filter.Status {
		query.Add("status", value)
	}
	for _, value := range filter.Namespaces {
		query.Add("namespaces", value)
	}

	return query
}
