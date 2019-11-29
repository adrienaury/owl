package realm

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

var (
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers ...
func SetDrivers(r realm.Driver, c credentials.Driver) {
	realmDriver = r
	credentialsDriver = c
}

// NewCommand implements the cli realm command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "realm {set,list,login} [arguments ...]",
		Short:   "Manage realms",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm set dev ldap://dev.my-company.com/dc=example,dc=com", fullName),
		Aliases: []string{"rlm"},
	}
	cmd.AddCommand(newSetCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
