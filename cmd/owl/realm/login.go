package realm

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// newLoginCommand implements the cli realm set command
func newLoginCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login [ID]",
		Short:   "Login to the realm, --realm option will be implied on next commands",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s realm login dev", fullName),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			realm, e := realmDriver.Get(id)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			if realm == nil {
				fmt.Fprintf(err, "Unknown realm '%v'.", id)
				fmt.Fprintln(err)
				os.Exit(1)
			}

			// get existing crendentials
			creds, e := credentialsDriver.Get(realm.URL(), realm.Username())
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			// no existing credentials, ask password to user
			if creds == nil {
				creds = credentials.NewCredentials(realm.URL(), realm.Username(), askPassword(err, out, in))
			}

			// test credentials
			ok, e := credentialsDriver.Test(creds)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			// credentials are invalid, ask password to user again
			if !ok {
				creds = credentials.NewCredentials(creds.URL(), creds.Username(), askPassword(err, out, in))
			}

			// test credentials (again)
			ok, e = credentialsDriver.Test(creds)
			if e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			// credentials are invalid, report error and exit
			if !ok {
				fmt.Fprintln(err, "Invalid credentials.")
				os.Exit(1)
			}

			// success! store credentials for next time
			if e := credentialsDriver.Set(realm.URL(), realm.Username(), creds.Password()); e != nil {
				fmt.Fprintln(err, e.Error())
				os.Exit(1)
			}

			credentialsDriver.Use(creds)
			globalSession.Realm = realm.ID()

			fmt.Fprintf(out, "Connected to realm '%v' as user '%v'.", realm.ID(), realm.Username())
			fmt.Fprintln(out)
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}

func askPassword(err *os.File, out *os.File, in *os.File) string {
	fmt.Fprint(out, "Password: ")
	bytePassword, e := terminal.ReadPassword(int(in.Fd()))
	fmt.Fprintln(out)
	if e != nil {
		fmt.Fprintln(err, e.Error())
		os.Exit(1)
	}
	return string(bytePassword)
}
