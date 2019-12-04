package main

import (
	"os"

	domain "github.com/adrienaury/owl/pkg/domain/password"
	infra "github.com/adrienaury/owl/pkg/infra"
)

func newMailService() domain.SecretPusher {
	// TODO : viper
	addr := os.Getenv("SMTP_ADDR")
	from := os.Getenv("SMTP_FROM")
	return infra.NewMailService(addr, from)
}

func newPasswordDriver(backend domain.Backend) domain.Driver {
	return domain.NewDriver(backend, newMailService())
}
