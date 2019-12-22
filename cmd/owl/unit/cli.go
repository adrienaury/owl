package unit

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/policy"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

var (
	globalSession     *session.Session
	unitDriver        unit.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
	policyDriver      policy.Driver

	// init
	curRealm realm.Realm
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(u unit.Driver, r realm.Driver, c credentials.Driver, pold policy.Driver) {
	unitDriver = u
	realmDriver = r
	credentialsDriver = c
	policyDriver = pold
}

// SetSession inject the global session in the command.
func SetSession(s *session.Session) {
	globalSession = s
}

// InitCommand initialize the cli unit command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit [ID]",
		Short:   "Show or select current unit",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s unit dev`, parentCmd.Root().Name()),
		Args:    cobra.MaximumNArgs(1),
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				id := args[0]

				if id == "default" || id == "-" {
					unitDriver.Use("")
					globalSession.Unit = ""
					cmd.PrintErrf("Using default unit for next commands.")
				} else {
					policyName := curRealm.Policy()

					policy, err := policyDriver.Get(policyName)
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}

					unit, err := unitDriver.Get(id, policy.Objects()["unit"])
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}

					if unit == nil {
						cmd.PrintErrf("Unknown unit '%v'.", id)
						cmd.PrintErrln()
						os.Exit(1)
					}

					unitDriver.Use(id)
					globalSession.Unit = unit.ID()

					cmd.PrintErrf("Using unit '%v' for next commands.", unit.ID())
				}
			} else {
				currentUnit := ""
				flagUnit := cmd.Flag("unit")
				if flagUnit != nil {
					currentUnit = flagUnit.Value.String()
				}

				if currentUnit == "" {
					cmd.PrintErrf("Using default unit.")
				} else {
					cmd.PrintErrf("Using unit '%v'.", currentUnit)
				}
			}
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}

func initCredentials(cmd *cobra.Command, args []string) {
	flagRealm := cmd.Flag("realm")
	if flagRealm == nil || strings.TrimSpace(flagRealm.Value.String()) == "" {
		cmd.PrintErrf("No active realm connection, use '--realm' flag or '%v realm login' command.", cmd.Root().Name())
		cmd.PrintErrln()
		os.Exit(1)
	}

	var err error
	curRealm, err = realmDriver.Get(flagRealm.Value.String())
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	if curRealm == nil {
		cmd.PrintErrf("No realm with id '%v'.", flagRealm.Value.String())
		cmd.PrintErrln()
		os.Exit(1)
	}

	creds, err := credentialsDriver.Get(curRealm.URL(), curRealm.Username())
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
