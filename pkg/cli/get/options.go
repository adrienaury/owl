package get

import (
	"github.com/adrienaury/owl/pkg/helpers/options"
	"github.com/adrienaury/owl/pkg/helpers/paths"
	"github.com/adrienaury/owl/pkg/helpers/session"
	"github.com/spf13/cobra"
)

// Options is a helper
type Options struct {
	*session.Session

	CommandName string

	options.IOStreams
}

// NewOptions todo
func NewOptions(streams options.IOStreams) *Options {
	return &Options{
		IOStreams: streams,
	}
}

// LoadSession get current session from disk (if it exists)
func (o *Options) LoadSession(cmd *cobra.Command, args []string) error {
	session, err := session.NewSession().Restore(paths.Session)
	if err != nil {
		return err
	}

	o.Session = session

	return nil
}

// SaveSession all the information present in this helper to the session file.
func (o *Options) SaveSession() error {
	return o.Session.Dump(paths.Session)
}
