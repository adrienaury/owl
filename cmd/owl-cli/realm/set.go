package realm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newSetCommand implements the cli realm set command
func newSetCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [ID] [URL] [User]",
		Short:   "Set realm URL and user",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm set dev ldap://dev.my-company.com/dc=example,dc=com username", fullName),
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			url := args[1]
			username := args[2]

			e := realmDriver.Set(id, url, username)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			fmt.Fprintf(out, "Set realm '%v' to '%v'.", id, url)
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
