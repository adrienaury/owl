package imprt

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
	groupDriver       group.Driver
	userDriver        user.Driver
	unitDriver        unit.Driver
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers ...
func SetDrivers(gd group.Driver, us user.Driver, un unit.Driver, r realm.Driver, c credentials.Driver) {
	groupDriver = gd
	userDriver = us
	unitDriver = un
	realmDriver = r
	credentialsDriver = c
}

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

// InitCommand initialize the cli export command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "import [arguments ...]",
		Short:   "Import objects",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s import < objects.json`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
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
					cmd.PrintErrf("Updated unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
				} else {
					cmd.PrintErrf("Imported unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
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
						cmd.PrintErrf("Updated user '%v' in unit '%v' of realm '%v'.", us.ID, u.ID, flagRealm.Value)
					} else {
						cmd.PrintErrf("Imported user '%v' in unit '%v' of realm '%v'.", us.ID, u.ID, flagRealm.Value)
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
						cmd.PrintErrf("Updated group '%v' in unit '%v' of realm '%v'.", g.ID, u.ID, flagRealm.Value)
					} else {
						cmd.PrintErrf("Imported group '%v' in unit '%v' of realm '%v'.", g.ID, u.ID, flagRealm.Value)
					}
					cmd.PrintErrln()
				}
			}
		},
	}
	parentCmd.AddCommand(cmd)
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
