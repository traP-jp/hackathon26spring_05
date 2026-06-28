package env

import (
	"cmp"
	"os"
)

type Env struct {
	Environment     string
	TraqHost        string
	TraqAccessToken string
}

const (
	envKeyEnvironment     = "ENVIRONMENT"
	envKeyTraqHost        = "TRAQ_HOST"
	envKeyTraqAccessToken = "TRAQ_ACCESS_TOKEN"
)

const (
	environmentLocal      = "local"
	environmentProduction = "production"
	defaultTraqHost       = "https://q.trap.jp"
)

func (e *Env) IsLocal() bool {
	return e.Environment == environmentLocal
}

func (e *Env) IsProduction() bool {
	return e.Environment == environmentProduction
}

func GetEnv() Env {
	environment := cmp.Or(os.Getenv(envKeyEnvironment), environmentLocal)
	if environment != environmentLocal && environment != environmentProduction {
		environment = environmentLocal
	}
	return Env{
		Environment:     environment,
		TraqHost:        cmp.Or(os.Getenv(envKeyTraqHost), defaultTraqHost),
		TraqAccessToken: os.Getenv(envKeyTraqAccessToken),
	}
}
