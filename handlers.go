package standard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

var (
	errActionUnknown      = errors.New("action unknown")
	errIncorrectArguments = errors.New("incorrect arguments")
	errValidationFailed   = errors.New("validation failed")

	standardSpecs = []Spec{ // nolint:gochecknoglobals
		{
			Name:        "description",
			Description: "The description of this run",
			Source:      Parameter,
		},
		{
			Name:        "id",
			Description: "The id of this run",
			Source:      Parameter,
		},
		{
			Name:        "name",
			Description: "The name of this run",
			Source:      Parameter,
		},
		{
			Name:        "owner",
			Description: "Who the service launched this run on behalf of",
			Source:      Parameter,
		},
		{
			Name:        "service",
			Description: "The service that launched this run",
			Source:      Parameter,
		},
		{
			Name:        "url",
			Description: "Homepage to the service run",
			Source:      Parameter,
		},
	}
)

const configFilename = "config.yaml"

func Run(app Application) {
	if err := handle(app); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", app.Name, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func handle(app Application) error {
	if len(os.Args) != 2 {
		return errIncorrectArguments
	}

	switch os.Args[1] {
	case "create":
		return handleCreate(app)

	case "create-check":
		return handleCreateCheck(app)

	case "destroy":
		return handleDestroy(app)

	case "destroy-check":
		return handleDestroyCheck(app)

	case "manifest":
		return handleManifest(app)

	case "version":
		return handleVersion(app)

	default:
		return errActionUnknown
	}
}

func handleCreateCheck(app Application) error {
	return handleCheck(app.Create.Inputs)
}

func handleDestroyCheck(app Application) error {
	return handleCheck(app.Destroy.Inputs)
}

func handleCheck(inputs []Spec) error {
	cfg, err := LoadConfig(configFilename)
	if err != nil {
		return err
	}

	specs, err := combineWithStandardSpecs(inputs)
	if err != nil {
		return err
	}

	return ValidatePretty(specs, cfg)
}

func handleCreate(app Application) error {
	return handleAction(app.Create)
}

func handleDestroy(app Application) error {
	return handleAction(app.Destroy)
}

func handleAction(action ActionConfiguration) error {
	cfg, err := LoadConfig(configFilename)
	if err != nil {
		return err
	}

	resolved, err := resolveSpecs(cfg, action.Inputs)
	if err != nil {
		return err
	}

	ctx, cancel := lifespanContext(action.Timeout)
	defer cancel()
	return action.Handler(ctx, resolved)
}

func handleManifest(app Application) error {
	combinedCreateSpecs, err := combineWithStandardSpecs(app.Create.Inputs)
	if err != nil {
		return err
	}

	combinedDestroySpecs, err := combineWithStandardSpecs(app.Destroy.Inputs)
	if err != nil {
		return err
	}

	sortSpecs(combinedCreateSpecs)
	sortSpecs(combinedDestroySpecs)

	manifest := Manifest{
		Create: ActionManifest{
			Inputs: combinedCreateSpecs,
		},
		Destroy: ActionManifest{
			Inputs: combinedDestroySpecs,
		},
		Metadata: Metadata{
			Name:        app.Name,
			Description: app.Description,
			Version:     app.Version,
			Homepage:    app.Homepage,
		},
		Version: "v1.0",
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	fmt.Printf("%s", string(data))
	return nil
}

func handleVersion(app Application) error {
	fmt.Println(app.Version)
	return nil
}

func resolveSpecs(cfg Config, inputs []Spec) (map[string]string, error) {
	specs, err := combineWithStandardSpecs(inputs)
	if err != nil {
		return nil, err
	}

	if err := ValidatePretty(specs, cfg); err != nil {
		return nil, err
	}

	resolved := make(map[string]string, len(specs))
	for _, spec := range specs {
		value, err := spec.Resolve(cfg)
		if err != nil {
			return nil, err
		}
		resolved[spec.Name] = value
	}

	return resolved, nil
}

func lifespanContext(lifespan time.Duration) (context.Context, context.CancelFunc) {
	if lifespan == 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), lifespan)
}

func Validate(specs []Spec, cfg Config) error {
	results := validateSpecs(specs, cfg)
	for _, result := range results {
		if result.err != nil {
			return errValidationFailed
		}
	}

	return nil
}

func ValidatePretty(specs []Spec, cfg Config) error {
	results := validateSpecs(specs, cfg)
	sort.Slice(results, func(i, j int) bool {
		return results[i].name < results[j].name
	})

	var errorsEncountered bool
	for _, result := range results {
		if result.err != nil {
			errorsEncountered = true
			color.Red("[FAIL] %s", result.message)
			color.Red("       ↳ %v", result.err)
		} else {
			color.Green("[PASS] %s", result.message)
		}
	}
	if errorsEncountered {
		return errValidationFailed
	}
	return nil
}
