package login

import (
	"errors"

	"github.com/adrienaury/owl/pkg/helpers/cmdutil"
	"github.com/spf13/cobra"
)

// Validate user inputs
func (o *Options) Validate(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Only the server URL may be specified as an argument")
	}

	serverFlag := cmdutil.GetFlagString(cmd, "server")
	if (len(serverFlag) > 0) && (len(args) == 1) {
		return errors.New("--server and passing the server URL as an argument are mutually exclusive")
	}

	return nil
}
