package policy

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	storage Storage
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(storage Storage) Driver {
	d := Driver{storage}
	return d
}

// Get the policy with name.
func (d Driver) Get(name string) (Policy, error) {
	p, err := d.storage.GetPolicy(name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return defaultPolicy, nil
	}
	return p, nil
}
