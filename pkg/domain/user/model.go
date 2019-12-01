package user

import (
	"encoding/json"
	"fmt"
)

// User ...
type User interface {
	ID() string
	FirstNames() []string
	LastNames() []string
	Emails() []string
	Groups() []string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type user struct {
	id         string
	firstNames []string
	lastNames  []string
	emails     []string
	groups     []string
}

// NewUser ...
func NewUser(id string, firstNames []string, lastNames []string, emails []string, groups []string) User {
	return user{
		id:         id,
		firstNames: firstNames,
		lastNames:  lastNames,
		emails:     emails,
		groups:     groups,
	}
}

func (u user) ID() string           { return u.id }
func (u user) FirstNames() []string { return u.firstNames }
func (u user) LastNames() []string  { return u.lastNames }
func (u user) Emails() []string     { return u.emails }
func (u user) Groups() []string     { return u.groups }
func (u user) String() string {
	return fmt.Sprintf("%v %v %v %v %v", u.ID(), u.FirstNames(), u.LastNames(), u.Emails(), u.Groups())
}
func (u user) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct {
		ID         string
		FirstNames []string
		LastNames  []string
		Emails     []string
		Groups     []string
	}{
		u.id,
		u.firstNames,
		u.lastNames,
		u.emails,
		u.groups,
	})
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (u user) MarshalYAML() (interface{}, error) {
	return struct {
		ID         string
		FirstNames []string
		LastNames  []string
		Emails     []string
		Groups     []string
	}{
		u.id,
		u.firstNames,
		u.lastNames,
		u.emails,
		u.groups,
	}, nil
}

// List ...
type List interface {
	All() []User
	Index(idx uint) User
	Len() uint
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type userlist struct {
	slice []User
}

// NewList ...
func NewList(slice []User) List {
	return userlist{slice}
}

func (l userlist) All() []User         { return l.slice }
func (l userlist) Index(idx uint) User { return l.slice[idx] }
func (l userlist) Len() uint           { return uint(len(l.slice)) }
func (l userlist) String() string      { return fmt.Sprintf("%v users(s)", l.Len()) }
func (l userlist) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(l.slice)
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (l userlist) MarshalYAML() (interface{}, error) {
	return l.slice, nil
}
