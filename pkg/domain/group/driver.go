package group

// Driver ...
type Driver struct {
	backend Backend
}

// NewDriver ...
func NewDriver(backend Backend) Driver {
	return Driver{backend}
}

// List ...
func (d Driver) List() (List, error) {
	list, err := d.backend.ListGroups()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get ...
func (d Driver) Get(id string) (Group, error) {
	group, err := d.backend.GetGroup(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Create ...
func (d Driver) Create(u Group) error {
	err := d.backend.CreateGroup(u)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteGroup(id)
	if err != nil {
		return err
	}
	return nil
}

// AddMembers ...
func (d Driver) AddMembers(id string, memberIDs ...string) error {
	err := d.backend.AddToGroup(id, memberIDs...)
	if err != nil {
		return err
	}
	return nil
}

// RemoveMembers ...
func (d Driver) RemoveMembers(id string, memberIDs ...string) error {
	err := d.backend.RemoveFromGroup(id, memberIDs...)
	if err != nil {
		return err
	}
	return nil
}
