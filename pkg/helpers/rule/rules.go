package rule

// Rule to apply to an attribute
type Rule interface {
	Apply(value string) string
	Validate(value string) bool
}

var rules = map[string]Rule{
	"Capitalize": Capitalize{},
	"HashSha256": HashSha256{},
}
