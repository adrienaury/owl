package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli apply group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Create or replace groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s apply group <<< '{"ID": "good-guys", "Members": ["batman"]}'`, parentCmd.Root().Name()),
		Args:    cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			groups := []exportedGroup{}

			if len(args) > 1 {
				g := exportedGroup{}
				g.ID = args[0]
				for _, arg := range args[1:] {
					argparts := strings.Split(arg, "=")
					if len(argparts) == 2 {
						if argparts[0] == "member" {
							g.Members = append(g.Members, argparts[1])
						} else {
							cmd.PrintErrln("Invalid attribute:", argparts[0])
							os.Exit(1)
						}
					}
				}
				groups = append(groups, g)
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
				modified, err := groupDriver.Apply(group.NewGroup(g.ID, g.Members...))
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				if modified {
					cmd.PrintErrf("Replaced group '%v' in unit '%v' of realm '%v'.", g.ID, flagUnit.Value, flagRealm.Value)
				} else {
					cmd.PrintErrf("Created group '%v' in unit '%v' of realm '%v'.", g.ID, flagUnit.Value, flagRealm.Value)
				}
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
