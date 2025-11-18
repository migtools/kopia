// Package main implements the OADP VM Data Protection CLI wrapper around Kopia.
//
// This is a thin wrapper that provides OADP-specific command names and environment
// variables while using Kopia's core functionality underneath.
//
// Usage:
//
//	$ oadp-vmdp [<flags>] <subcommand> [<args> ...]
//
// Use 'oadp-vmdp help' to see more details.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"

	"github.com/kopia/kopia/cli"
	"github.com/kopia/kopia/internal/logfile"
	"github.com/kopia/kopia/repo"
)

func main() {
	// Step 1: Check if we should show custom help
	helpHandler := newHelpHandler()
	if helpHandler.checkAndShowHelp() {
		os.Exit(0)
	}

	// Step 2: Check if command should be blocked (hidden commands)
	interceptor := newCommandInterceptor()
	if interceptor.checkAndBlock() {
		os.Exit(1)
	}

	// Step 3: Process arguments (S3 prefix normalization, etc.)
	processor := newArgsProcessor()
	if err := processor.processAllArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Translate command aliases (bsl → repository, backup → snapshot)
	translateArgs()

	// Step 5: Translate OADP_* environment variables to KOPIA_* for compatibility
	translateEnvironmentVariables()

	// Step 5: Create the base Kopia CLI app
	app := cli.NewApp()

	// Step 6: Create kingpin application with OADP branding
	kp := kingpin.New("oadp-vmdp", "OADP VM Data Protection - Virtual Machine Data Protection for OpenShift Virtualization").
		Author("Red Hat, Inc. <https://www.redhat.com/>")

	kp.Version(repo.BuildVersion + " build: " + repo.BuildInfo + " from: " + repo.BuildGitHubRepo)

	// Step 7: Attach logfile support
	logfile.Attach(app, kp)

	// Step 8: Setup the OADP wrapper with customizations
	wrapper := newOADPWrapper(app, kp)
	wrapper.setup()

	// Step 9: Configure kingpin output
	kp.ErrorWriter(os.Stderr)
	kp.UsageWriter(os.Stdout)

	// Step 10: Setup the Kopia CLI
	app.Attach(kp)

	// Step 11: Parse arguments and execute
	kingpin.MustParse(kp.Parse(os.Args[1:]))
}
