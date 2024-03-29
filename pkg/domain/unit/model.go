package unit

import (
	"encoding/json"
	"fmt"
)

// Unit object contains information about an organizational unit.
type Unit interface {
	ID() string
	Description() string
	String() string
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
}

type unit struct {
	id         string
	decription string
}

// NewUnit create a new unit object.
func NewUnit(id, decription string) Unit {
	return unit{
		id:         id,
		decription: decription,
	}
}

func (u unit) ID() string          { return u.id }
func (u unit) Description() string { return u.decription }
func (u unit) String() string      { return fmt.Sprintf("%v", u.ID()) }
func (u unit) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(struct {
		ID          string
		Description string
	}{u.id, u.decription})
	if e != nil {
		return nil, e
	}
	return b, nil
}
func (u unit) MarshalYAML() (interface{}, error) {
	return struct {
		ID          string
		Description string
	}{u.id, u.decription}, nil
}

// List of unit objects.
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

// NewList create a new list of units object.
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
