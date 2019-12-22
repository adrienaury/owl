package unit

import (
	"fmt"
	"strings"
)

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	backend Backend
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(backend Backend) Driver {
	return Driver{backend}
}

// List all units.
func (d Driver) List(p Policy) (List, error) {
	list, err := d.backend.ListUnits(p.BackendObject(), p.BackendFields())
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get the unit with id.
func (d Driver) Get(id string, p Policy) (Unit, error) {
	unit, err := d.backend.GetUnit(id, p.BackendObject(), p.BackendFields())
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// Create a new unit.
func (d Driver) Create(u Unit, p Policy) error {
	err := d.backend.CreateUnit(u, p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}
	return nil
}

// Apply new attributes to an existing unit, or create a new one.
func (d Driver) Apply(u Unit, p Policy) (updated bool, err error) {
	unit, err := d.backend.GetUnit(u.ID(), p.BackendObject(), p.BackendFields())
	if err != nil {
		return false, err
	}

	exists := false
	if unit != nil {
		exists = true
		err = d.backend.UpdateUnit(u, p.BackendObject(), p.BackendFields())
	} else {
		err = d.backend.CreateUnit(u, p.BackendObject(), p.BackendFields())
	}

	if err != nil {
		return exists, err
	}

	return exists, nil
}

// Update the unit with id.
func (d Driver) Update(u Unit, p Policy) error {
	unit, err := d.backend.GetUnit(u.ID(), p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}

	if unit == nil {
		return fmt.Errorf("unit %v doesn't exist", u.ID())
	}

	description := u.Description()
	if strings.TrimSpace(description) == "" {
		description = unit.Description()
	}
	merged := NewUnit(u.ID(), description)

	err = d.backend.UpdateUnit(merged, p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}

	return nil
}

// Upsert the unit with id (update or create).
func (d Driver) Upsert(u Unit, p Policy) (updated bool, err error) {
	unit, err := d.backend.GetUnit(u.ID(), p.BackendObject(), p.BackendFields())
	if err != nil {
		return false, err
	}

	if unit == nil {
		return false, d.Create(u, p)
	}

	description := u.Description()
	if strings.TrimSpace(description) == "" {
		description = unit.Description()
	}
	merged := NewUnit(u.ID(), description)

	err = d.backend.UpdateUnit(merged, p.BackendObject(), p.BackendFields())
	if err != nil {
		return false, err
	}

	return true, nil
}

// Append attributes to the unit with id.
func (d Driver) Append(u Unit, p Policy) error {
	err := d.backend.AppendUnit(u, p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}
	return nil
}

// Remove attributes from the unit with id.
func (d Driver) Remove(u Unit, p Policy) error {
	err := d.backend.RemoveUnit(u, p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}
	return nil
}

// Delete the unit with id.
func (d Driver) Delete(id string, p Policy) error {
	err := d.backend.DeleteUnit(id, p.BackendObject(), p.BackendFields())
	if err != nil {
		return err
	}
	return nil
}

// Use ask the backend to start using this unit as default unit.
func (d Driver) Use(id string) {
	d.backend.UseUnit(id)
}
