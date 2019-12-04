package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/realm"
	infra "github.com/adrienaury/owl/pkg/infra/realm"
)

func newRealmStorage() domain.Storage {
	return infra.NewYAMLStorage()
}

func newRealmDriver() domain.Driver {
	return domain.NewDriver(newRealmStorage())
}
