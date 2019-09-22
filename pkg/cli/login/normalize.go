package login

import (
	"github.com/adrienaury/owl/pkg/helpers/normalizer"
	"github.com/spf13/cobra"
)

// Normalize user inputs
func (o *Options) Normalize(cmd *cobra.Command, args []string, commandName string) error {
	// normalize the provided server to a format expected by config
	serverNormalized, err := normalizer.NormalizeLdapServerURL(o.Server)
	if err != nil {
		return err
	}
	o.Server = serverNormalized

	o.CommandName = commandName

	return nil
}
