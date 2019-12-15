package user

import (
	"fmt"
)

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
	user, err := d.backend.GetUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
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
func (d Driver) Apply(u User) (updated bool, err error) {
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

// Update the user with id.
func (d Driver) Update(u User) error {
	user, err := d.backend.GetUser(u.ID())
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user %v doesn't exist", u.ID())
	}

	firstNames := u.FirstNames()
	if len(firstNames) == 0 {
		firstNames = user.FirstNames()
	}

	lastNames := u.LastNames()
	if len(lastNames) == 0 {
		lastNames = user.LastNames()
	}

	emails := u.Emails()
	if len(emails) == 0 {
		emails = user.Emails()
	}

	merged := NewUser(u.ID(), firstNames, lastNames, emails)

	err = d.backend.UpdateUser(merged)
	if err != nil {
		return err
	}

	return nil
}

// Upsert the user with id (update or create).
func (d Driver) Upsert(u User) (updated bool, err error) {
	user, err := d.backend.GetUser(u.ID())
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, d.Create(u)
	}

	firstNames := u.FirstNames()
	if len(firstNames) == 0 {
		firstNames = user.FirstNames()
	}

	lastNames := u.LastNames()
	if len(lastNames) == 0 {
		lastNames = user.LastNames()
	}

	emails := u.Emails()
	if len(emails) == 0 {
		emails = user.Emails()
	}

	merged := NewUser(u.ID(), firstNames, lastNames, emails)

	err = d.backend.UpdateUser(merged)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Append attributes to the user with id.
func (d Driver) Append(u User) error {
	err := d.backend.AppendUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Remove attributes from the user with id.
func (d Driver) Remove(u User) error {
	err := d.backend.RemoveUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Delete the user with id.
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
