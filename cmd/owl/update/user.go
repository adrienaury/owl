package update

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli update user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Replace specified attributes on users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s update user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
