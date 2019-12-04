package user

import (
	"encoding/json"
	"fmt"
)

// User object contains informations about a user.
type User interface {
	ID() string
	FirstNames() []string
	LastNames() []string
	Emails() []string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type user struct {
	id         string
	firstNames []string
	lastNames  []string
	emails     []string
}

// NewUser create a new user object.
func NewUser(id string, firstNames []string, lastNames []string, emails []string) User {
	return user{
		id:         id,
		firstNames: firstNames,
		lastNames:  lastNames,
		emails:     emails,
	}
}

func (u user) ID() string           { return u.id }
func (u user) FirstNames() []string { return u.firstNames }
func (u user) LastNames() []string  { return u.lastNames }
func (u user) Emails() []string     { return u.emails }
func (u user) String() string {
	return fmt.Sprintf("%v %v %v %v", u.ID(), u.FirstNames(), u.LastNames(), u.Emails())
}
func (u user) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct {
		ID         string
		FirstNames []string
		LastNames  []string
		Emails     []string
	}{
		u.id,
		u.firstNames,
		u.lastNames,
		u.emails,
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
	}{
		u.id,
		u.firstNames,
		u.lastNames,
		u.emails,
	}, nil
}

// List of user objects.
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

// NewList create a new list of users object.
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
