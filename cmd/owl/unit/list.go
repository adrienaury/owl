package unit

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
)

// initListCommand initialize the cli unit list command
func initListCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List units",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s unit list", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			units, err := unitDriver.List()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			b, err := json.Marshal(struct{ Units unit.List }{units})
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			cmd.Println(string(b))
		},
	}
	parentCmd.AddCommand(cmd)
}
