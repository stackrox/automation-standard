package standard

import (
	"context"
	"time"
)

// Application represents the full configuration for a runnable automation
// flavor.
type Application struct {
	// Name is a human-readable name for this application.
	//
	// Example: "Example Cluster"
	Name string

	// Description is a human-readable Description for this application.
	//
	// Example: "Provisions an Example cluster"
	Description string

	// Homepage is a URL that links to this application's GitHub project.
	//
	// Example: "https://github.com/stackrox/example-cluster-automation"
	Homepage string

	// Version is the specific version of this application.
	//
	// Example: "v1.2.3"
	Version string

	// Create is the extended configuration for this applications "create"
	// action.
	Create ActionConfiguration

	// Destroy is the extended configuration for this applications "destroy"
	// action.
	Destroy ActionConfiguration
}

type ActionConfiguration struct {
	Inputs  []Spec
	Handler Handler
	Timeout time.Duration
}

type Handler func(context.Context, map[string]string) error
