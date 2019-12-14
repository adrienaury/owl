package realm

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

var (
	realmDriver realm.Driver
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(r realm.Driver) {
	realmDriver = r
}

// InitCommand initialize the cli realm command
func InitCommand(parentCmd *cobra.Command) {
	realmCmd := &cobra.Command{
		Use:     "realm [ID] [URL] [User]",
		Short:   "Configure realms",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm dev ldap://dev.my-company.com/dc=example,dc=com username", parentCmd.Root().Name()),
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
	parentCmd.AddCommand(realmCmd)
}
