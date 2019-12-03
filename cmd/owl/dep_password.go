package main

import (
	domain "github.com/adrienaury/owl/pkg/domain/password"
	infra "github.com/adrienaury/owl/pkg/infra"
)

func newMailService() domain.MailService {
	return infra.NewMailService()
}

func newPasswordDriver() domain.Driver {
	return domain.NewDriver(newMailService())
}
