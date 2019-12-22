package delete

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type singleUnit struct {
	ID          string
	Description string
}

// initUnitCommand initialize the cli delete unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Delete units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s delete unit <<< '{"ID": "my-unit"}'`, parentCmd.Root().Name()),
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
				for _, arg := range args {
					u := singleUnit{}
					u.ID = arg
					units = append(units, u)
				}
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
				err := unitDriver.Delete(u.ID, policy.Objects()["unit"])
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.PrintErrf("Deleted unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
