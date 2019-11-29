package main

import "github.com/adrienaury/owl/pkg/infra"

func backend() infra.BackendLDAP {
	return infra.NewBackendLDAP()
}
