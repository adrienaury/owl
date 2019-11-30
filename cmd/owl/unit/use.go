package unit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newUseCommand implements the cli unit use command
func newUseCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use [ID]",
		Short:   "Use as default unit, --unit option will be implied on next commands",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s unit use my-unit", fullName),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			unit, e := unitDriver.Get(id)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			if unit == nil {
				fmt.Fprintf(err, "Unknown unit '%v'.", id)
				fmt.Fprintln(err)
				os.Exit(1)
			}

			unitDriver.Use(id)
			globalSession.Unit = unit.ID()

			fmt.Fprintf(out, "Using unit '%v' for next commands.", unit.ID())
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
