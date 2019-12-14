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

// initUserCommand initialize the cli get user command
func initUserCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "user ID",
		Short:   "Get user",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s get user batman", parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			user, err := userDriver.Get(id)
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
				b, err := json.Marshal(user)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "yaml":
				b, err := yaml.Marshal(user)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "table":
				w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "ID", "First Names", "Last Names", "E-mails")
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", user.ID(), strings.Join(user.FirstNames(), ", "), strings.Join(user.LastNames(), ", "), strings.Join(user.Emails(), ", "))
				w.Flush()
			default:
				cmd.PrintErrf("Invalid output format : %v", output)
				cmd.PrintErrln()
			}
		},
	}
	parentCmd.AddCommand(cmd)
}
