package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

// initCreateCommand initialize the cli user create command
func initCreateCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create user",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'
  %[1]s user create batman firstname=Bruce lastname=Wayne`, parentCmd.Root().Name()),
		Args:   cobra.ArbitraryArgs,
		PreRun: initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			u := struct {
				ID         string
				FirstNames []string
				LastNames  []string
				Emails     []string
			}{}

			if len(args) > 0 {
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
			} else {
				b, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}

				err = json.Unmarshal(b, &u)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
			}

			err := userDriver.Create(user.NewUser(u.ID, u.FirstNames, u.LastNames, u.Emails))
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Created user '%v' in unit '%v' of realm '%v'.", u.ID, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
