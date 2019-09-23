package rule

import "strings"

// HashSha256 todo
type HashSha256 struct{}

// Apply todo
func (HashSha256) Apply(value string) string {
	return strings.Title(value)
}

// Validate todo
func (HashSha256) Validate(value string) bool {
	return strings.Title(value) == value
}
