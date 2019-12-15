package infra

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"

	"gopkg.in/ldap.v3"
)

// BackendLDAP implements a backend over LDAP (OpenLDAP schema).
type BackendLDAP struct {
	conn      *ldap.Conn
	baseDN    string
	unit      string
	userUnit  string
	groupUnit string
}

// NewBackendLDAP creates a new BackendLDAP.
func NewBackendLDAP() BackendLDAP {
	subunit := os.Getenv("OWL_BACKEND_SUBUNIT") != "false"
	userUnit := os.Getenv("OWL_BACKEND_SUBUNIT_USER")
	groupUnit := os.Getenv("OWL_BACKEND_SUBUNIT_GROUP")
	if !subunit {
		userUnit = ""
		groupUnit = ""
	} else {
		if strings.TrimSpace(userUnit) == "" {
			userUnit = "ou=users,"
		} else {
			userUnit = "ou=" + userUnit + ","
		}
		if strings.TrimSpace(groupUnit) == "" {
			groupUnit = "ou=groups,"
		} else {
			groupUnit = "ou=" + groupUnit + ","
		}
	}
	return BackendLDAP{userUnit: userUnit, groupUnit: groupUnit}
}

// OpenConnection connect to the LDAP server with given credentials.
func (b *BackendLDAP) OpenConnection(c credentials.Credentials) error {
	if b.conn != nil {
		b.conn.Close()
	}
	u, _ := url.Parse(c.URL())
	b.baseDN = strings.Trim(u.EscapedPath(), "/")
	conn, err := b.initConnection(c)
	if err != nil {
		return err
	}

	b.conn = conn
	if !b.authenticateToServer(c, b.conn) {
		return fmt.Errorf("invalid credentials")
	}

	return err
}

// UseUnit sets the default organizational unit for all operations.
func (b *BackendLDAP) UseUnit(id string) {
	b.unit = id
}

// TestCredentials tests specific credentials.
func (b BackendLDAP) TestCredentials(c credentials.Credentials) (bool, error) {
	// try to TCP connect to the server to make sure it's reachable, and discover
	// about the need of certificates or insecure TLS
	if conn, err := b.dialToServer(c); err != nil {
		switch err.(type) {
		// certificate authority unknown, check or prompt if we want an insecure
		// connection or if we already have a cluster stanza that tells us to
		// connect to this particular server insecurely
		case x509.UnknownAuthorityError, x509.HostnameError, x509.CertificateInvalidError:
			return false, err
		// TLS record header errors, like oversized record which usually means
		// the server only supports "http"
		case tls.RecordHeaderError:
			return false, err
		default:
			if _, ok := err.(*net.OpError); ok {
				return false, fmt.Errorf("%v - verify you have provided the correct host and port and that the server is currently running", err)
			}
			return false, err
		}
	} else {
		defer conn.Close()
		return b.authenticateToServer(c, conn), nil
	}
}

// ListUnits list all units contained in the LDAP server.
func (b BackendLDAP) ListUnits() (unit.List, error) {
	sr, err := b.search(b.baseDN, "(objectClass=organizationalUnit)", "ou", "description")
	if err != nil {
		return nil, err
	}

	units := make([]unit.Unit, len(sr.Entries))
	for idx, entry := range sr.Entries {
		units[idx] = unit.NewUnit(
			entry.GetAttributeValue("ou"),
			entry.GetAttributeValue("description"),
		)
	}

	return unit.NewList(units), nil
}

// GetUnit returns a specific unit with id.
func (b BackendLDAP) GetUnit(id string) (unit.Unit, error) {
	sr, err := b.search(b.baseDN, "(ou="+id+")", "ou", "description")
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("assertion failed, should have exactly 1 unit named %v but got %v", id, len(sr.Entries))
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	entry := sr.Entries[0]

	return unit.NewUnit(
		entry.GetAttributeValue("ou"),
		entry.GetAttributeValue("description"),
	), nil
}

// CreateUnit creates a new unit in the server.
func (b BackendLDAP) CreateUnit(u unit.Unit) error {
	dn := "ou=" + u.ID() + "," + b.baseDN

	addRequest := ldap.NewAddRequest(dn, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})
	addRequest.Attribute("ou", []string{u.ID()})
	addRequest.Attribute("description", []string{u.Description()})

	err := b.conn.Add(addRequest)
	if err != nil {
		return err
	}

	dnUsers := b.userUnit + dn

	addRequest = ldap.NewAddRequest(dnUsers, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})

	err = b.conn.Add(addRequest)
	if err != nil {
		return err
	}

	dnGroups := b.groupUnit + dn

	addRequest = ldap.NewAddRequest(dnGroups, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})

	err = b.conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUnit modify an existing unit.
func (b BackendLDAP) UpdateUnit(u unit.Unit) error {
	dn := "ou=" + u.ID() + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("description", []string{u.Description()})

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// AppendUnit add attributes to the unit with id.
func (b BackendLDAP) AppendUnit(u unit.Unit) error {
	dn := "ou=" + u.ID() + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	// nothing to add on this object

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// RemoveUnit remove attributes from the unit with id.
func (b BackendLDAP) RemoveUnit(u unit.Unit) error {
	dn := "ou=" + u.ID() + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	// nothing to remove on this object

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUnit deletes a unit.
func (b BackendLDAP) DeleteUnit(id string) error {
	dnUsers := b.userUnit + "ou=" + id + "," + b.baseDN
	delRequest := ldap.NewDelRequest(dnUsers, []ldap.Control{})
	err := b.conn.Del(delRequest)
	if err != nil && !ldap.IsErrorWithCode(err, 32) {
		return err
	}

	dnGroups := b.groupUnit + "ou=" + id + "," + b.baseDN
	delRequest = ldap.NewDelRequest(dnGroups, []ldap.Control{})
	err = b.conn.Del(delRequest)
	if err != nil && !ldap.IsErrorWithCode(err, 32) {
		return err
	}

	dn := "ou=" + id + "," + b.baseDN
	delRequest = ldap.NewDelRequest(dn, []ldap.Control{})
	err = b.conn.Del(delRequest)
	if err != nil && !ldap.IsErrorWithCode(err, 32) {
		return err
	}

	return nil
}

// ListUsers list all users contained in the LDAP server.
func (b BackendLDAP) ListUsers() (user.List, error) {
	var searchDN string
	if strings.TrimSpace(b.unit) == "" {
		searchDN = b.userUnit + b.baseDN
	} else {
		searchDN = b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	sr, err := b.search(searchDN, "(objectClass=inetOrgPerson)", "cn", "givenName", "sn", "mail")
	if err != nil {
		return nil, err
	}

	users := make([]user.User, len(sr.Entries))
	for idx, entry := range sr.Entries {
		users[idx] = user.NewUser(
			entry.GetAttributeValue("cn"),
			entry.GetAttributeValues("givenName"),
			entry.GetAttributeValues("sn"),
			entry.GetAttributeValues("mail"),
		)
	}

	return user.NewList(users), nil
}

// GetUser returns a specific user with id.
func (b BackendLDAP) GetUser(id string) (user.User, error) {
	var searchDN string
	if strings.TrimSpace(b.unit) == "" {
		searchDN = b.userUnit + b.baseDN
	} else {
		searchDN = b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	sr, err := b.search(searchDN, "(cn="+id+")", "cn", "givenName", "sn", "mail")
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("assertion failed, should have exactly 1 user named %v but got %v", id, len(sr.Entries))
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	entry := sr.Entries[0]

	return user.NewUser(
		entry.GetAttributeValue("cn"),
		entry.GetAttributeValues("givenName"),
		entry.GetAttributeValues("sn"),
		entry.GetAttributeValues("mail"),
	), nil
}

// CreateUser creates a new user in the default unit.
func (b BackendLDAP) CreateUser(u user.User) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + u.ID() + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + u.ID() + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	addRequest := ldap.NewAddRequest(dn, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"inetOrgPerson"})
	addRequest.Attribute("cn", []string{u.ID()})
	if len(u.FirstNames()) > 0 {
		addRequest.Attribute("givenName", u.FirstNames())
	}
	if len(u.LastNames()) > 0 {
		addRequest.Attribute("sn", u.LastNames())
	}
	if len(u.Emails()) > 0 {
		addRequest.Attribute("mail", u.Emails())
	}

	err := b.conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser modify an existing user.
func (b BackendLDAP) UpdateUser(u user.User) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + u.ID() + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + u.ID() + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("givenName", u.FirstNames())
	modRequest.Replace("sn", u.LastNames())
	modRequest.Replace("mail", u.Emails())

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user.
func (b BackendLDAP) DeleteUser(id string) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + id + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + id + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	delRequest := ldap.NewDelRequest(dn, []ldap.Control{})
	err := b.conn.Del(delRequest)
	if err != nil && !ldap.IsErrorWithCode(err, 32) {
		return err
	}

	return nil
}

// AppendUser add attributes to the user with id.
func (b BackendLDAP) AppendUser(u user.User) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + u.ID() + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + u.ID() + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	if len(u.FirstNames()) > 0 {
		modRequest.Add("givenName", u.FirstNames())
	}
	if len(u.LastNames()) > 0 {
		modRequest.Add("sn", u.LastNames())
	}
	if len(u.Emails()) > 0 {
		modRequest.Add("mail", u.Emails())
	}

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// RemoveUser remove attributes from the user with id.
func (b BackendLDAP) RemoveUser(u user.User) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + u.ID() + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + u.ID() + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	if len(u.FirstNames()) > 0 {
		modRequest.Delete("givenName", u.FirstNames())
	}
	if len(u.LastNames()) > 0 {
		modRequest.Delete("sn", u.LastNames())
	}
	if len(u.Emails()) > 0 {
		modRequest.Delete("mail", u.Emails())
	}

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// GetPrincipalEmail returns the main e-mail (of index 0) of a specific unit with id.
func (b BackendLDAP) GetPrincipalEmail(userID string) (string, error) {
	user, err := b.GetUser(userID)
	if err != nil {
		return "", err
	}
	if len(user.Emails()) == 0 {
		return "", fmt.Errorf("user has no e-mail")
	}
	return user.Emails()[0], nil
}

// SetUserPassword sets hashed password to user.
func (b BackendLDAP) SetUserPassword(userID string, hashedPassword string) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + userID + "," + b.userUnit + b.baseDN
	} else {
		dn = "cn=" + userID + "," + b.userUnit + "ou=" + b.unit + "," + b.baseDN
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("userPassword", []string{hashedPassword})

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// ListGroups list all groups contained in the LDAP server.
func (b BackendLDAP) ListGroups() (group.List, error) {
	var searchDN string
	if strings.TrimSpace(b.unit) == "" {
		searchDN = b.groupUnit + b.baseDN
	} else {
		searchDN = b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	sr, err := b.search(searchDN, "(objectClass=groupOfUniqueNames)", "cn", "uniqueMember")
	if err != nil {
		return nil, err
	}

	userDn := b.userUnit + "ou=" + b.unit + "," + b.baseDN
	groups := make([]group.Group, len(sr.Entries))
	for idx, entry := range sr.Entries {
		usersInGroup := []string{}
		for _, user := range entry.GetAttributeValues("uniqueMember") {
			if strings.HasSuffix(user, userDn) {
				usersInGroup = append(usersInGroup, strings.TrimPrefix(strings.TrimSuffix(user, ","+userDn), "cn="))
			}
		}
		groups[idx] = group.NewGroup(
			entry.GetAttributeValue("cn"),
			usersInGroup...,
		)
	}

	return group.NewList(groups), nil
}

// GetGroup returns a specific group with id.
func (b BackendLDAP) GetGroup(id string) (group.Group, error) {
	var searchDN string
	if strings.TrimSpace(b.unit) == "" {
		searchDN = b.groupUnit + b.baseDN
	} else {
		searchDN = b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	sr, err := b.search(searchDN, "(cn="+id+")", "cn", "uniqueMember")
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("assertion failed, should have exactly 1 group named %v but got %v", id, len(sr.Entries))
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	entry := sr.Entries[0]
	members := entry.GetAttributeValues("uniqueMember")

	userDn := b.userUnit + "ou=" + b.unit + "," + b.baseDN
	userIDs := make([]string, len(members))
	for idx, member := range members {
		if strings.HasSuffix(member, userDn) {
			userIDs[idx] = strings.TrimPrefix(strings.TrimSuffix(member, ","+userDn), "cn=")
		}
	}

	return group.NewGroup(
		entry.GetAttributeValue("cn"),
		userIDs...,
	), nil
}

// CreateGroup creates a new group in the default unit.
func (b BackendLDAP) CreateGroup(g group.Group) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + g.ID() + "," + b.groupUnit + b.baseDN
	} else {
		dn = "cn=" + g.ID() + "," + b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	members := []string{}
	for _, member := range g.Members() {
		members = append(members, "cn="+member+","+b.userUnit+"ou="+b.unit+","+b.baseDN)
	}

	addRequest := ldap.NewAddRequest(dn, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"groupOfUniqueNames"})
	addRequest.Attribute("cn", []string{g.ID()})
	addRequest.Attribute("uniqueMember", members)

	err := b.conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGroup modify an existing group.
func (b BackendLDAP) UpdateGroup(g group.Group) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + g.ID() + "," + b.groupUnit + b.baseDN
	} else {
		dn = "cn=" + g.ID() + "," + b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	members := []string{}
	for _, member := range g.Members() {
		members = append(members, "cn="+member+","+b.userUnit+"ou="+b.unit+","+b.baseDN)
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("uniqueMember", members)

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroup deletes a group.
func (b BackendLDAP) DeleteGroup(id string) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + id + "," + b.groupUnit + b.baseDN
	} else {
		dn = "cn=" + id + "," + b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	delRequest := ldap.NewDelRequest(dn, []ldap.Control{})

	err := b.conn.Del(delRequest)
	if err != nil && !ldap.IsErrorWithCode(err, 32) {
		return err
	}

	return nil
}

// AppendGroup add attributes to the group with id.
func (b BackendLDAP) AppendGroup(g group.Group) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + g.ID() + "," + b.groupUnit + b.baseDN
	} else {
		dn = "cn=" + g.ID() + "," + b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	members := []string{}
	for _, member := range g.Members() {
		members = append(members, "cn="+member+","+b.userUnit+"ou="+b.unit+","+b.baseDN)
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Add("uniqueMember", members)

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// RemoveGroup remove attributes from the group with id.
func (b BackendLDAP) RemoveGroup(g group.Group) error {
	var dn string
	if strings.TrimSpace(b.unit) == "" {
		dn = "cn=" + g.ID() + "," + b.groupUnit + b.baseDN
	} else {
		dn = "cn=" + g.ID() + "," + b.groupUnit + "ou=" + b.unit + "," + b.baseDN
	}

	members := []string{}
	for _, member := range g.Members() {
		members = append(members, "cn="+member+","+b.userUnit+"ou="+b.unit+","+b.baseDN)
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Delete("uniqueMember", members)

	err := b.conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

func (b BackendLDAP) search(dn string, filter string, attributes ...string) (*ldap.SearchResult, error) {
	sr, err := b.conn.Search(
		ldap.NewSearchRequest(
			dn,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			filter,
			attributes,
			nil,
		),
	)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func (b BackendLDAP) initConnection(c credentials.Credentials) (*ldap.Conn, error) {
	if c == nil {
		return nil, fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(c)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// dialToServer takes the Server URL and dials to make sure the server is reachable.
func (b BackendLDAP) dialToServer(c credentials.Credentials) (*ldap.Conn, error) {
	l, err := ldap.DialURL(c.URL())
	if err != nil {
		return nil, err
	}
	return l, nil
}

// authenticateToServer takes the provided credentials try to bind to make sure they are correct.
func (b BackendLDAP) authenticateToServer(c credentials.Credentials, conn *ldap.Conn) bool {
	return conn.Bind(c.Username(), c.Password()) == nil
}
