package apply

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli apply user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Create or replace users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s apply user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
