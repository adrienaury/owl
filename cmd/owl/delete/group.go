package delete

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli delete group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Delete groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s delete group <<< '{"ID": "good-guys"}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			groups := []exportedGroup{}

			if len(args) > 0 {
				for _, arg := range args {
					g := exportedGroup{}
					g.ID = arg
					groups = append(groups, g)
				}
			} else {
				b, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				err = json.Unmarshal(b, &groups)
				if err != nil {
					g := exportedGroup{}
					err = json.Unmarshal(b, &g)
					if err != nil || g.ID == "" {
						tmp := struct{ Groups []exportedGroup }{}
						err = json.Unmarshal(b, &tmp)
						if err != nil {
							cmd.PrintErrln(err)
							os.Exit(1)
						}
						groups = tmp.Groups
					} else {
						groups = append(groups, g)
					}
				}
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			for _, g := range groups {
				err := groupDriver.Delete(g.ID)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.PrintErrf("Deleted group '%v' in unit '%v' of realm '%v'.", g.ID, flagUnit.Value, flagRealm.Value)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
