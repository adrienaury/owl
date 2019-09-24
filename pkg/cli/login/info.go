package login

import (
	"github.com/adrienaury/owl/pkg/helpers/credentials"
	ldap "gopkg.in/ldap.v3"
)

// GatherInfo from LDAP Server
func (o *Options) GatherInfo() error {
	conn, err := ldap.DialURL(o.Server)
	if err != nil {
		return err
	}
	defer conn.Close()

	creds, err := credentials.GetCredentials(o.Server, o.Username)
	if err != nil {
		return err
	}

	err = conn.Bind(o.Username, creds.Secret)
	if err != nil {
		return err
	}

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			"",
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=top)",
			[]string{"namingContexts"},
			nil,
		),
	)

	if err != nil {
		return err
	}

	//sr.PrettyPrint(0)
	o.BaseDN = sr.Entries[0].Attributes[0].Values[0]

	return nil
}
