package upsert

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli upsert group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Create or replace attributes on groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s upsert group <<< '{"ID": "good-guys", "Members": ["batman"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
