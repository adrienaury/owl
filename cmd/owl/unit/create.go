package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// initCreateCommand initialize the cli unit list command
func initCreateCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create unit",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s unit create <<< '{"id": "my-unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			u := struct {
				ID string
			}{}

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

			err = unitDriver.Create(unit.NewUnit(u.ID))
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Created unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
