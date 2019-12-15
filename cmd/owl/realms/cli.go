package realms

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
)

var (
	realmDriver realm.Driver
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(r realm.Driver) {
	realmDriver = r
}

// InitCommand initialize the cli realm command
func InitCommand(parentCmd *cobra.Command) {
	realmCmd := &cobra.Command{
		Use:     "realms",
		Short:   "List realms",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realms", parentCmd.Root().Name()),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			realms, err := realmDriver.List()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			currentRealm := ""
			flagRealm := cmd.Flag("realm")
			if flagRealm != nil {
				currentRealm = flagRealm.Value.String()
			}
			w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 2, ' ', 0)
			for _, realm := range realms.All() {
				if realm.ID() == currentRealm {
					fmt.Fprintf(w, "%v*\t%v\t%v\n", realm.ID(), realm.Username(), realm.URL())
				} else {
					fmt.Fprintf(w, "%v\t%v\t%v\n", realm.ID(), realm.Username(), realm.URL())
				}
			}
			w.Flush()
		},
	}
	parentCmd.AddCommand(realmCmd)
}
