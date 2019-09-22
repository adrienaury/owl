package ldap

import (
	"log"

	ldap "gopkg.in/ldap.v3"
)

// Connection manage the LDAP connection
type Connection struct {
	conn *ldap.Conn
	sess *Session
}

// NewConnection todo
func NewConnection(s *Session) *Connection {
	l, err := ldap.DialURL(s.URL)
	if err != nil {
		log.Fatal(err)
	}
	err = c.conn.Bind(s.Username, pw)
	if err != nil {
		log.Fatal(err)
	}
	return &Connection{l, s}
}

// Close todo
func (c *Connection) Close() {
	c.conn.Close()
}
