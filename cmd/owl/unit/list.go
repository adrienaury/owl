package unit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

			output := "json"
			flagOutput := cmd.Flag("output")
			if flagOutput != nil && strings.TrimSpace(flagOutput.Value.String()) != "" {
				output = flagOutput.Value.String()
			}

			switch output {
			case "json":
				b, err := json.Marshal(struct{ Units unit.List }{units})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "yaml":
				b, err := yaml.Marshal(struct{ Units unit.List }{units})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "table":
				w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "%v\n", "ID")
				for _, unit := range units.All() {
					fmt.Fprintf(w, "%v\n", unit.ID())
				}
				w.Flush()
			default:
				cmd.PrintErrf("Invalid output format : %v", output)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
