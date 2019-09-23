package rule

import "strings"

// Capitalize todo
type Capitalize struct{}

// Apply todo
func (Capitalize) Apply(value string) string {
	return strings.Title(value)
}

// Validate todo
func (Capitalize) Validate(value string) bool {
	return strings.Title(value) == value
}
