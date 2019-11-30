package unit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initUseCommand initialize the cli unit use command
func initUseCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "use [ID]",
		Short:   "Use default unit, --unit option will be implied on next commands",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s unit use my-unit", parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			unit, err := unitDriver.Get(id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if unit == nil {
				cmd.PrintErrf("Unknown unit '%v'.", id)
				cmd.PrintErrln()
				os.Exit(1)
			}

			unitDriver.Use(id)
			globalSession.Unit = unit.ID()

			cmd.PrintErrf("Using unit '%v' for next commands.", unit.ID())
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
