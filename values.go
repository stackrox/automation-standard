package standard

import "fmt"

// Values represents a set of named parameter values. Isolated to more
// restrictive access.
type Values struct {
	values map[string]string
}

// Get returns the named parameter value, or panics of that parameter does not
// exist. This is intended to guard against silent bugs where a non-existent
// parameter is referenced, and an empty string is consumed.
func (g Values) Get(name string) string {
	if value, found := g.values[name]; found {
		return value
	}
	panic(fmt.Sprintf("unknown parameter %q", name))
}
