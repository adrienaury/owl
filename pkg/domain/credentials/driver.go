package credentials

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	storage Storage
	backend Backend
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(storage Storage, backend Backend) Driver {
	return Driver{storage, backend}
}

// Set ...
func (d Driver) Set(url string, username string, password string) error {
	if err := d.storage.SetCredentials(NewCredentials(url, username, password)); err != nil {
		return err
	}
	return nil
}

// Remove ...
func (d Driver) Remove(url string, username string) error {
	if err := d.storage.RemoveCredentials(url, username); err != nil {
		return err
	}
	return nil
}

// Get ...
func (d Driver) Get(url string, username string) (Credentials, error) {
	r, err := d.storage.GetCredentials(url, username)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Test ...
func (d Driver) Test(creds Credentials) (bool, error) {
	ok, err := d.backend.TestCredentials(creds)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Use ...
func (d Driver) Use(creds Credentials) {
	d.backend.UseCredentials(creds)
}
