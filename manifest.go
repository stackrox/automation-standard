package standard

// Manifest represents a (JSON-marshalable) manifest of an application's
// declared configuration and parameters.
type Manifest struct {
	Create   Action   `json:"create"`
	Destroy  Action   `json:"destroy"`
	Metadata Metadata `json:"metadata"`
	Version  string   `json:"version"`
}

// Metadata represents an application's metadata.
type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Homepage    string `json:"homepage"`
}
