package login

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/realm"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	globalSession     *session.Session
	realmDriver       realm.Driver
	credentialsDriver credentials.Driver
)

// SetDrivers inject required domain drivers in the command.
func SetDrivers(r realm.Driver, c credentials.Driver) {
	realmDriver = r
	credentialsDriver = c
}

// SetSession inject the global session in the command.
func SetSession(s *session.Session) {
	globalSession = s
}

// InitCommand initialize the cli login command
func InitCommand(parentCmd *cobra.Command) {
	loginCmd := &cobra.Command{
		Use:     "login [ID or URL]",
		Short:   "Login to realm",
		Long:    "Login to realm which means --realm option will be implied on next commands.",
		Example: fmt.Sprintf("  %[1]s realm login dev", parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			realm, err := realmDriver.Get(id)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if realm == nil {
				cmd.PrintErrf("Unknown realm '%v'.", id)
				cmd.PrintErrln()
				os.Exit(1)
			}

			// get existing crendentials
			creds, err := credentialsDriver.Get(realm.URL(), realm.Username())
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			// no existing credentials, ask password to user
			if creds == nil {
				creds = credentials.NewCredentials(realm.URL(), realm.Username(), askPassword(cmd))
			}

			// test credentials
			ok, err := credentialsDriver.Test(creds)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			// credentials are invalid, ask password to user again
			if !ok {
				creds = credentials.NewCredentials(creds.URL(), creds.Username(), askPassword(cmd))
			}

			// test credentials (again)
			ok, err = credentialsDriver.Test(creds)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			// credentials are invalid, report error and exit
			if !ok {
				cmd.PrintErrln("Invalid credentials.")
				os.Exit(1)
			}

			// success! store credentials for next time
			if err := credentialsDriver.Set(realm.URL(), realm.Username(), creds.Password()); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if err := credentialsDriver.Use(creds); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			globalSession.Realm = realm.ID()

			cmd.PrintErrf("Connected to realm '%v' as user '%v'.", realm.ID(), realm.Username())
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(loginCmd)
}

func askPassword(cmd *cobra.Command) string {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		cmd.Printf("Password: ")
		bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		cmd.Println()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		return string(bytePassword)
	}
	cmd.PrintErrln("No credentials.")
	os.Exit(1)
	return ""
}
