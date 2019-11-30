package unit

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// newListCommand implements the cli unit list command
func newListCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List units",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s unit list", fullName),
		Args:    cobra.NoArgs,
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			units, e := unitDriver.List()
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			b, e := json.Marshal(struct{ Units unit.List }{units})
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
