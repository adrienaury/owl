package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

var (
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver

	unitDriver  unit.Driver
	userDriver  user.Driver
	groupDriver group.Driver

	// local flags
	flagAllUnits bool
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(rd realm.Driver, cd credentials.Driver, und unit.Driver, usd user.Driver, gd group.Driver) {
	realmDriver = rd
	credentialsDriver = cd
	unitDriver = und
	userDriver = usd
	groupDriver = gd
}

type exportedUnit struct {
	ID          string
	Description string
	Users       user.List
	Groups      group.List
}

type exportedStruct struct {
	Units []exportedUnit
}

// InitCommand initialize the cli list command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List objects",
		Long:             "",
		Aliases:          []string{"ls"},
		Example:          fmt.Sprintf("  %[1]s list", parentCmd.Root().Name()),
		Args:             cobra.MaximumNArgs(1),
		PersistentPreRun: initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			var units unit.List
			var err error
			if flagAllUnits {
				units, err = unitDriver.List()
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
			} else {
				aunit, err := unitDriver.Get(cmd.Flag("unit").Value.String())
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				units = unit.NewList([]unit.Unit{aunit})
			}

			exportedUnits := []exportedUnit{}
			for _, u := range units.All() {
				unitDriver.Use(u.ID())

				users, err := userDriver.List()
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}

				groups, err := groupDriver.List()
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}

				exportedUnits = append(exportedUnits, exportedUnit{
					ID:          u.ID(),
					Description: u.Description(),
					Users:       users,
					Groups:      groups,
				})
			}

			b, err := json.Marshal(exportedStruct{exportedUnits})
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			cmd.Println(string(b))
		},
	}
	cmd.PersistentFlags().BoolVar(&flagAllUnits, "all-units", false, "export all units")
	parentCmd.AddCommand(cmd)
	initUserCommand(cmd)
	initUnitCommand(cmd)
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
