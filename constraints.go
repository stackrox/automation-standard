package standard

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// Constraint represents a restriction that can be placed on a parameter.
type Constraint struct {
	Name        string `json:"name"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description"`
}

const (
	constraintNameIntMinimum  = "int-minimum"
	constraintNameIntMaximum  = "int-maximum"
	constraintNameDockerImage = "docker-image"

	regexDockerImage = `^[a-z0-9-]+(\.[a-z0-9-]+)+/[a-z0-9_.-]+(/[a-z0-9_.-]+)*:[a-z0-9_.-]{1,128}$`
)

var (
	constraintFuncs = map[string]struct { // nolint:gochecknoglobals
		description string
		fn          func(value string, given string) error
	}{
		constraintNameIntMinimum: {
			description: "ensure value minimum",
			fn: func(value string, given string) error {
				minimumInt, givenInt, err := asInts(value, given)
				if err != nil {
					return err
				}
				if givenInt < minimumInt {
					return fmt.Errorf("input %d was less than %d", givenInt, minimumInt)
				}
				return nil
			},
		},
		constraintNameIntMaximum: {
			description: "ensure value maximum",
			fn: func(value string, given string) error {
				maximum, givenInt, err := asInts(value, given)
				if err != nil {
					return err
				}
				if maximum < givenInt {
					return fmt.Errorf("input %d was greater than %d", givenInt, maximum)
				}
				return nil
			},
		},
		constraintNameDockerImage: {
			description: "ensure value is a docker image name",
			fn: func(_ string, given string) error {
				if matched, _ := regexp.MatchString(regexDockerImage, given); !matched {
					return fmt.Errorf("input %s was not a docker image", given)
				}
				return nil
			},
		},
	}
)

// ConstraintByName returns the named constraint if it exists, and an error
// otherwise.
func ConstraintByName(name string, value string) (Constraint, error) {
	if pair, found := constraintFuncs[name]; found {
		return Constraint{
			Name:        name,
			Value:       value,
			Description: pair.description,
		}, nil
	}

	return Constraint{}, fmt.Errorf("unknown constraint %q", name)
}

// mustConstraintByName returns the named constraint if it exists, and panics
// otherwise. A panic indicates a bug in this library.
func mustConstraintByName(name string, value string) Constraint {
	constraint, err := ConstraintByName(name, value)
	if err != nil {
		panic(errors.Wrapf(err, "failed to lookup built-in constraint"))
	}

	return constraint
}

// ConstraintIntMinimum creates a constraint for ensuring a parameter has a
// minimum integer value.
func ConstraintIntMinimum(minimum int) Constraint {
	c := mustConstraintByName(constraintNameIntMinimum, fmt.Sprint(minimum))
	return c
}

// ConstraintIntMaximum creates a constraint for ensuring a parameter has a
// maximum integer value.
func ConstraintIntMaximum(maximum int) Constraint {
	c := mustConstraintByName(constraintNameIntMaximum, fmt.Sprint(maximum))
	return c
}

// ConstraintDockerImage creates a constraint for ensuring a parameter is a
// valid Docker image name.
func ConstraintDockerImage() Constraint {
	c := mustConstraintByName(constraintNameDockerImage, "")
	return c
}

// Check validates that the given value conforms to this constraint.
func (c Constraint) Check(given string) error {
	if pair, found := constraintFuncs[c.Name]; found {
		return errors.Wrapf(pair.fn(c.Value, given), c.Name)
	}

	return fmt.Errorf("unknown constraint %q", c.Name)
}

// checkAll validates that the given value conforms to all given constraints.
func checkAll(cs []Constraint, given string) []error {
	var errs []error
	for _, c := range cs {
		if err := c.Check(given); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func asInts(value string, given string) (int, int, error) {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, errors.New("constraint value was not an integer")
	}

	givenInt, err := strconv.Atoi(given)
	if err != nil {
		return 0, 0, errors.New("given value was not an integer")
	}

	return valueInt, givenInt, nil
}
