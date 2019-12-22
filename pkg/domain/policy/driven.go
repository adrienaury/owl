package policy

// Storage interface for reading policies.
type Storage interface {
	GetPolicy(name string) (Policy, error)
}
