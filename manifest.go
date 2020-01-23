package standard

import (
	"encoding/json"
	"io/ioutil"
)

type Manifest struct {
	Create   ActionManifest `json:"create"`
	Destroy  ActionManifest `json:"destroy"`
	Metadata Metadata       `json:"metadata"`
	Version  string         `json:"version"`
}

type ActionManifest struct {
	Inputs []Spec `json:"inputs"`
}

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Homepage    string `json:"homepage"`
}

func LoadManifest(filename string) (Manifest, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}
