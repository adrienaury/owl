package delete

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli delete user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Delete users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s delete user <<< '{"ID": "batman"}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			users := []exportedUser{}

			if len(args) > 0 {
				for _, arg := range args {
					u := exportedUser{}
					u.ID = arg
					users = append(users, u)
				}
			} else {
				b, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				err = json.Unmarshal(b, &users)
				if err != nil {
					u := exportedUser{}
					err = json.Unmarshal(b, &u)
					if err != nil || u.ID == "" {
						tmp := struct{ Users []exportedUser }{}
						err = json.Unmarshal(b, &tmp)
						if err != nil {
							cmd.PrintErrln(err)
							os.Exit(1)
						}
						users = tmp.Users
					} else {
						users = append(users, u)
					}
				}
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			for _, u := range users {
				err := userDriver.Delete(u.ID)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.PrintErrf("Deleted user '%v' in unit '%v' of realm '%v'.", u.ID, flagUnit.Value, flagRealm.Value)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
