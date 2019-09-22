package login

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/helpers/cmdutil"
	"github.com/adrienaury/owl/pkg/helpers/errutil"
	"github.com/adrienaury/owl/pkg/helpers/normalizer"
	"github.com/adrienaury/owl/pkg/helpers/options"
	"github.com/adrienaury/owl/pkg/helpers/templates"

	"github.com/spf13/cobra"
)

var (
	loginLong = templates.LongDesc(`
		Log in to your server and save login for subsequent use
		First-time users of the client should run this command to connect to a server,
		establish an authenticated session, and save connection to the configuration file. The
		default configuration will be saved to your home directory under ".owl/config".
		The information required to login -- like username and password, or
		the server details -- can be provided through flags. If not provided, the command will
		prompt for user input as needed.`)

	loginExample = templates.Examples(`
		# Log in interactively
		%[1]s login
		# Log in to the given server with the given credentials (will not prompt interactively)
		%[1]s login localhost:389 --username=myuser --password=mypass`)
)

// NewCommand implements the cli login command
func NewCommand(fullName string, streams options.IOStreams) *cobra.Command {
	o := NewOptions(streams)
	cmds := &cobra.Command{
		Use:     "login [URL]",
		Short:   "Log in to a server",
		Long:    loginLong,
		Example: fmt.Sprintf(loginExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.LoadSession(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Normalize(cmd, args, fullName))
			cmdutil.CheckErr(o.Validate(cmd, args))

			if err := o.Run(); errutil.IsUnauthorized(err) {
				fmt.Fprintln(streams.Out, "Login failed (401 Unauthorized)")
				fmt.Fprintln(streams.Out, "Verify you have provided correct credentials.")

				os.Exit(1)

			} else {
				cmdutil.CheckErr(err)
			}
		},
	}

	// Login is the only command that can negotiate a session token against the auth server using basic auth
	cmds.Flags().StringVarP(&o.Username, "username", "u", o.Username, "Username, will prompt if not provided")
	cmds.Flags().StringVarP(&o.Password, "password", "p", o.Password, "Password, will prompt if not provided")

	return cmds
}

// Run contains all the necessary functionality for the cli login command
func (o *Options) Run() error {
	if err := o.DialAndAuthenticate(); err != nil {
		return err
	}

	normalizedURL, err := normalizer.NormalizeLdapServerURL(o.Server)
	if err != nil {
		return err
	}

	o.Session.Server = normalizedURL
	o.Session.Username = o.Username

	err = o.SaveSession()
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Welcome! See '%s help' to get started.\n", o.CommandName)
	return nil
}
