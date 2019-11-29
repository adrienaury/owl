package realm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

// newListCommand implements the cli realm list command
func newListCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List realms",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s realm list", fullName),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			realms, e := realmDriver.List()
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			b, e := json.Marshal(struct{ Realms realm.List }{realms})
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
