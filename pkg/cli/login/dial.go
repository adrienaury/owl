package login

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/adrienaury/owl/pkg/helpers/errutil"
	"github.com/adrienaury/owl/pkg/helpers/term"
	ldap "gopkg.in/ldap.v3"
)

// DialAndAuthenticate try to connect and authenticate to the server.
func (o *Options) DialAndAuthenticate() error {

	// try to TCP connect to the server to make sure it's reachable, and discover
	// about the need of certificates or insecure TLS
	if conn, err := o.dialToServer(); err != nil {
		switch err.(type) {
		// certificate authority unknown, check or prompt if we want an insecure
		// connection or if we already have a cluster stanza that tells us to
		// connect to this particular server insecurely
		case x509.UnknownAuthorityError, x509.HostnameError, x509.CertificateInvalidError:
			return errutil.GetPrettyErrorForServer(err, o.Server)
		// TLS record header errors, like oversized record which usually means
		// the server only supports "http"
		case tls.RecordHeaderError:
			return errutil.GetPrettyErrorForServer(err, o.Server)
		default:
			if _, ok := err.(*net.OpError); ok {
				return fmt.Errorf("%v - verify you have provided the correct host and port and that the server is currently running", err)
			}
			return err
		}
	} else {
		defer conn.Close()
		return o.authenticateToServer(conn)
	}
}

// dialToServer takes the Server URL and dials to make sure the server is reachable.
func (o *Options) dialToServer() (*ldap.Conn, error) {

	l, err := ldap.DialURL(o.Server)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// authenticateToServer takes the provided credentials try to bind to make sure they are correct.
func (o *Options) authenticateToServer(conn *ldap.Conn) error {
	err := conn.Bind(o.Username, o.Password)
	if err != nil {
		if term.IsTerminal(o.In) {
			fmt.Fprintf(o.Out, "Invalid credentials!\n")
			for !o.usernameProvided() {
				o.Username = term.PromptForString(o.In, o.Out, "Username: ")
			}
			for !o.passwordProvided() {
				o.Password = term.PromptForPasswordString(o.In, o.Out, "Password: ")
			}
			return o.authenticateToServer(conn)
		}
		err = fmt.Errorf("%v - verify you have provided the correct credentials", err)
	}

	return err
}
