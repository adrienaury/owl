package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
)

// newCreateCommand implements the cli user list command
func newCreateCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create user",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'`, fullName),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			u := struct {
				ID         string
				FirstNames []string
				LastNames  []string
				Emails     []string
				Groups     []string
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

			e = userDriver.Create(user.NewUser(u.ID, u.FirstNames, u.LastNames, u.Emails, u.Groups))
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			fmt.Fprintf(out, "Created user '%v' in unit '%v' of realm '%v'.", u.ID, flagUnit.Value, flagRealm.Value)
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
