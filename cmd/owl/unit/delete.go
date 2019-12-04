package unit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initDeleteCommand initialize the cli unit delete command
func initDeleteCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete [ID]",
		Short:   "Delete unit",
		Long:    "",
		Aliases: []string{"del"},
		Example: fmt.Sprintf(`  %[1]s unit delete my-unit`, parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			err := unitDriver.Delete(id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Deleted unit '%v' in realm '%v'.", id, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
