package group

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/spf13/cobra"
)

// initCreateCommand initialize the cli group create command
func initCreateCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create group",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s group create <<< '{"ID": "my-group"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			g := struct {
				ID      string
				Members []string
			}{}

			b, err := ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = json.Unmarshal(b, &g)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = groupDriver.Create(group.NewGroup(g.ID, g.Members...))
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Created group '%v' in unit '%v' of realm '%v'.", g.ID, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
