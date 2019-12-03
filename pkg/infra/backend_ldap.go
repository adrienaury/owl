package infra

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/adrienaury/owl/pkg/domain/group"
	"github.com/adrienaury/owl/pkg/domain/unit"
	"github.com/adrienaury/owl/pkg/domain/user"

	"gopkg.in/ldap.v3"
)

// BackendLDAP ...
type BackendLDAP struct {
	creds  credentials.Credentials
	baseDN string
	unit   string
}

// NewBackendLDAP ...
func NewBackendLDAP() BackendLDAP {
	return BackendLDAP{}
}

// UseCredentials ...
func (b *BackendLDAP) UseCredentials(c credentials.Credentials) {
	b.creds = c
	u, _ := url.Parse(c.URL())
	b.baseDN = strings.Trim(u.EscapedPath(), "/")
}

// UseUnit ...
func (b *BackendLDAP) UseUnit(id string) {
	b.unit = id
}

// TestCredentials ...
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

// ListUnits ...
func (b BackendLDAP) ListUnits() (unit.List, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			b.baseDN,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=organizationalUnit)",
			[]string{"ou", "description"},
			nil,
		),
	)
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

// GetUnit ...
func (b BackendLDAP) GetUnit(id string) (unit.Unit, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			b.baseDN,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(ou="+id+")",
			[]string{"ou", "description"},
			nil,
		),
	)
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

// CreateUnit ...
func (b BackendLDAP) CreateUnit(u unit.Unit) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "ou=" + u.ID() + "," + b.baseDN

	addRequest := ldap.NewAddRequest(dn, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})
	addRequest.Attribute("ou", []string{u.ID()})
	addRequest.Attribute("description", []string{u.Description()})

	err = conn.Add(addRequest)
	if err != nil {
		return err
	}

	dnUsers := "ou=users," + dn

	addRequest = ldap.NewAddRequest(dnUsers, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})
	addRequest.Attribute("ou", []string{"users"})

	err = conn.Add(addRequest)
	if err != nil {
		return err
	}

	dnGroups := "ou=groups," + dn

	addRequest = ldap.NewAddRequest(dnGroups, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"organizationalUnit"})
	addRequest.Attribute("ou", []string{"groups"})

	err = conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUnit ...
func (b BackendLDAP) UpdateUnit(u unit.Unit) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "ou=" + u.ID() + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("description", []string{u.Description()})

	err = conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUnit ...
func (b BackendLDAP) DeleteUnit(id string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dnUsers := "ou=users,ou=" + id + "," + b.baseDN
	delRequest := ldap.NewDelRequest(dnUsers, []ldap.Control{})
	err = conn.Del(delRequest)
	if err != nil {
		return err
	}

	dnGroups := "ou=groups,ou=" + id + "," + b.baseDN
	delRequest = ldap.NewDelRequest(dnGroups, []ldap.Control{})
	err = conn.Del(delRequest)
	if err != nil {
		return err
	}

	dn := "ou=" + id + "," + b.baseDN
	delRequest = ldap.NewDelRequest(dn, []ldap.Control{})
	err = conn.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}

// ListUsers ...
func (b BackendLDAP) ListUsers() (user.List, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return nil, fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	dn := "ou=users,ou=" + b.unit + "," + b.baseDN

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			dn,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=inetOrgPerson)",
			[]string{"cn", "givenName", "sn", "mail"},
			nil,
		),
	)
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

// GetUser ..
func (b BackendLDAP) GetUser(id string) (user.User, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return nil, fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	dn := "ou=users,ou=" + b.unit + "," + b.baseDN

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			dn,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(cn="+id+")",
			[]string{"cn", "givenName", "sn", "mail"},
			nil,
		),
	)
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

// CreateUser ...
func (b BackendLDAP) CreateUser(u user.User) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + u.ID() + ",ou=users,ou=" + b.unit + "," + b.baseDN

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

	err = conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser ...
func (b BackendLDAP) UpdateUser(u user.User) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + u.ID() + ",ou=users,ou=" + b.unit + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	if len(u.FirstNames()) > 0 {
		modRequest.Replace("givenName", u.FirstNames())
	}
	if len(u.LastNames()) > 0 {
		modRequest.Replace("sn", u.LastNames())
	}
	if len(u.Emails()) > 0 {
		modRequest.Replace("mail", u.Emails())
	}

	err = conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser ...
func (b BackendLDAP) DeleteUser(id string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + id + ",ou=users,ou=" + b.unit + "," + b.baseDN

	delRequest := ldap.NewDelRequest(dn, []ldap.Control{})

	err = conn.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}

// SetUserPassword ...
func (b BackendLDAP) SetUserPassword(userID string, hashedPassword string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + userID + ",ou=users,ou=" + b.unit + "," + b.baseDN

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Replace("userPassword", []string{hashedPassword})

	err = conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// ListGroups ...
func (b BackendLDAP) ListGroups() (group.List, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return nil, fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	dn := "ou=groups,ou=" + b.unit + "," + b.baseDN

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			dn,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=groupOfUniqueNames)",
			[]string{"cn", "uniqueMember"},
			nil,
		),
	)
	if err != nil {
		return nil, err
	}

	userDn := "ou=users,ou=" + b.unit + "," + b.baseDN
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

// GetGroup ...
func (b BackendLDAP) GetGroup(id string) (group.Group, error) {
	if b.creds == nil {
		return nil, fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return nil, fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return nil, fmt.Errorf("invalid credentials")
	}

	dn := "ou=groups,ou=" + b.unit + "," + b.baseDN

	sr, err := conn.Search(
		ldap.NewSearchRequest(
			dn,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(cn="+id+")",
			[]string{"cn", "uniqueMember"},
			nil,
		),
	)
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

	userDn := "ou=users,ou=" + b.unit + "," + b.baseDN
	userIDs := make([]string, len(members))
	for _, member := range members {
		if strings.HasSuffix(member, userDn) {
			userIDs = append(userIDs, strings.TrimPrefix(strings.TrimSuffix(member, ","+userDn), "cn="))
		}
	}

	return group.NewGroup(
		entry.GetAttributeValue("cn"),
		userIDs...,
	), nil
}

// CreateGroup ...
func (b BackendLDAP) CreateGroup(g group.Group) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + g.ID() + ",ou=groups,ou=" + b.unit + "," + b.baseDN

	members := []string{}
	for _, member := range g.Members() {
		members = append(members, "cn="+member+",ou=users,ou="+b.unit+","+b.baseDN)
	}

	addRequest := ldap.NewAddRequest(dn, []ldap.Control{})
	addRequest.Attribute("objectClass", []string{"groupOfUniqueNames"})
	addRequest.Attribute("cn", []string{g.ID()})
	addRequest.Attribute("uniqueMember", members)

	err = conn.Add(addRequest)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroup ...
func (b BackendLDAP) DeleteGroup(id string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + id + ",ou=groups,ou=" + b.unit + "," + b.baseDN

	delRequest := ldap.NewDelRequest(dn, []ldap.Control{})

	err = conn.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}

// AddToGroup ...
func (b BackendLDAP) AddToGroup(id string, memberIDs ...string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + id + ",ou=groups,ou=" + b.unit + "," + b.baseDN

	members := []string{}
	for _, member := range memberIDs {
		members = append(members, "cn="+member+",ou=users,ou="+b.unit+","+b.baseDN)
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Add("uniqueMember", members)

	err = conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFromGroup ...
func (b BackendLDAP) RemoveFromGroup(id string, memberIDs ...string) error {
	if b.creds == nil {
		return fmt.Errorf("no credentials")
	}

	if strings.TrimSpace(b.unit) == "" {
		return fmt.Errorf("no unit selected")
	}

	conn, err := b.dialToServer(b.creds)
	if err != nil {
		return err
	}

	defer conn.Close()

	if !b.authenticateToServer(b.creds, conn) {
		return fmt.Errorf("invalid credentials")
	}

	dn := "cn=" + id + ",ou=groups,ou=" + b.unit + "," + b.baseDN

	members := []string{}
	for _, member := range memberIDs {
		members = append(members, "cn="+member+",ou=users,ou="+b.unit+","+b.baseDN)
	}

	modRequest := ldap.NewModifyRequest(dn, []ldap.Control{})
	modRequest.Delete("uniqueMember", members)

	err = conn.Modify(modRequest)
	if err != nil {
		return err
	}

	return nil
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
	err := conn.Bind(c.Username(), c.Password())
	if err != nil {
		return false
	}

	return true
}
