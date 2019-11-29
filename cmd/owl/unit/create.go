package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// newCreateCommand implements the cli unit list command
func newCreateCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create unit",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s unit create <<< '{"id": "my-unit"}'`, fullName),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			u := struct {
				ID string
			}{}

			b, e := ioutil.ReadAll(in)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			e = json.Unmarshal(b, &u)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			e = unitDriver.Create(unit.NewUnit(u.ID))
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
