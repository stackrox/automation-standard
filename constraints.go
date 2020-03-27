package standard

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type Constraint struct {
	Name        string `json:"name"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description"`
}

var (
	constraintNamePass        = "pass"
	constraintNameFail        = "fail"
	constraintNameIntMinimum  = "int-minimum"
	constraintNameIntMaximum  = "int-maximum"
	constraintNameDockerImage = "docker-image"

	constraintFuncs = map[string]struct {
		description string
		fn          func(value string, given string) error
	}{
		constraintNamePass: {
			description: "always passes",
			fn: func(string, string) error {
				return nil
			},
		},
		constraintNameFail: {
			description: "always fails",
			fn: func(string, string) error {
				return fmt.Errorf("failed")
			},
		},
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
				if matched, _ := regexp.MatchString(`^[a-z0-9-]+(\.[a-z0-9-]+)+/[a-z0-9_.-]+(/[a-z0-9_.-]+)*:[a-z0-9_.-]{1,128}$`, given); !matched {
					return fmt.Errorf("input %s was not a docker image", given)
				}
				return nil
			},
		},
	}
)

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

func ConstraintIntMinimum(minumum int) Constraint {
	c, _ := ConstraintByName(constraintNameIntMinimum, fmt.Sprint(minumum))
	return c
}

func ConstraintIntMaximum(maximum int) Constraint {
	c, _ := ConstraintByName(constraintNameIntMaximum, fmt.Sprint(maximum))
	return c
}

func ConstraintDockerImage() Constraint {
	c, _ := ConstraintByName(constraintNameDockerImage, "")
	return c
}

func (c Constraint) Check(given string) error {
	if pair, found := constraintFuncs[c.Name]; found {
		return errors.Wrapf(pair.fn(c.Value, given), c.Name)
	}

	return fmt.Errorf("unknown constraint %q", c.Name)
}

func CheckAll(cs []Constraint, given string) error {
	for _, c := range cs {
		if err := c.Check(given); err != nil {
			return err
		}
	}

	return nil
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
