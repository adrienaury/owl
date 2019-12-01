package unit

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
	list, err := d.backend.ListUnits()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get ...
func (d Driver) Get(id string) (Unit, error) {
	unit, err := d.backend.GetUnit(id)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// Create ...
func (d Driver) Create(u Unit) error {
	err := d.backend.CreateUnit(u)
	if err != nil {
		return err
	}
	return nil
}

// Apply ...
func (d Driver) Apply(u Unit) (bool, error) {
	user, err := d.backend.GetUnit(u.ID())
	if err != nil {
		return false, err
	}

	exists := false
	if user != nil {
		exists = true
		err = d.backend.UpdateUnit(u)
	} else {
		err = d.backend.CreateUnit(u)
	}

	if err != nil {
		return exists, err
	}

	return exists, nil
}

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUnit(id)
	if err != nil {
		return err
	}
	return nil
}

// Use ...
func (d Driver) Use(id string) {
	d.backend.UseUnit(id)
}
