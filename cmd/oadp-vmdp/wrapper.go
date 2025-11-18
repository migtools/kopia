package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/kopia/kopia/cli"
)

// oadpAppWrapper wraps the Kopia CLI app with OADP-specific customizations.
type oadpAppWrapper struct {
	kopiaApp *cli.App
	kp       *kingpin.Application
}

// newOADPWrapper creates a new OADP wrapper around a Kopia CLI app.
func newOADPWrapper(app *cli.App, kp *kingpin.Application) *oadpAppWrapper {
	return &oadpAppWrapper{
		kopiaApp: app,
		kp:       kp,
	}
}

// setup configures the OADP-specific customizations.
func (w *oadpAppWrapper) setup() {
	// Apply OADP branding and help text customizations
	w.customizeHelpText()

	// Configure OADP-specific defaults
	w.setOADPDefaults()
}

// customizeHelpText modifies help text to use OADP terminology.
func (w *oadpAppWrapper) customizeHelpText() {
	// The kingpin application already has OADP branding from main.go
	// Help is handled by help_handler.go which shows clean, filtered output
}

// setOADPDefaults sets OADP-specific default values.
func (w *oadpAppWrapper) setOADPDefaults() {
	// OADP defaults are primarily handled through environment variables
	// and command-line flag defaults in the base Kopia CLI.
	// This function is reserved for future OADP-specific defaults.
}
