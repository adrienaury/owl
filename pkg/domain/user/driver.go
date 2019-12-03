package user

import (
	"fmt"
)

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

// Get ...
func (d Driver) Get(id string) (User, error) {
	unit, err := d.backend.GetUser(id)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// Create ...
func (d Driver) Create(u User) error {
	err := d.backend.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Apply ...
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

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

// AssignPassword ...
func (d Driver) AssignPassword(userID string, password string) error {
	user, err := d.backend.GetUser(userID)
	if err != nil {
		return err
	}

	if len(user.Emails()) <= 0 {
		return fmt.Errorf("user has no e-mail, password change is forbidden")
	}

	if err := d.backend.SetUserPassword(userID, password); err != nil {
		return err
	}

	return nil
}
