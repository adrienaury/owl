package append

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUserCommand initialize the cli append user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Append attributes to existing users",
		Long:    "",
		Aliases: []string{"users"},
		Example: fmt.Sprintf(`  %[1]s append user <<< '{"ID": "batman", "Emails": ["bruce.wayne@gotham.dc"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
