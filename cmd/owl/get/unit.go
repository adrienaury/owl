package get

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initUnitCommand initialize the cli get unit command
func initUnitCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "unit ID",
		Short:   "Get unit",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s get unit my-unit", parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			policyName := curRealm.Policy()

			policy, err := policyDriver.Get(policyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			unit, err := unitDriver.Get(id, policy.Objects()["unit"])
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
				b, err := json.Marshal(unit)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "yaml":
				b, err := yaml.Marshal(unit)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "table":
				w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "%v\t%v\n", "ID", "Description")
				fmt.Fprintf(w, "%v\t%v\n", unit.ID(), unit.Description())
				w.Flush()
			default:
				cmd.PrintErrf("Invalid output format : %v", output)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
