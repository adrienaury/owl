package unit

// Backend ...
type Backend interface {
	ListUnits() (List, error)
	GetUnit(id string) (Unit, error)
	CreateUnit(Unit) error
	DeleteUnit(id string) error
	UseUnit(id string)
}
