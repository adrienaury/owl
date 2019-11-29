package unit

// Backend ...
type Backend interface {
	ListUnits() (List, error)
}
