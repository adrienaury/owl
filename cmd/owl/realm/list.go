package realm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

// initListCommand initialize the cli realm list command
func initListCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List realms",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s realm list", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			realms, err := realmDriver.List()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			b, err := json.Marshal(struct{ Realms realm.List }{realms})
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.Println(string(b))
		},
	}
	parentCmd.AddCommand(cmd)
}
