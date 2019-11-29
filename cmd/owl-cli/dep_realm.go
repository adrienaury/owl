package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/realm"
	infra "github.com/adrienaury/owl/pkg/infra/realm"
)

func realmStorage() domain.Storage {
	return infra.NewYAMLStorage()
}

func realmDriver() domain.Driver {
	return domain.NewDriver(realmStorage())
}
