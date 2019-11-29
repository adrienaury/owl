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

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.backend.DeleteUnit(id)
	if err != nil {
		return err
	}
	return nil
}
