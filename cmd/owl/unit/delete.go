package unit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newDeleteCommand implements the cli unit list command
func newDeleteCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [ID]",
		Short:   "Delete unit",
		Long:    "",
		Aliases: []string{"del"},
		Example: fmt.Sprintf(`  %[1]s unit delete my-unit`, fullName),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			e := unitDriver.Delete(id)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			flagRealm := cmd.Flag("realm")

			fmt.Fprintf(out, "Deleted unit '%v' in realm '%v'.", id, flagRealm.Value)
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
