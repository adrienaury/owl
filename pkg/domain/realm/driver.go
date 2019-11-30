package realm

// Driver ...
type Driver struct {
	storage Storage
}

// NewDriver ...
func NewDriver(storage Storage) Driver {
	return Driver{storage}
}

// Set ...
func (d Driver) Set(id string, url string, username string) error {
	if err := d.storage.CreateOrUpdateRealm(NewRealm(id, url, username)); err != nil {
		return err
	}
	return nil
}

// Get ...
func (d Driver) Get(id string) (Realm, error) {
	r, err := d.storage.GetRealm(id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (d Driver) Delete(id string) error {
	err := d.storage.DeleteRealm(id)
	if err != nil {
		return err
	}
	return nil
}

// List ...
func (d Driver) List() (List, error) {
	list, err := d.storage.ListRealms()
	if err != nil {
		return nil, err
	}
	return list, nil
}
