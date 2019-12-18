package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/cmd/owl/paths"
	infra "github.com/adrienaury/owl/pkg/infra/realm"
)

func newRealmStorage() domain.Storage {
	return infra.NewYAMLStorage(paths.LocalDir)
}

func newRealmDriver() domain.Driver {
	return domain.NewDriver(newRealmStorage())
}
