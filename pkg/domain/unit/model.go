package unit

import (
	"encoding/json"
	"fmt"
)

// Unit ...
type Unit interface {
	ID() string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type unit struct {
	id string
}

// NewUnit ...
func NewUnit(id string) Unit {
	return unit{
		id: id,
	}
}

func (u unit) ID() string     { return u.id }
func (u unit) String() string { return fmt.Sprintf("%v", u.ID()) }
func (u unit) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct{ ID string }{u.id})
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (u unit) MarshalYAML() (interface{}, error) {
	return struct{ ID string }{u.id}, nil
}

// List ...
type List interface {
	All() []Unit
	Index(idx uint) Unit
	Len() uint
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type unitlist struct {
	slice []Unit
}

// NewList ...
func NewList(slice []Unit) List {
	return unitlist{slice}
}

func (l unitlist) All() []Unit         { return l.slice }
func (l unitlist) Index(idx uint) Unit { return l.slice[idx] }
func (l unitlist) Len() uint           { return uint(len(l.slice)) }
func (l unitlist) String() string      { return fmt.Sprintf("%v unit(s)", l.Len()) }
func (l unitlist) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(l.slice)
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (l unitlist) MarshalYAML() (interface{}, error) {
	return l.slice, nil
}
