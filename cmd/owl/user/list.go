package user

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

// initListCommand initialize the cli user list command
func initListCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s user list", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			users, err := userDriver.List()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			b, err := json.Marshal(struct{ Users user.List }{users})
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.Println(string(b))
		},
	}
	parentCmd.AddCommand(cmd)
}
