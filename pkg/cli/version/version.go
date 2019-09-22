package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// NewCommand provides a shim around version for
// non-client packages that require version information
func NewCommand(fullName string, version string, commit string, buildDate string, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  "Display version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "%s %v %v\n%v\n", fullName, version, commit, buildDate)
		},
	}

	return cmd
}
