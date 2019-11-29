package user

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
	list, err := d.backend.ListUsers()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Create ...
func (d Driver) Create(u User) error {
	err := d.backend.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
