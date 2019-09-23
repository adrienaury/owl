package rule

import "fmt"

// Rule to apply to an attribute
type Rule interface {
	Apply(value string) string
	Validate(value string) bool
}

var rules = map[string]Rule{
	"Capitalize": Capitalize{},
	"SchemeSSHA": SchemeSSHA{},
}

// Apply todo
func Apply(rulename, value string) (string, error) {
	rule, ok := rules[rulename]
	if ok {
		return rule.Apply(value), nil
	}
	return "", fmt.Errorf("todo")
}

// Validate todo
func Validate(rulename, value string) (bool, error) {
	rule, ok := rules[rulename]
	if ok {
		return rule.Validate(value), nil
	}
	return false, fmt.Errorf("todo")
}
