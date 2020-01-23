package standard

import (
	"encoding/json"
	"errors"
)

var errorUnknownSource = errors.New("unknown source")

type source int

const (
	Environment source = iota + 1
	File
	Parameter
)

var (
	_ json.Marshaler   = (*source)(nil)
	_ json.Unmarshaler = (*source)(nil)
)

func (src source) MarshalJSON() ([]byte, error) {
	switch src {
	case Environment:
		return []byte(`"ENVIRONMENT_VARIABLE"`), nil
	case File:
		return []byte(`"FILE"`), nil
	case Parameter:
		return []byte(`"CONFIGURATION_PARAMETER"`), nil
	default:
		return nil, errorUnknownSource
	}
}

func (src *source) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"ENVIRONMENT_VARIABLE"`:
		*src = Environment
	case `"FILE"`:
		*src = File
	case `"CONFIGURATION_PARAMETER"`:
		*src = Parameter
	default:
		return errorUnknownSource
	}
	return nil
}
