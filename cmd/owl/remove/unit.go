package remove

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

type singleUnit struct {
	ID          string
	Description string
}

// initUnitCommand initialize the cli remove unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Remove attributes from existing units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s remove unit <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			policyName := curRealm.Policy()

			policy, err := policyDriver.Get(policyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			units := []singleUnit{}

			if len(args) > 0 {
				description := ""
				if len(args) > 1 {
					description = args[1]
				}
				units = append(units, singleUnit{args[0], description})
			} else {
				b, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				err = json.Unmarshal(b, &units)
				if err != nil {
					u := singleUnit{}
					err = json.Unmarshal(b, &u)
					if err != nil || u.ID == "" {
						tmp := struct{ Units []singleUnit }{}
						err = json.Unmarshal(b, &tmp)
						if err != nil {
							cmd.PrintErrln(err)
							os.Exit(1)
						}
						units = tmp.Units
					} else {
						units = append(units, u)
					}
				}
			}

			flagRealm := cmd.Flag("realm")

			for _, u := range units {
				err := unitDriver.Remove(unit.NewUnit(u.ID, u.Description), policy.Objects()["unit"])
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.PrintErrf("Removed from unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
