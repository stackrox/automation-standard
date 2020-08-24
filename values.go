package standard

import "fmt"

// Values represents a set of named parameter values. Isolated to more
// restrictive access.
type Values struct {
	values map[string]string
}

// Get returns the named parameter value, or panics if that parameter does not
// exist. This is intended to guard against silent bugs where a non-existent
// parameter is referenced, and an empty string is consumed.
func (g Values) Get(name string) string {
	if value, found := g.values[name]; found {
		return value
	}
	panic(fmt.Sprintf("unknown parameter %q", name))
}

// GetMany returns the named parameter values, or panics if a parameter does not
// exist. This is intended to guard against silent bugs where a non-existent
// parameter is referenced, and an empty string is consumed.
func (g Values) GetMany(names ...string) []string {
	var values []string
	for _, name := range names {
		values = append(values, g.Get(name))
	}
	return values
}
