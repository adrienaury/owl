package remove

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUnitCommand initialize the cli remove unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Remove attributes from existing units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s remove unit <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
