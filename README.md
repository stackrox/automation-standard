# Automation Standard

🤖 A micro-framework for building standardized cluster automation entrypoints

## Installing

You can fetch this library by running the following:

```bash
$ go get -u github.com/stackrox/automation-standard
```

## Motivations

Given that there is a need for various cluster automations (GKE, OpenShift, Istio, KOPS, etc), standardizing on how those automations are configured is beneficial.

Additionally, integrating these automations into hard-to-debug-systems leads to things failing in strange ways because they were partially or misconfigured.

This library is opinionated, explicit, and strict as an effort to eliminate as many unknowns across these service boundaries as possible.

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/stackrox/automation-standard"
)

func main() {
	cfg := standard.Application{
		Name:        "example",
		Description: "Automation Standard example application",
		Homepage:    "https://github.com/stackrox/automation-standard",
		Version:     "v0.0.0",
		
		Create: standard.ActionConfiguration{
			Inputs: []standard.Spec{
				{
					Name:        "GOOGLE_APPLICATION_CREDENTIALS",
					Description: "Location of GCP service account credential file",
					Source:      standard.Environment,
				},
				{
					Name:        "main-image",
					Description: "Main image tag",
					Source:      standard.Parameter,
				},
			},
			Handler: create,
		},
		
		Destroy: standard.ActionConfiguration{
			Inputs: []standard.Spec{
				{
					Name:        "GOOGLE_APPLICATION_CREDENTIALS",
					Description: "Location of GCP service account credential file",
					Source:      standard.Environment,
				},
			},
			Handler: destroy,
		},
	}

	standard.Run(cfg)
}

func create(ctx context.Context, parameters map[string]string) error {
	fmt.Println("Hello from create()")
	fmt.Printf("Deploying %s\n", parameters["main-image"])
	return nil
}

func destroy(ctx context.Context, parameters map[string]string) error {
	fmt.Println("Hello from destroy()")
	return nil
}
```
