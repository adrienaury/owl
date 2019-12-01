package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// initApplyCommand initialize the cli unit apply command
func initApplyCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Create or update unit with provided attributes",
		Long:  "",
		Example: fmt.Sprintf(`  %[1]s unit apply <<< '{"ID": "my-unit", "Description": "Test unit"}'
  %[1]s unit apply my-unit "New description"`, parentCmd.Root().Name()),
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

			modified, err := unitDriver.Apply(unit.NewUnit(u.ID, u.Description))
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagRealm := cmd.Flag("realm")

			if modified {
				cmd.PrintErrf("Modified unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
			} else {
				cmd.PrintErrf("Created unit '%v' in realm '%v'.", u.ID, flagRealm.Value)
			}
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
