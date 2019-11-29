package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/credentials"
	infra "github.com/adrienaury/owl/pkg/infra/credentials"
)

func newCredentialsStorage() domain.Storage {
	return infra.NewHelperStorage()
}

func newCredentialsDriver(backend domain.Backend) domain.Driver {
	return domain.NewDriver(newCredentialsStorage(), backend)
}
