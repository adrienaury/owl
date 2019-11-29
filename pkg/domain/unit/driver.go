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

// Create ...
func (d Driver) Create(u Unit) error {
	err := d.backend.CreateUnit(u)
	if err != nil {
		return err
	}
	return nil
}
