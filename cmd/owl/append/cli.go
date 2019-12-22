package append

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/policy"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

var (
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
	policyDriver      policy.Driver

	unitDriver  unit.Driver
	userDriver  user.Driver
	groupDriver group.Driver

	// init
	curRealm realm.Realm
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(rd realm.Driver, cd credentials.Driver, pold policy.Driver, und unit.Driver, usd user.Driver, gd group.Driver) {
	realmDriver = rd
	credentialsDriver = cd
	policyDriver = pold
	unitDriver = und
	userDriver = usd
	groupDriver = gd
}

type exportedUser struct {
	ID         string
	FirstNames []string
	LastNames  []string
	Emails     []string
}

type exportedUsers struct {
	Users []exportedUser
}

type exportedGroup struct {
	ID      string
	Members []string
}

type exportedGroups struct {
	Groups []exportedGroup
}

type exportedUnit struct {
	ID          string
	Description string
	Users       []exportedUser
	Groups      []exportedGroup
}

type exportedUnits struct {
	Units []exportedUnit
}

func appendGroups(cmd *cobra.Command, groups []exportedGroup, realmID, unitID string) {
	unitDriver.Use(unitID)
	for _, g := range groups {
		err := groupDriver.Append(group.NewGroup(g.ID, g.Members...))
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.PrintErrf("Appended to group '%v' in unit '%v' of realm '%v'.", g.ID, unitID, realmID)
		cmd.PrintErrln()
	}
}

func appendUsers(cmd *cobra.Command, users []exportedUser, realmID, unitID string) {
	unitDriver.Use(unitID)
	for _, u := range users {
		err := userDriver.Append(user.NewUser(u.ID, u.FirstNames, u.LastNames, u.Emails))
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.PrintErrf("Appended to user '%v' in unit '%v' of realm '%v'.", u.ID, unitID, realmID)
		cmd.PrintErrln()
	}
}

func appendUnits(cmd *cobra.Command, units []exportedUnit, realmID string, p policy.Policy) {
	for _, u := range units {
		err := unitDriver.Append(unit.NewUnit(u.ID, u.Description), p.Objects()["unit"])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.PrintErrf("Appended to unit '%v' in realm '%v'.", u.ID, realmID)
		cmd.PrintErrln()
		appendUsers(cmd, u.Users, realmID, u.ID)
		appendGroups(cmd, u.Groups, realmID, u.ID)
	}
}

// InitCommand initialize the cli append command
func InitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:              "append",
		Short:            "Append attributes to existing objects",
		Long:             "",
		Aliases:          []string{"add"},
		Example:          fmt.Sprintf("  %[1]s append", parentCmd.Root().Name()),
		Args:             cobra.NoArgs,
		PersistentPreRun: initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			policyName := curRealm.Policy()

			policy, err := policyDriver.Get(policyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			b, err := ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			realmID := cmd.Flag("realm").Value.String()
			unitID := cmd.Flag("unit").Value.String()

			format1 := exportedUnits{}
			format2 := exportedUsers{}
			format3 := exportedGroups{}

			err = json.Unmarshal(b, &format1)
			if err == nil && len(format1.Units) > 0 {
				appendUnits(cmd, format1.Units, realmID, policy)
				return
			}

			err = json.Unmarshal(b, &format2)
			if err == nil && len(format2.Users) > 0 {
				appendUsers(cmd, format2.Users, realmID, unitID)
				return
			}

			err = json.Unmarshal(b, &format3)
			if err == nil && len(format3.Groups) > 0 {
				appendGroups(cmd, format3.Groups, realmID, unitID)
				return
			}

			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.PrintErrln("No valid data.")
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

func initUnit(cmd *cobra.Command) {
	flagUnit := cmd.Flag("unit")
	if flagUnit == nil || strings.TrimSpace(flagUnit.Value.String()) == "" {
		unitDriver.Use("")
	} else {
		unitDriver.Use(flagUnit.Value.String())
	}
}
