package rule

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jsimonetti/pwscheme/ssha"
)

// SchemeSSHA todo
type SchemeSSHA struct{}

// Apply todo
func (SchemeSSHA) Apply(value string) string {
	scheme, err := ssha.Generate(value, 8)
	if err != nil {
		fmt.Println(err)
		// TODO
	}
	return scheme
}

// Validate todo
func (SchemeSSHA) Validate(value string) bool {
	_, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "{SSHA}"))
	return strings.HasPrefix(value, "{SSHA}") && err == nil
}
