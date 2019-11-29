package main

import domain "github.com/adrienaury/owl/pkg/domain/user"

func newUserDriver(backend domain.Backend) domain.Driver {
	return domain.NewDriver(backend)
}
