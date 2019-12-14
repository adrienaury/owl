package upsert

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli upsert user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Create or replace attributes on users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s upsert user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
