package standard

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Create   ActionManifest `yaml:"create"`
	Destroy  ActionManifest `yaml:"destroy"`
	Metadata Metadata       `yaml:"metadata"`
	Version  string         `yaml:"version"`
}

type ActionManifest struct {
	Inputs []Spec `yaml:"inputs"`
}

type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Homepage    string `yaml:"homepage"`
}

func LoadManifest(filename string) (Manifest, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := yaml.UnmarshalStrict(data, &manifest); err != nil {
		return Manifest{}, nil
	}

	return manifest, nil
}
