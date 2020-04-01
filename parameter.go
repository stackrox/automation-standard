package standard

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// Parameter represents a single parameter that is consumed by an application.
type Parameter struct {
	// Name is the name for this parameter.
	//
	// Example: "node-count"
	Name string `json:"name"`

	// Description is a human-readable Description for this parameter.
	//
	// Example: "number of nodes to provision"
	Description string `json:"description"`

	// Source is how this parameter is configured, as a flag, environment
	// variable, or file.
	Source source `json:"source"`

	// Constraints is a list of constraints used vor validating parameter
	// values.
	Constraints []Constraint `json:"constraints,omitempty"`
}

// Resolve obtains the value of the given parameter from the appropriate source.
func (spec Parameter) Resolve(cmd *cobra.Command) (string, []error) {
	var resolved string
	switch spec.Source {
	case Flag:
		value, err := cmd.Flags().GetString(spec.Name)
		if err != nil {
			return "", []error{fmt.Errorf("the parameter %q was not found", spec.Name)}
		}
		resolved = value

	case Environment:
		value, found := os.LookupEnv(spec.Name)
		if !found {
			return "", []error{fmt.Errorf("the environment variable %q was not found", spec.Name)}
		}
		resolved = value

	case File:
		if _, err := os.Stat(spec.Name); err != nil {
			return "", []error{fmt.Errorf("the file %q was not found", spec.Name)}
		}
		resolved = spec.Name

	default:
		return "", []error{errorUnknownSource}
	}

	if errs := checkAll(spec.Constraints, resolved); len(errs) != 0 {
		return "", errs
	}

	return resolved, nil
}

func sortSpecs(specs []Parameter) {
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Name < specs[j].Name
	})
}
