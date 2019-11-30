package realm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initDeleteCommand initialize the cli realm delete command
func initDeleteCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete [ID]",
		Short:   "Delete realm",
		Long:    "",
		Aliases: []string{"del"},
		Example: fmt.Sprintf(`  %[1]s realm delete my-realm`, parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			realm, err := realmDriver.Get(id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if realm == nil {
				cmd.PrintErrf("Unknown realm '%v'.", id)
				cmd.PrintErrln()
				os.Exit(1)
			}

			err = realmDriver.Delete(realm.ID())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.PrintErrf("Deleted realm '%v'.", id)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
