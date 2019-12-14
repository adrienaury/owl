package delete

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUnitCommand initialize the cli delete unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Delete units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s delete unit <<< '{"ID": "my-unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
