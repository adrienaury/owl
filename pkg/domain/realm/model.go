package realm

import (
	"encoding/json"
	"fmt"
)

// Realm object contains informations about a realm.
type Realm interface {
	ID() string
	URL() string
	Username() string
	Policy() string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type realm struct {
	id       string
	url      string
	username string
	policy   string
}

// NewRealm create a new realm object.
func NewRealm(id, url, username string) Realm {
	return realm{
		id:       id,
		url:      url,
		username: username,
		policy:   "default",
	}
}

// NewRealmWithPolicy create a new realm object.
func NewRealmWithPolicy(id, url, username, policy string) Realm {
	return realm{
		id:       id,
		url:      url,
		username: username,
		policy:   policy,
	}
}

func (r realm) ID() string       { return r.id }
func (r realm) URL() string      { return r.url }
func (r realm) Username() string { return r.username }
func (r realm) Policy() string   { return r.policy }
func (r realm) String() string   { return fmt.Sprintf("%v %v %v", r.ID(), r.URL(), r.Username()) }
func (r realm) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct{ ID, URL, Username string }{r.id, r.url, r.username})
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (r realm) MarshalYAML() (interface{}, error) {
	return struct{ ID, URL, Username string }{r.id, r.url, r.username}, nil
}

// List of realm objects.
type List interface {
	All() []Realm
	Index(idx uint) Realm
	Len() uint
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type realmlist struct {
	slice []Realm
}

// NewList create a new list of realms object.
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
func (l realmlist) MarshalYAML() (interface{}, error) {
	return l.slice, nil
}
