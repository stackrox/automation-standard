package standard

import (
	"encoding/json"
	"errors"
)

var errorUnknownSource = errors.New("unknown source")

type source int

const (
	// Environment represents a parameter that exists in the current working
	// environment.
	Environment source = iota + 1

	// File represents a parameter that exists as a file.
	File

	// Flag represents a parameter that is given as a command line flag.
	Flag
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
	case Flag:
		return []byte(`"FLAG"`), nil
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
	case `"FLAG"`:
		*src = Flag
	default:
		return errorUnknownSource
	}
	return nil
}
