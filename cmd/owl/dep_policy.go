package main

import (
	"github.com/adrienaury/owl/cmd/owl/paths"
	"github.com/adrienaury/owl/pkg/domain/policy"
	domain "github.com/adrienaury/owl/pkg/domain/policy"
	infra "github.com/adrienaury/owl/pkg/infra/policy"
)

func newPolicyStorage() domain.Storage {
	return infra.NewYAMLStorage(paths.LocalDir)
}

func newPolicyDriver() policy.Driver {
	return domain.NewDriver(newPolicyStorage())
}
