package unit

// Backend interface.
type Backend interface {
	ListUnits() (List, error)
	GetUnit(id string) (Unit, error)
	CreateUnit(Unit) error
	UpdateUnit(Unit) error
	AppendUnit(Unit) error
	RemoveUnit(Unit) error
	DeleteUnit(id string) error
	UseUnit(id string)
}
