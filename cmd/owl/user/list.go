package user

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/adrienaury/owl/pkg/domain/user"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initListCommand initialize the cli user list command
func initListCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users",
		Long:    "",
		Aliases: []string{"ls"},
		Example: fmt.Sprintf("  %[1]s user list", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			users, err := userDriver.List()
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
				b, err := json.Marshal(struct{ Users user.List }{users})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "yaml":
				b, err := yaml.Marshal(struct{ Users user.List }{users})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				cmd.Println(string(b))
			case "table":
				w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "ID", "First Names", "Last Names", "E-mails")
				for _, user := range users.All() {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", user.ID(), strings.Join(user.FirstNames(), ", "), strings.Join(user.LastNames(), ", "), strings.Join(user.Emails(), ", "))
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
