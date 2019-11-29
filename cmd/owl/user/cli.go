package user

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

var (
	userDriver        user.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers ...
func SetDrivers(u user.Driver, r realm.Driver, c credentials.Driver) {
	userDriver = u
	realmDriver = r
	credentialsDriver = c
}

// NewCommand implements the cli unit command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user {list,create,delete,apply} [arguments ...]",
		Short:   "Manage users",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, fullName),
		Aliases: []string{"us"},
	}
	cmd.AddCommand(newListCommand(fullName, err, out, in))
	cmd.AddCommand(newCreateCommand(fullName, err, out, in))
	cmd.AddCommand(newDeleteCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
