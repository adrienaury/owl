package user

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/password"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

var (
	userDriver        user.Driver
	passwordDriver    password.Driver
	unitDriver        unit.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(us user.Driver, pd password.Driver, un unit.Driver, r realm.Driver, c credentials.Driver) {
	userDriver = us
	passwordDriver = pd
	unitDriver = un
	realmDriver = r
	credentialsDriver = c
}

// InitCommand initialize the cli user command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user {list,create,delete,apply} [arguments ...]",
		Short:   "Manage users",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
	}
	parentCmd.AddCommand(cmd)
	initListCommand(cmd)
	initCreateCommand(cmd)
	initApplyCommand(cmd)
	initDeleteCommand(cmd)
	initPasswordCommand(cmd)
}

func initCredentialsAndUnit(cmd *cobra.Command, args []string) {
	initCredentials(cmd)
	initUnit(cmd)
}

func initCredentials(cmd *cobra.Command) {
	flagRealm := cmd.Flag("realm")
	if flagRealm == nil || strings.TrimSpace(flagRealm.Value.String()) == "" {
		cmd.PrintErrf("No active realm connection, use '--realm' flag or '%v realm login <realm ID>' command.", cmd.Root().Name())
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

	if err := credentialsDriver.Use(creds); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func initUnit(cmd *cobra.Command) {
	flagUnit := cmd.Flag("unit")
	if flagUnit == nil || strings.TrimSpace(flagUnit.Value.String()) == "" {
		cmd.PrintErrf("No unit selected, use '--unit' flag or '%v unit use <unit ID>' command.", cmd.Root().Name())
		cmd.PrintErrln()
		os.Exit(1)
	}

	unitDriver.Use(flagUnit.Value.String())
}
