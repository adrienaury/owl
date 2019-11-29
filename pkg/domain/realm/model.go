package realm

import (
	"encoding/json"
	"fmt"
)

// Realm ...
type Realm interface {
	ID() string
	URL() string
	Username() string
	String() string
	MarshalJSON() ([]byte, error)
}

type realm struct {
	id       string
	url      string
	username string
}

// NewRealm ...
func NewRealm(id, url, username string) Realm {
	return realm{
		id:       id,
		url:      url,
		username: username,
	}
}

func (r realm) ID() string       { return r.id }
func (r realm) URL() string      { return r.url }
func (r realm) Username() string { return r.username }
func (r realm) String() string   { return fmt.Sprintf("%v %v %v", r.ID(), r.URL(), r.Username()) }
func (r realm) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct{ ID, URL, Username string }{r.id, r.url, r.username})
	if e != nil {
		return nil, e
	}
	return b, nil
}

// List ...
type List interface {
	All() []Realm
	Index(idx uint) Realm
	Len() uint
	String() string
	MarshalJSON() ([]byte, error)
}

type realmlist struct {
	slice []Realm
}

// NewList ...
func NewList(slice []Realm) List {
	return realmlist{slice}
}

func (l realmlist) All() []Realm         { return l.slice }
func (l realmlist) Index(idx uint) Realm { return l.slice[idx] }
func (l realmlist) Len() uint            { return uint(len(l.slice)) }
func (l realmlist) String() string       { return fmt.Sprintf("%v realm(s)", l.Len()) }
func (l realmlist) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(l.slice)
	if e != nil {
		return nil, e
	}
	return b, nil
}
