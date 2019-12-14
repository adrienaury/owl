package apply

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli apply group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Create or replace groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s apply group <<< '{"ID": "good-guys", "Members": ["batman"]}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
