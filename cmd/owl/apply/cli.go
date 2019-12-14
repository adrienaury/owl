package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
)

type exportedUser struct {
	ID         string
	FirstNames []string
	LastNames  []string
	Emails     []string
}

type exportedGroup struct {
	ID      string
	Members []string
}

type exportedUnit struct {
	ID          string
	Description string
	Users       []exportedUser
	Groups      []exportedGroup
}

type exportedStruct struct {
	Units []exportedUnit
}

// SetDrivers inject required domain drivers in the command.
func SetDrivers(rd realm.Driver, cd credentials.Driver, und unit.Driver, usd user.Driver, gd group.Driver) {
	realmDriver = rd
	credentialsDriver = cd
	unitDriver = und
	userDriver = usd
	groupDriver = gd
}

// InitCommand initialize the cli create command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:              "apply",
		Short:            "Create or replace objects",
		Long:             "",
		Aliases:          []string{"replace", "ap"},
		Example:          fmt.Sprintf("  %[1]s list", parentCmd.Root().Name()),
		Args:             cobra.MaximumNArgs(1),
		PersistentPreRun: initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			b, err := ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			s := exportedStruct{[]exportedUnit{}}
			err = json.Unmarshal(b, &s)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagRealm := cmd.Flag("realm")
			for _, u := range s.Units {
				modified, err := unitDriver.Apply(unit.NewUnit(u.ID, u.Description))
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				if modified {
					cmd.PrintErrf("Replaced unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
				} else {
					cmd.PrintErrf("Created unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
				}
				cmd.PrintErrln()
				unitDriver.Use(u.ID)
				for _, us := range u.Users {
					modified, err := userDriver.Apply(user.NewUser(us.ID, us.FirstNames, us.LastNames, us.Emails))
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}
					if modified {
						cmd.PrintErrf("Replaced user '%v' in unit '%v' of realm '%v'.", us.ID, u.ID, flagRealm.Value)
					} else {
						cmd.PrintErrf("Created user '%v' in unit '%v' of realm '%v'.", us.ID, u.ID, flagRealm.Value)
					}
					cmd.PrintErrln()
				}
				for _, g := range u.Groups {
					modified, err := groupDriver.Apply(group.NewGroup(g.ID, g.Members...))
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}
					if modified {
						cmd.PrintErrf("Replaced group '%v' in unit '%v' of realm '%v'.", g.ID, u.ID, flagRealm.Value)
					} else {
						cmd.PrintErrf("Created group '%v' in unit '%v' of realm '%v'.", g.ID, u.ID, flagRealm.Value)
					}
					cmd.PrintErrln()
				}
			}
		},
	}
	parentCmd.AddCommand(cmd)
	initUnitCommand(cmd)
	initUserCommand(cmd)
	initGroupCommand(cmd)
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
