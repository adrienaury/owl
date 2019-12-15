package group

import "fmt"

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	backend Backend
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(backend Backend) Driver {
	return Driver{backend}
}

// List all groups.
func (d Driver) List() (List, error) {
	list, err := d.backend.ListGroups()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get the group with id.
func (d Driver) Get(id string) (Group, error) {
	group, err := d.backend.GetGroup(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Create a new group.
func (d Driver) Create(g Group) error {
	err := d.backend.CreateGroup(g)
	if err != nil {
		return err
	}
	return nil
}

// Apply new attributes to an existing group, or create a new one.
func (d Driver) Apply(g Group) (updated bool, err error) {
	group, err := d.backend.GetGroup(g.ID())
	if err != nil {
		return false, err
	}

	exists := false
	if group != nil {
		exists = true
		err = d.backend.UpdateGroup(g)
	} else {
		err = d.backend.CreateGroup(g)
	}

	if err != nil {
		return exists, err
	}

	return exists, nil
}

// Update the group with id.
func (d Driver) Update(g Group) error {
	group, err := d.backend.GetGroup(g.ID())
	if err != nil {
		return err
	}

	if group == nil {
		return fmt.Errorf("group %v doesn't exist", g.ID())
	}

	members := g.Members()
	if len(members) == 0 {
		members = group.Members()
	}

	merged := NewGroup(g.ID(), members...)

	err = d.backend.UpdateGroup(merged)
	if err != nil {
		return err
	}

	return nil
}

// Upsert the group with id (update or create).
func (d Driver) Upsert(g Group) (updated bool, err error) {
	group, err := d.backend.GetGroup(g.ID())
	if err != nil {
		return false, err
	}

	if group == nil {
		return false, d.Create(g)
	}

	members := g.Members()
	if len(members) == 0 {
		members = group.Members()
	}

	merged := NewGroup(g.ID(), members...)

	err = d.backend.UpdateGroup(merged)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete the group with id.
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteGroup(id)
	if err != nil {
		return err
	}
	return nil
}

// Append attributes (members) to the group with id.
func (d Driver) Append(g Group) error {
	err := d.backend.AppendGroup(g)
	if err != nil {
		return err
	}
	return nil
}

// Remove attributes (members) from the group with id.
func (d Driver) Remove(g Group) error {
	err := d.backend.RemoveGroup(g)
	if err != nil {
		return err
	}
	return nil
}
