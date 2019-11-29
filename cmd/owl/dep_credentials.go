package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/credentials"
	infra "github.com/adrienaury/owl/pkg/infra/credentials"
)

func credentialsStorage() domain.Storage {
	return infra.NewHelperStorage()
}

func credentialsDriver() domain.Driver {
	return domain.NewDriver(credentialsStorage(), nil)
}
