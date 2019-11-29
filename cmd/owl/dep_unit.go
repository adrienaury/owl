package main

import domain "github.com/adrienaury/owl/pkg/domain/unit"

func newUnitDriver(backend domain.Backend) domain.Driver {
	return domain.NewDriver(backend)
}
