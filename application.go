package standard

import (
	"context"
	"time"
)

type Application struct {
	Name        string
	Description string
	Homepage    string
	Version     string
	Create      ActionConfiguration
	Destroy     ActionConfiguration
}

type ActionConfiguration struct {
	Inputs  []Spec
	Handler Handler
	Timeout time.Duration
}

type Handler func(context.Context, map[string]string) error
