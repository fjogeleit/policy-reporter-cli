package model

import "github.com/kyverno/policy-reporter-cli/pkg/policyreporter"

type Group struct {
	Label string
	List  []policyreporter.PolicyReportResult
}
