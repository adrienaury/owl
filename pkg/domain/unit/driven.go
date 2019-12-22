package unit

// Backend interface.
type Backend interface {
	ListUnits(object string, fields map[string]string) (List, error)
	GetUnit(id string, object string, fields map[string]string) (Unit, error)
	CreateUnit(u Unit, object string, fields map[string]string) error
	UpdateUnit(u Unit, object string, fields map[string]string) error
	AppendUnit(u Unit, object string, fields map[string]string) error
	RemoveUnit(u Unit, object string, fields map[string]string) error
	DeleteUnit(id string, object string, fields map[string]string) error
	UseUnit(id string)
}
