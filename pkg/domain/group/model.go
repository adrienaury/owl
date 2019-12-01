package group

import (
	"encoding/json"
	"fmt"
)

// Group ...
type Group interface {
	ID() string
	Members() []string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type group struct {
	id      string
	members []string
}

// NewGroup ...
func NewGroup(id string, members ...string) Group {
	return group{
		id:      id,
		members: members,
	}
}

func (g group) ID() string        { return g.id }
func (g group) Members() []string { return g.members }
func (g group) String() string    { return fmt.Sprintf("%v", g.ID()) }
func (g group) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct {
		ID      string
		Members []string
	}{g.id, g.members})
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (g group) MarshalYAML() (interface{}, error) {
	return struct {
		ID      string
		Members []string
	}{g.id, g.members}, nil
}

// List ...
type List interface {
	All() []Group
	Index(idx uint) Group
	Len() uint
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type grouplist struct {
	slice []Group
}

// NewList ...
func NewList(slice []Group) List {
	return grouplist{slice}
}

func (l grouplist) All() []Group         { return l.slice }
func (l grouplist) Index(idx uint) Group { return l.slice[idx] }
func (l grouplist) Len() uint            { return uint(len(l.slice)) }
func (l grouplist) String() string       { return fmt.Sprintf("%v group(s)", l.Len()) }
func (l grouplist) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(l.slice)
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (l grouplist) MarshalYAML() (interface{}, error) {
	return l.slice, nil
}
