package standard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Run executes and takes full control of the configured application. This
// function never returns, and will exit with an appropriate status code
// indicating success or failure.
func Run(app Application) {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          app.Name,
		Version:      app.Version,
		Long:         app.Description + "\n\n" + app.Homepage,
	}

	cmd.AddCommand(
		commandCreate(app),
		commandDestroy(app),
		commandManifest(app),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1) // nolint:gomnd
	}
	os.Exit(0)
}

func commandCreate(app Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create cluster",
		Long:  "Creates a cluster",

		RunE: func(cmd *cobra.Command, _ []string) error {
			values, err := resolveSpecs(cmd, app.Create.Inputs)
			if err != nil {
				return err
			}

			return app.Create.Handler(context.Background(), values)
		},
	}

	addCommandFlags(app.Create.Inputs, cmd)
	return cmd
}

func commandDestroy(app Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy cluster",
		Long:  "Destroys a cluster",

		RunE: func(cmd *cobra.Command, _ []string) error {
			values, err := resolveSpecs(cmd, app.Destroy.Inputs)
			if err != nil {
				return err
			}

			return app.Destroy.Handler(context.Background(), values)
		},
	}

	addCommandFlags(app.Destroy.Inputs, cmd)
	return cmd
}

func commandManifest(app Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Flavor manifest",
		Long:  "Print a JSON manifest",

		RunE: func(cmd *cobra.Command, _ []string) error {
			sortSpecs(app.Create.Inputs)
			sortSpecs(app.Destroy.Inputs)

			manifest := Manifest{
				Create:  app.Create,
				Destroy: app.Destroy,
				Metadata: Metadata{
					Name:        app.Name,
					Description: app.Description,
					Version:     app.Version,
					Homepage:    app.Homepage,
				},
				Version: "v1.0",
			}

			data, err := json.MarshalIndent(manifest, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}

	return cmd
}

func addCommandFlags(specs []Parameter, cmd *cobra.Command) {
	for _, spec := range specs {
		if spec.Source != Flag {
			continue
		}

		cmd.Flags().String(spec.Name, "", spec.Description)
		cmd.MarkFlagRequired(spec.Name) // nolint:errcheck
	}
}

func resolveSpecs(cmd *cobra.Command, specs []Parameter) (map[string]string, error) {
	if err := validate(specs, cmd); err != nil {
		return nil, err
	}

	resolved := make(map[string]string, len(specs))
	for _, spec := range specs {
		value, err := spec.Resolve(cmd)
		if err != nil {
			return nil, err
		}
		resolved[spec.Name] = value
	}

	return resolved, nil
}

func validate(specs []Parameter, cmd *cobra.Command) error {
	results := validateSpecs(specs, cmd)
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
		return errors.New("validation failed")
	}
	return nil
}
