package unit

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

var (
	unitDriver        unit.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers ...
func SetDrivers(u unit.Driver, r realm.Driver, c credentials.Driver) {
	unitDriver = u
	realmDriver = r
	credentialsDriver = c
}

// NewCommand implements the cli unit command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unit {list,create,apply,use} [arguments ...]",
		Short:   "Manage realms",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm set dev ldap://dev.my-company.com/dc=example,dc=com", fullName),
		Aliases: []string{"rlm"},
	}
	cmd.AddCommand(newListCommand(fullName, err, out, in))
	cmd.AddCommand(newCreateCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
