package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// initCreateCommand initialize the cli unit create command
func initCreateCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create unit",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s unit create <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args: func(cmd *cobra.Command, args []string) error {
			err1 := cobra.ExactArgs(2)(cmd, args)
			err2 := cobra.NoArgs(cmd, args)
			if err1 != nil && err2 != nil {
				return fmt.Errorf("accepts 0 or 2 arg(s), received %d", len(args))
			}
			return nil
		},
		PreRun: initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			u := struct {
				ID          string
				Description string
			}{}

			if len(args) > 0 {
				u.ID = args[0]
				u.Description = args[1]
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

			err := unitDriver.Create(unit.NewUnit(u.ID, u.Description))
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
