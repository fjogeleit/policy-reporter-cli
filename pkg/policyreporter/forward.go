package policyreporter

import (
	"context"

	"github.com/kyverno/policy-reporter-cli/pkg/forwarder"
	"k8s.io/client-go/rest"
)

type ForwardConnection struct {
	Port  uint16
	Close func()
}

func Forward(ctx context.Context, options []*forwarder.Option, kubeConfig *rest.Config) (*ForwardConnection, error) {
	ret, err := forwarder.Exec(ctx, options, kubeConfig)
	if err != nil {
		return nil, err
	}

	ports, err := ret.Ready()
	if err != nil {
		return nil, err
	}

	if len(ports) < 1 {
		return nil, err
	}

	return &ForwardConnection{ports[0][0].Local, ret.Close}, nil
}
