package main

import "github.com/adrienaury/owl/pkg/infra"

func newBackend() infra.BackendLDAP {
	return infra.NewBackendLDAP()
}
