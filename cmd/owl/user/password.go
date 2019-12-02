package user

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/password"
	"github.com/spf13/cobra"
)

// initPasswordCommand initialize the cli user password command
func initPasswordCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "new-password [ID]",
		Short:   "Assign a new password to user",
		Long:    "",
		Aliases: []string{"passwd"},
		Example: fmt.Sprintf(`  %[1]s user new-password joker`, parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			err := passwordDriver.AssignRandomPassword("SSHA256", password.Standard, 12, id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Assigned new random password to user '%v' in unit '%v' of realm '%v'.", id, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
