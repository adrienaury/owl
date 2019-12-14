package apply

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initUnitCommand initialize the cli apply unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "Create or replace units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf(`  %[1]s apply unit <<< '{"ID": "my-unit", "Description": "Test unit"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	parentCmd.AddCommand(cmd)
}
