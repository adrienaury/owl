package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUnitCommand initialize the cli create unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Create units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s create unit <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}