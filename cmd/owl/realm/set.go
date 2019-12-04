package realm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initSetCommand initialize the cli realm set command
func initSetCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "set [ID] [URL] [User]",
		Short:   "Set realm URL and user",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm set dev ldap://dev.my-company.com/dc=example,dc=com username", parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			url := args[1]
			username := args[2]

			err := realmDriver.Set(id, url, username)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.PrintErrf("Set realm '%v' to '%v'.", id, url)
			cmd.Println()
		},
	}
	parentCmd.AddCommand(cmd)
}
