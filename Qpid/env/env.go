package env

import (
	"cmp"
	"os"
)

type Env struct {
	Environment string
}

const (
	EnvironmentLocal      = "local"
	EnvironmentProduction = "production"
)

func (e *Env) IsLocal() bool {
	return e.Environment == EnvironmentLocal
}

func (e *Env) IsProduction() bool {
	return e.Environment == EnvironmentProduction
}

func GetEnv() Env {
	environment := cmp.Or(os.Getenv("ENVIRONMENT"), EnvironmentLocal)
	if environment != EnvironmentLocal && environment != EnvironmentProduction {
		environment = EnvironmentLocal
	}
	return Env{
		Environment: environment,
	}
}
