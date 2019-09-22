package policies

// A Policy configure how data is queried in LDAP
type Policy struct {
	Attributes []string
}

// DefaultPolicy is used when no policy is configured
var DefaultPolicy = Policy{
	Attributes: []string{"cn", "description"},
}
