package cmdutil

import (
	"github.com/spf13/cobra"
)

// GetFlagString todo
func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		// TODO log
		// fmt.Errorf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}
