package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stackrox/automation-standard"
)

var version = "development"

type args struct {
	action   string
	config   string
	manifest string
	quiet    bool
	version  bool
}

func main() {
	var args args
	flag.StringVar(&args.config, "config", "", "path to config file")
	flag.StringVar(&args.manifest, "manifest", "", "path to manifest file")
	flag.BoolVar(&args.quiet, "quiet", false, "silence all output")
	flag.BoolVar(&args.version, "version", false, fmt.Sprintf("print the version %q and exit", version))
	flag.Parse()
	args.action = flag.Arg(0)

	if err := mainCmd(args); err != nil {
		if !args.quiet {
			fmt.Fprintf(os.Stderr, "standard-validate: %v\n", err)
		}
		os.Exit(1)
	}
}

func mainCmd(args args) error {
	// The version flag "-version" was set. Print version and exit.
	if args.version {
		fmt.Println(version)
		return nil
	}

	// Load the manifest file, which contains create and destroy specs.
	manifest, err := standard.LoadManifest(args.manifest)
	if err != nil {
		return err
	}

	// Load the configuration file, which contains arbitrary runtime
	// configuration data.
	configuration, err := standard.LoadConfig(args.config)
	if err != nil {
		return err
	}

	var specs []standard.Spec
	switch args.action {
	case "create":
		// Use the create specs for validation.
		specs = manifest.Create.Inputs

	case "destroy":
		// Use the destroy specs for validation.
		specs = manifest.Destroy.Inputs

	default:
		return fmt.Errorf("unknown action")
	}

	// The quiet flag "-quiet" was set, so silently validate.
	if args.quiet {
		return standard.Validate(specs, configuration)
	}

	// Validate and pretty print results.
	return standard.ValidatePretty(specs, configuration)
}
