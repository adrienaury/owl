package user

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	backend Backend
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(backend Backend) Driver {
	return Driver{backend}
}

// List all users.
func (d Driver) List() (List, error) {
	list, err := d.backend.ListUsers()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get the user with id.
func (d Driver) Get(id string) (User, error) {
	unit, err := d.backend.GetUser(id)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// Create a new user.
func (d Driver) Create(u User) error {
	err := d.backend.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Apply new attributes to an existing user, or create a new one.
func (d Driver) Apply(u User) (bool, error) {
	user, err := d.backend.GetUser(u.ID())
	if err != nil {
		return false, err
	}

	exists := false
	if user != nil {
		exists = true
		err = d.backend.UpdateUser(u)
	} else {
		err = d.backend.CreateUser(u)
	}

	if err != nil {
		return exists, err
	}

	return exists, nil
}

// Delete the user with id.
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
