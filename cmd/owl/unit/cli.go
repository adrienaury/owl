package unit

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

var (
	globalSession     *session.Session
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

// SetSession ...
func SetSession(s *session.Session) {
	globalSession = s
}

// NewCommand implements the cli unit command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unit {list,create,apply,delete,use} [arguments ...]",
		Short:   "Manage realms",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s unit create <<< '{"ID": "my-unit"}'`, fullName),
	}
	cmd.AddCommand(newListCommand(fullName, err, out, in))
	cmd.AddCommand(newCreateCommand(fullName, err, out, in))
	cmd.AddCommand(newDeleteCommand(fullName, err, out, in))
	cmd.AddCommand(newUseCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}

func initCredentials(cmd *cobra.Command, args []string) {
	flagRealm := cmd.Flag("realm")
	if flagRealm == nil || strings.TrimSpace(flagRealm.Value.String()) == "" {
		cmd.PrintErrf("No active realm connection, use '--realm' flag or '%v realm login' command.", cmd.Root().Name())
		cmd.PrintErrln()
		os.Exit(1)
	}

	realm, err := realmDriver.Get(flagRealm.Value.String())
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	if realm == nil {
		cmd.PrintErrf("No realm with id '%v'.", flagRealm.Value.String())
		cmd.PrintErrln()
		os.Exit(1)
	}

	creds, err := credentialsDriver.Get(realm.URL(), realm.Username())
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	if creds == nil {
		cmd.PrintErrf("No credentials for realm '%[1]v', use '%[2]v realm login %[1]v' command.", flagRealm.Value.String(), cmd.Root().Name())
		cmd.PrintErrln()
		os.Exit(1)
	}

	credentialsDriver.Use(creds)
}
