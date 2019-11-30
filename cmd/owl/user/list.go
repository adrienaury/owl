package user

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

// newListCommand implements the cli user list command
func newListCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s user list", fullName),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			users, e := userDriver.List()
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			b, e := json.Marshal(struct{ Users user.List }{users})
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			fmt.Fprintln(out, string(b))
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
