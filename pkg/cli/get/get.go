package get

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/helpers/cmdutil"
	"github.com/adrienaury/owl/pkg/helpers/credentials"
	"github.com/adrienaury/owl/pkg/helpers/errutil"
	"github.com/adrienaury/owl/pkg/helpers/options"
	"github.com/adrienaury/owl/pkg/helpers/templates"
	"gopkg.in/ldap.v3"

	"github.com/spf13/cobra"
)

var (
	getLong = templates.LongDesc(`
		TODO`)

	getExample = templates.Examples(`
		TODO`)
)

// NewCommand implements the cli get command
func NewCommand(fullName string, streams options.IOStreams) *cobra.Command {
	o := NewOptions(streams)
	cmds := &cobra.Command{
		Use:     "get [filter]",
		Short:   "Print objects from LDAP server",
		Long:    getLong,
		Example: fmt.Sprintf(getExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.LoadSession(cmd, args))

			if err := o.Run(); errutil.IsUnauthorized(err) {
				fmt.Fprintln(streams.Out, "Login failed (401 Unauthorized)")
				fmt.Fprintln(streams.Out, "Verify you have provided correct credentials.")

				os.Exit(1)

			} else {
				cmdutil.CheckErr(err)
			}
		},
	}

	return cmds
}

// Run contains all the necessary functionality for the cli get command
func (o *Options) Run() error {
	conn, err := ldap.DialURL(o.Session.Server)
	if err != nil {
		return err
	}

	creds, err := credentials.GetCredentials(o.Session.Server, o.Session.Username)
	if err != nil {
		return err
	}

	err = conn.Bind(o.Session.Username, creds.Secret)
	if err != nil {
		return err
	}

	sr, err := conn.Search(ldap.NewSearchRequest("dc=example,dc=org", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"cn"}, nil))
	if err != nil {
		return err
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}

	err = o.SaveSession()
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Success\n")
	return nil
}
