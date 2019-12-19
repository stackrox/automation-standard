package standard

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return Config{}, nil
	}

	return cfg, err
}
