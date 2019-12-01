package group

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initDeleteCommand initialize the cli group delete command
func initDeleteCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete [ID]",
		Short:   "Delete group",
		Long:    "",
		Aliases: []string{"del"},
		Example: fmt.Sprintf(`  %[1]s group delete my-group`, parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			err := groupDriver.Delete(id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Deleted group '%v' in unit '%v' of realm '%v'.", id, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
