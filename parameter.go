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
func (spec Parameter) Resolve(cmd *cobra.Command) (string, error) {
	var resolved string
	switch spec.Source {
	case Flag:
		value, err := cmd.Flags().GetString(spec.Name)
		if err != nil {
			return "", fmt.Errorf("the parameter %q was not found", spec.Name)
		}
		resolved = value

	case Environment:
		value, found := os.LookupEnv(spec.Name)
		if !found {
			return "", fmt.Errorf("the environment variable %q was not found", spec.Name)
		}
		resolved = value

	case File:
		if _, err := os.Stat(spec.Name); err != nil {
			return "", fmt.Errorf("the file %q was not found", spec.Name)
		}
		resolved = spec.Name

	default:
		return "", errorUnknownSource
	}

	if err := CheckAll(spec.Constraints, resolved); err != nil {
		return "", err
	}

	return resolved, nil
}

type result struct {
	name    string
	message string
	err     error
}

func validateSpecs(specs []Parameter, cmd *cobra.Command) []result {
	results := make([]result, 0, len(specs))
	paramNames := make(map[string]struct{})

	// Filter a set of spec names, where the source is of the parameter type.
	for _, spec := range specs {
		if spec.Source != Flag {
			continue
		}
		paramNames[spec.Name] = struct{}{}
	}

	// Examine the given config file, and verify that all of the parameters
	// have a matching spec. This avoids having extraneous values in the config file.
	//for name := range config {
	//	if _, found := paramNames[name]; !found {
	//		results = append(results, result{
	//			name:    name,
	//			message: name,
	//			err:     fmt.Errorf("the parameter %q did not have a matching spec", name),
	//		})
	//	}
	//}

	// Resolve each spec, and report an error if the file/env/parameter doesn't exist.
	for _, spec := range specs {
		_, err := spec.Resolve(cmd)
		results = append(results, result{
			name:    spec.Name,
			message: fmt.Sprintf("%s (%s)", spec.Name, spec.Description),
			err:     err,
		})
	}

	return results
}

func sortSpecs(specs []Parameter) {
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Name < specs[j].Name
	})
}
