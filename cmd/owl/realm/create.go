package realm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// initCreateCommand initialize the cli unit list command
func initCreateCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create realm",
		Long:    "",
		Aliases: []string{"add"},
		Example: fmt.Sprintf(`  %[1]s realm create <<< '{"ID":"test","URL":"ldap://192.168.99.108:389/dc=example,dc=org","Username":"cn=admin,dc=example,dc=org"}'`, parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			r := struct {
				ID       string
				URL      string
				Username string
			}{}

			b, err := ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = json.Unmarshal(b, &r)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = realmDriver.Set(r.ID, r.URL, r.Username)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.PrintErrf("Created realm '%v'.", r.ID)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
