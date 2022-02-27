package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

const (
	PolicyReporterServiceEnv  = "POLICY_REPORTER_SERVICE"
	PolicyReporterNamespacEnv = "POLICY_REPORTER_NAMESPACE"
	PolicyReporterPortEnv     = "POLICY_REPORTER_PORT"
)

type PolicyReporter struct {
	Namespace string `mapstructure:"namespace"`
	Service   string `mapstructure:"service"`
	Port      int    `mapstructure:"port"`
}

// Config of the PolicyReporter
type Config struct {
	PolicyReporter PolicyReporter `mapstructure:"policyreporter"`
}

func LoadConfig() *Config {
	v := viper.New()
	v.SetDefault("policyreporter.service", "svc/policy-reporter")
	v.SetDefault("policyreporter.namespace", "policy-reporter")
	v.SetDefault("policyreporter.port", 8080)

	v.AutomaticEnv()
	v.ReadInConfig()

	c := &Config{}

	v.Unmarshal(c)

	if value, present := os.LookupEnv(PolicyReporterNamespacEnv); present {
		c.PolicyReporter.Namespace = value
	}
	if value, present := os.LookupEnv(PolicyReporterServiceEnv); present {
		c.PolicyReporter.Service = fmt.Sprintf("svc/%s", value)
	}
	if value, present := os.LookupEnv(PolicyReporterPortEnv); present {
		port, err := strconv.Atoi(value)
		if err == nil {
			c.PolicyReporter.Port = port
		} else {
			fmt.Printf("[WARNING] Unable to parse port '%s' using default 8080\n", value)
		}
	}

	return c
}
