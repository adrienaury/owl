package apply

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

// initUnitCommand initialize the cli apply unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Create or replace units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s apply unit <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
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
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
