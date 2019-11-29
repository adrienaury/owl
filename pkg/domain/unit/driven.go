package unit

// Backend ...
type Backend interface {
	ListUnits() (List, error)
	CreateUnit(Unit) error
	DeleteUnit(id string) error
}
