package login

import (
	"errors"
	"fmt"

	"github.com/adrienaury/owl/pkg/helpers/cmdutil"
	"github.com/adrienaury/owl/pkg/helpers/credentials"
	"github.com/adrienaury/owl/pkg/helpers/term"
	"github.com/spf13/cobra"
)

// Complete user inputs with missing informations
func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	if err := o.completeServer(cmd, args); err != nil {
		return err
	}

	if err := o.completeUsername(cmd, args); err != nil {
		return err
	}

	if err := o.completePassword(cmd, args); err != nil {
		return err
	}

	return nil
}

func (o *Options) completeServer(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		o.Server = args[0]
	} else if serverFlag := cmdutil.GetFlagString(cmd, "server"); serverFlag != "" {
		o.Server = serverFlag
	} else if o.Session.Server != "" {
		o.Server = o.Session.Server
	} else if term.IsTerminal(o.In) {
		for !o.serverProvided() {
			defaultServer := defaultURL
			promptMsg := fmt.Sprintf("Server [%s]: ", defaultServer)
			o.Server = term.PromptForStringWithDefault(o.In, o.Out, defaultServer, promptMsg)
		}
	} else {
		return errors.New("A server URL must be specified")
	}
	return nil
}

func (o *Options) completeUsername(cmd *cobra.Command, args []string) error {
	if !o.usernameProvided() {
		if o.Session.Username != "" {
			o.Username = o.Session.Username
		} else if term.IsTerminal(o.In) {
			for !o.usernameProvided() {
				o.Username = term.PromptForString(o.In, o.Out, "Username: ")
			}
		} else {
			return errors.New("A username must be specified")
		}
	}
	return nil
}

func (o *Options) completePassword(cmd *cobra.Command, args []string) error {
	if !o.passwordProvided() {
		if creds, err := credentials.GetCredentials(o.Server, o.Username); err == nil {
			o.Password = creds.Secret
		} else if credentials.IsErrCredentialsNotFound(err) {
			if term.IsTerminal(o.In) {
				for !o.passwordProvided() {
					o.Password = term.PromptForPasswordString(o.In, o.Out, "Password: ")
				}
			}
		} else {
			return errors.New("A password must be specified")
		}
	}
	return nil
}
