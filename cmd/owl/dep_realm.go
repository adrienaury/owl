package main

import (
	"github.com/adrienaury/owl/cmd/owl/paths"
	domain "github.com/adrienaury/owl/pkg/domain/realm"
	infra "github.com/adrienaury/owl/pkg/infra/realm"
)

func newRealmStorage() domain.Storage {
	return infra.NewYAMLStorage(paths.LocalDir)
}

func newRealmDriver() domain.Driver {
	return domain.NewDriver(newRealmStorage())
}
