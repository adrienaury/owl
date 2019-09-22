package login

import (
	"github.com/adrienaury/owl/pkg/helpers/credentials"
	"github.com/adrienaury/owl/pkg/helpers/options"
	"github.com/adrienaury/owl/pkg/helpers/paths"
	"github.com/adrienaury/owl/pkg/helpers/session"
	"github.com/spf13/cobra"
)

const defaultURL = "ldap://localhost:389"

// Options is a helper for the login and setup process, gathers all information required for a
// successful login and eventual update of config files.
// Depending on the Reader present it can be interactive, asking for terminal input in
// case of any missing information.
// Notice that some methods mutate this object so it should not be reused. The Config
// provided as a pointer will also mutate (handle new auth tokens, etc).
type Options struct {
	Server string

	// flags and printing helpers
	Username string
	Password string

	// infra
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
		// LOG
	}

	o.Session = session

	return nil
}

// SaveSession all the information present in this helper to the session file.
func (o *Options) SaveSession() error {
	err := credentials.SetCredentials(o.Session.Server, o.Session.Username, o.Password)
	if err != nil {
		return err
	}
	return o.Session.Dump(paths.Session)
}

func (o *Options) usernameProvided() bool {
	return len(o.Username) > 0
}

func (o *Options) passwordProvided() bool {
	return len(o.Password) > 0
}

func (o *Options) serverProvided() bool {
	return (len(o.Server) > 0)
}
