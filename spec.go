package standard

import (
	"fmt"
	"os"
	"sort"
)

type Spec struct {
	Name        string
	Description string
	Source      source
}

func (spec Spec) Resolve(config Config) (string, error) {
	switch spec.Source {
	case Parameter:
		value, found := config[spec.Name]
		if !found {
			return "", fmt.Errorf("the parameter %q was not found", spec.Name)
		}
		return value, nil

	case Environment:
		value, found := os.LookupEnv(spec.Name)
		if !found {
			return "", fmt.Errorf("the environment variable %q was not found", spec.Name)
		}
		return value, nil

	case File:
		if _, err := os.Stat(spec.Name); err != nil {
			return "", fmt.Errorf("the file %q was not found", spec.Name)
		}
		return spec.Name, nil

	default:
		return "", errorUnknownSource
	}
}

func joinSpecs(first, second []Spec) ([]Spec, error) {
	names := make(map[string]struct{}, len(first)+len(second))
	result := make([]Spec, 0, len(first)+len(second))
	for _, spec := range first {
		names[spec.Name] = struct{}{}
		result = append(result, spec)
	}

	for _, spec := range second {
		if _, found := names[spec.Name]; found {
			return nil, fmt.Errorf("duplicate spec name %q", spec.Name)
		}
		names[spec.Name] = struct{}{}
		result = append(result, spec)
	}

	return result, nil
}

type result struct {
	name    string
	message string
	err     error
}

func validateSpecs(specs []Spec, config Config) []result {
	results := make([]result, 0, len(specs))
	paramNames := make(map[string]struct{})

	// Filter a set of spec names, where the source is of the parameter type.
	for _, spec := range specs {
		if spec.Source != Parameter {
			continue
		}
		paramNames[spec.Name] = struct{}{}
	}

	// Examine the given config file, and verify that all of the parameters
	// have a matching spec. This avoids having extraneous values in the config file.
	for name := range config {
		if _, found := paramNames[name]; !found {
			results = append(results, result{
				name:    name,
				message: name,
				err:     fmt.Errorf("the parameter %q did not have a matching spec", name),
			})
		}
	}

	// Resolve each spec, and report an error if the file/env/parameter doesn't exist.
	for _, spec := range specs {
		_, err := spec.Resolve(config)
		results = append(results, result{
			name:    spec.Name,
			message: fmt.Sprintf("%s (%s)", spec.Name, spec.Description),
			err:     err,
		})
	}

	return results
}

func sortSpecs(specs []Spec) {
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Name < specs[j].Name
	})
}
