package group

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

var (
	groupDriver       group.Driver
	unitDriver        unit.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(gd group.Driver, un unit.Driver, r realm.Driver, c credentials.Driver) {
	groupDriver = gd
	unitDriver = un
	realmDriver = r
	credentialsDriver = c
}

// InitCommand initialize the cli group command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group {list,create,delete,member} [arguments ...]",
		Short:   "Manage groups",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s group create <<< '{"ID": "bad-guys", "Members": ["batman", "joker"]}'`, parentCmd.Root().Name()),
		Aliases: []string{"grp"},
	}
	parentCmd.AddCommand(cmd)
	initListCommand(cmd)
	initCreateCommand(cmd)
	initDeleteCommand(cmd)
	initMemberCommand(cmd)
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

	credentialsDriver.Use(creds)
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
