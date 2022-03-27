package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyverno/policy-reporter-cli/pkg/forwarder"
	"github.com/kyverno/policy-reporter-cli/pkg/k8s"
	"github.com/kyverno/policy-reporter-cli/pkg/policyreporter"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientConfig *rest.Config
	kubeConfig   clientcmd.ClientConfig
)

type Resolver struct {
	config *Config
}

func (r *Resolver) KubeConfig() clientcmd.ClientConfig {
	if kubeConfig == nil {
		defferedConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		kubeConfig = defferedConfig
	}

	return kubeConfig
}

func (r *Resolver) ClientConfig() (*rest.Config, error) {
	if clientConfig == nil {
		config, err := r.KubeConfig().ClientConfig()
		if err != nil {
			return nil, err
		}
		clientConfig = config
	}

	return clientConfig, nil
}

func (r *Resolver) ForwardPolicyReporter(ctx context.Context) (*policyreporter.ForwardConnection, error) {
	kubeConfig, err := r.ClientConfig()
	if err != nil {
		return nil, err
	}

	prc := r.config.PolicyReporter

	options := []*forwarder.Option{
		{
			RemotePort: prc.Port,
			Source:     prc.Service,
			Namespace:  prc.Namespace,
		},
	}

	conn, err := policyreporter.Forward(ctx, options, kubeConfig)
	if err == forwarder.ErrServiceNotFound {
		fmt.Printf("Unable to connect to Policy Reporter with http://%s.%s:%d\n", strings.Split(prc.Service, "/")[1], prc.Namespace, prc.Port)
		fmt.Printf("Use the following env variables '%s', '%s', '%s' to customize your configuration\n", PolicyReporterNamespacEnv, PolicyReporterServiceEnv, PolicyReporterPortEnv)
	}

	return conn, err
}

func (r *Resolver) API(port uint16) policyreporter.API {
	return policyreporter.NewV1API(port)
}

func (r *Resolver) CurrentNamespace() (string, error) {
	namespace, _, err := r.KubeConfig().Namespace()

	return namespace, err
}

func (r *Resolver) K8sClient() (k8s.Client, error) {
	config, err := r.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8s.NewClient(client), nil
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{config: config}
}
