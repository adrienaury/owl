package list

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

// initUnitCommand initialize the cli list unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit",
		Short:   "List units",
		Long:    "",
		Aliases: []string{"units"},
		Example: fmt.Sprintf("  %[1]s list unit", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			policyName := curRealm.Policy()

			policy, err := policyDriver.Get(policyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			units, err := unitDriver.List(policy.Objects()["unit"])
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
				fmt.Fprintf(w, "%v\t%v\n", "ID", "Description")
				for _, unit := range units.All() {
					fmt.Fprintf(w, "%v\t%v\n", unit.ID(), unit.Description())
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
