package standard

import "errors"

type source int

const (
	Environment source = iota + 1
	File
	Parameter
)

var errorUnknownSource = errors.New("unknown source")

func (src source) MarshalYAML() (interface{}, error) {
	switch src {
	case Environment:
		return "ENVIRONMENT_VARIABLE", nil
	case File:
		return "FILE", nil
	case Parameter:
		return "CONFIGURATION_PARAMETER", nil
	default:
		return nil, errorUnknownSource
	}
}

func (src *source) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var literal string
	if err := unmarshal(&literal); err != nil {
		return err
	}

	switch literal {
	case "ENVIRONMENT_VARIABLE":
		*src = Environment
	case "FILE":
		*src = File
	case "CONFIGURATION_PARAMETER":
		*src = Parameter
	default:
		return errorUnknownSource
	}
	return nil
}
