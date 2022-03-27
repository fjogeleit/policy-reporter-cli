package policyreporter

// Result Enum defined for PolicyReport
type Result = string

// Severity Enum defined for PolicyReport
type Severity = string

// Enums for predefined values from the PolicyReport spec
const (
	Fail  Result = "fail"
	Warn  Result = "warn"
	Error Result = "error"
	Pass  Result = "pass"
	Skip  Result = "skip"

	Low    Severity = "low"
	Medium Severity = "medium"
	High   Severity = "high"
)

var (
	AllResults = []Result{Error, Fail, Warn, Pass, Skip}
)

type Target struct {
	Name                  string   `json:"name"`
	MinimumPriority       string   `json:"minimumPriority"`
	Sources               []string `json:"sources,omitempty"`
	SkipExistingOnStartup bool     `json:"skipExistingOnStartup"`
}

type Resource struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type PolicyReportResult struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace,omitempty"`
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Name       string            `json:"name"`
	Message    string            `json:"message"`
	Category   string            `json:"category"`
	Policy     string            `json:"policy"`
	Rule       string            `json:"rule"`
	Status     string            `json:"status"`
	Severity   string            `json:"severity,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}
