package standard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
			values, err := resolveSpecsAndReportErrors(cmd, app.Create.Inputs)
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
			values, err := resolveSpecsAndReportErrors(cmd, app.Destroy.Inputs)
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

func resolveSpecsAndReportErrors(cmd *cobra.Command, specs []Parameter) (map[string]string, error) {
	resolvedValues := make(map[string]string)
	var errorsEncountered bool

	sortSpecs(specs)
	for _, spec := range specs {
		// Resolve this single spec.
		if value, errs := spec.Resolve(cmd); len(errs) != 0 {
			// Report all of the errors.
			errorsEncountered = true
			color.Red("[FAIL] %s (%s)", spec.Name, spec.Description)
			for _, err := range errs {
				color.Red("       ↳ %v", err)
			}
		} else {
			// No errors to report.
			resolvedValues[spec.Name] = value
			color.Green("[PASS] %s (%s)", spec.Name, spec.Description)
		}
	}

	if errorsEncountered {
		return nil, errors.New("validation failed")
	}

	return resolvedValues, nil
}
