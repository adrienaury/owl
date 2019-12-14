package update

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli update group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Replace specified attributes on groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s update group <<< '{"ID": "good-guys", "Members": ["batman"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
