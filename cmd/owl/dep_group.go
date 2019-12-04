package main

import domain "github.com/adrienaury/owl/pkg/domain/group"

func newGroupDriver(backend domain.Backend) domain.Driver {
	return domain.NewDriver(backend)
}
