package delete

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initGroupCommand initialize the cli delete group command
func initGroupCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Delete groups",
		Long:    "",
		Aliases: []string{"groups"},
		Example: fmt.Sprintf(`  %[1]s delete group <<< '{"ID": "good-guys"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
