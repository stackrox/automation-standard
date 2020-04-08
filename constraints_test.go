package standard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraints(t *testing.T) {
	tests := []struct {
		constraint Constraint
		given      string
		error      bool
	}{
		{
			constraint: ConstraintBool(),
			given:      "true",
		},
		{
			constraint: ConstraintBool(),
			given:      "false",
		},
		{
			constraint: ConstraintBool(),
			given:      "yes",
			error:      true,
		},
		{
			constraint: ConstraintBool(),
			given:      "",
			error:      true,
		},

		{
			constraint: ConstraintEnum("foo", "bar"),
			given:      "foo",
		},
		{
			constraint: ConstraintEnum("foo", "bar"),
			given:      "bar",
		},
		{
			constraint: ConstraintEnum("foo", "bar"),
			given:      "baz",
			error:      true,
		},
		{
			constraint: ConstraintEnum("foo"),
			given:      "foo",
			error:      true,
		},

		{
			constraint: ConstraintIntMinimum(2),
			given:      "2",
		},
		{
			constraint: ConstraintIntMinimum(2),
			given:      "3",
		},
		{
			constraint: ConstraintIntMinimum(2),
			given:      "1",
			error:      true,
		},
		{
			constraint: ConstraintIntMinimum(2),
			given:      "two",
			error:      true,
		},

		{
			constraint: ConstraintIntMaximum(5),
			given:      "5",
		},
		{
			constraint: ConstraintIntMaximum(5),
			given:      "4",
		},
		{
			constraint: ConstraintIntMaximum(5),
			given:      "6",
			error:      true,
		},
		{
			constraint: ConstraintIntMaximum(5),
			given:      "five",
			error:      true,
		},

		{
			constraint: ConstraintDockerImage(),
			given:      "docker.io/nginx:1.2.3",
		},
		{
			constraint: ConstraintDockerImage(),
			given:      "stackrox.io/main:1.2.3",
		},
		{
			constraint: ConstraintDockerImage(),
			given:      "docker.io/nginx",
			error:      true,
		},
		{
			constraint: ConstraintDockerImage(),
			given:      "nginx:1.2.3",
			error:      true,
		},
		{
			constraint: ConstraintDockerImage(),
			given:      "nginx",
			error:      true,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("%d", index+1)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.constraint.Check(test.given)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
