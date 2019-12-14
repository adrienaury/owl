package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli update user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Replace specified attributes on users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s update user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			users := []exportedUser{}

			if len(args) > 0 {
				u := exportedUser{}
				u.ID = args[0]
				for _, arg := range args[1:] {
					argparts := strings.Split(arg, "=")
					if len(argparts) == 2 {
						switch argparts[0] {
						case "firstname":
							u.FirstNames = append(u.FirstNames, argparts[1])
						case "lastname":
							u.LastNames = append(u.LastNames, argparts[1])
						case "email":
							u.Emails = append(u.Emails, argparts[1])
						default:
							cmd.PrintErrln("Invalid attribute:", argparts[0])
							os.Exit(1)
						}
					}
				}
				users = append(users, u)
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
					if err != nil {
						cmd.PrintErrln(err)
						os.Exit(1)
					}
					users = append(users, u)
				}
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			for _, u := range users {
				err := userDriver.Update(user.NewUser(u.ID, u.FirstNames, u.LastNames, u.Emails))
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.PrintErrf("Updated user '%v' in unit '%v' of realm '%v'.", u.ID, flagUnit.Value, flagRealm.Value)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
