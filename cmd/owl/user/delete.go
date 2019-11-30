package user

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newDeleteCommand implements the cli user delete command
func newDeleteCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [ID]",
		Short:   "Delete user",
		Long:    "",
		Aliases: []string{"del"},
		Example: fmt.Sprintf(`  %[1]s user delete joker`, fullName),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			e := userDriver.Delete(id)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			fmt.Fprintf(out, "Deleted user '%v' in unit '%v' of realm '%v'.", id, flagUnit.Value, flagRealm.Value)
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
