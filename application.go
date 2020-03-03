// Package standard allows developers of different automation flavors to
// declaratively construct the create/destroy lifecycle of an application.
// This library will be augment that application with parameter sanity checking
// and other features.
package standard

import "context"

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
	Create Action

	// Destroy is the extended configuration for this applications "destroy"
	// action.
	Destroy Action
}

// Action represents the configuration of either the create or destroy action.
type Action struct {
	// Inputs is the list of parameters that this action requires.
	Inputs []Parameter `json:"inputs"`

	// Handler is the handler that is run for this action.
	Handler Handler `json:"-"`
}

// Handler represents the body of a runnable action.
type Handler func(context.Context, map[string]string) error
