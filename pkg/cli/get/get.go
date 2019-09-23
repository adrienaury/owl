package get

import (
	"errors"
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/helpers/cmdutil"
	"github.com/adrienaury/owl/pkg/helpers/credentials"
	"github.com/adrienaury/owl/pkg/helpers/errutil"
	"github.com/adrienaury/owl/pkg/helpers/options"
	"github.com/adrienaury/owl/pkg/helpers/paths"
	"github.com/adrienaury/owl/pkg/helpers/policies"
	"github.com/adrienaury/owl/pkg/helpers/printer"
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

			if err := o.Run(args); errutil.IsUnauthorized(err) {
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
func (o *Options) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("Require exactly 1 argument : filter")
	}

	filter := args[0]

	policies, err := policies.Get(paths.Home + "/policies.yaml")
	if err != nil {
		return err
	}

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

	policy, ok := policies.Policies[filter]
	if ok {
		sr, err := conn.Search(
			ldap.NewSearchRequest(
				"dc=example,dc=org",
				ldap.ScopeWholeSubtree,
				ldap.NeverDerefAliases,
				0, 0, false,
				policy.GetFilter(),
				policy.GetAttributes(),
				nil,
			),
		)

		if err != nil {
			return err
		}

		headers := policy.GetAttributes()
		data := [][]string{}
		for _, entry := range sr.Entries {
			dataline := make([]string, len(policy.Attributes))
			for i, attr := range policy.GetAttributes() {
				dataline[i] = entry.GetAttributeValue(attr)
			}
			data = append(data, dataline)
		}
		printer.PrintData(o.Out, headers, data)
	}

	err = o.SaveSession()
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "Success\n")
	return nil
}
