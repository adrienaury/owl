package realm

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	storage Storage
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(storage Storage) Driver {
	return Driver{storage}
}

// Set method create (if not exists) or update (if exists) the realm with id.
func (d Driver) Set(id string, url string, username string) error {
	if err := d.storage.CreateOrUpdateRealm(NewRealm(id, url, username)); err != nil {
		return err
	}
	return nil
}

// Get the realm with id.
func (d Driver) Get(id string) (Realm, error) {
	r, err := d.storage.GetRealm(id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete the realm with id.
func (d Driver) Delete(id string) error {
	err := d.storage.DeleteRealm(id)
	if err != nil {
		return err
	}
	return nil
}

// List all realms.
func (d Driver) List() (List, error) {
	list, err := d.storage.ListRealms()
	if err != nil {
		return nil, err
	}
	return list, nil
}
