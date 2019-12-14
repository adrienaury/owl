package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli create user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Create users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s create user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
