package realm

import (
	"fmt"

	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

var (
	globalSession     *session.Session
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers ...
func SetDrivers(r realm.Driver, c credentials.Driver) {
	realmDriver = r
	credentialsDriver = c
}

// SetSession ...
func SetSession(s *session.Session) {
	globalSession = s
}

// InitCommand initialize the cli realm command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "realm {set,create,delete,list,login} [arguments ...]",
		Short:   "Manage realms",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm set dev ldap://dev.my-company.com/dc=example,dc=com", parentCmd.Root().Name()),
		Aliases: []string{"rlm"},
	}
	parentCmd.AddCommand(cmd)
	initSetCommand(cmd)
	initCreateCommand(cmd)
	initListCommand(cmd)
	initLoginCommand(cmd)
	initDeleteCommand(cmd)
}
