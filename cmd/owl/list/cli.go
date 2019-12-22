package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/policy"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
	policyDriver      policy.Driver

	unitDriver  unit.Driver
	userDriver  user.Driver
	groupDriver group.Driver

	// local flags
	flagAllUnits bool

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
		Aliases:          []string{"ls", "export"},
		Example:          fmt.Sprintf("  %[1]s list", parentCmd.Root().Name()),
		Args:             cobra.MaximumNArgs(1),
		PersistentPreRun: initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			policyName := curRealm.Policy()

			policy, err := policyDriver.Get(policyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			var units unit.List
			if flagAllUnits {
				units, err = unitDriver.List(policy.Objects()["unit"])
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
			} else {
				if cmd.Flag("unit") != nil && strings.TrimSpace(cmd.Flag("unit").Value.String()) != "" {
					aunit, err := unitDriver.Get(cmd.Flag("unit").Value.String(), policy.Objects()["unit"])
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}
					if aunit == nil {
						cmd.PrintErrf("No unit named '%v'.", cmd.Flag("unit").Value.String())
						cmd.PrintErrln()
						os.Exit(1)
					}
					units = unit.NewList([]unit.Unit{aunit})
				} else {
					units = unit.NewList([]unit.Unit{unit.NewUnit("", "")})
				}
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

			output := "json"
			flagOutput := cmd.Flag("output")
			if flagOutput != nil && strings.TrimSpace(flagOutput.Value.String()) != "" {
				output = flagOutput.Value.String()
			}

			switch output {
			case "json":
				b, err := json.Marshal(exportedStruct{exportedUnits})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "yaml":
				b, err := yaml.Marshal(exportedStruct{exportedUnits})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "table":
				w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "Unit", "User", "First Names", "Last Names", "E-mails", "Member Of")
				for _, un := range exportedUnits {
					for _, us := range un.Users.All() {
						memberOf := []string{}
						for _, g := range un.Groups.All() {
							for _, member := range g.Members() {
								if member == us.ID() {
									memberOf = append(memberOf, g.ID())
									continue
								}
							}
						}
						unID := un.ID
						if unID == "" {
							unID = "default"
						}
						fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", unID, us.ID(), strings.Join(us.FirstNames(), ", "), strings.Join(us.LastNames(), ", "), strings.Join(us.Emails(), ", "), strings.Join(memberOf, ", "))
					}
				}
				w.Flush()
			default:
				cmd.PrintErrf("Invalid output format : %v", output)
				cmd.PrintErrln()
			}
		},
	}
	cmd.PersistentFlags().BoolVar(&flagAllUnits, "all-units", false, "export all units")
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
