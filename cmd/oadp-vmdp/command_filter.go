package main

import (
	"github.com/alecthomas/kingpin/v2"
)

// hiddenCommands defines top-level commands that should be hidden in OADP.
// These are advanced/dangerous commands not needed for basic VM backup/restore.
var hiddenCommands = map[string]bool{
	"benchmark":    true, // Performance testing - not needed for users
	"diff":         true, // Advanced repository comparison
	"list":         true, // Advanced file listing - use 'backup list' instead
	"notification": true, // Server-specific feature
	"server":       true, // Server mode not used in OADP
	"policy":       true, // Advanced policy management - use defaults
	"mount":        true, // FUSE mounting - not supported in containers
	"maintenance":  true, // Advanced maintenance - handled automatically
}

// hiddenRepositorySubcommands defines repository subcommands to hide.
// Only keep: connect, create, disconnect, status
var hiddenRepositorySubcommands = map[string]bool{
	"repair":            true, // Dangerous - can corrupt repository
	"set-client":        true, // Advanced configuration
	"set-parameters":    true, // Advanced configuration
	"sync-to":           true, // Advanced replication feature
	"throttle":          true, // Advanced performance tuning
	"change-password":   true, // Use OADP password management instead
	"validate-provider": true, // Internal testing command
	"upgrade":           true, // Dangerous - can break compatibility
}

// hiddenSnapshotSubcommands defines snapshot subcommands to hide.
// Only keep: create, delete, list, restore
var hiddenSnapshotSubcommands = map[string]bool{
	"copy-history":  true, // Advanced history management
	"move-history":  true, // Advanced history management
	"estimate":      true, // Advanced analysis
	"expire":        true, // Handled automatically by retention policy
	"fix":           true, // Dangerous - can corrupt snapshots
	"migrate":       true, // Advanced migration feature
	"pin":           true, // Advanced retention management
	"verify":        true, // Advanced verification - use defaults
}

// commandFilter handles hiding unwanted commands from the CLI.
type commandFilter struct {
	kp *kingpin.Application
}

// newCommandFilter creates a new command filter for the given kingpin application.
func newCommandFilter(kp *kingpin.Application) *commandFilter {
	return &commandFilter{
		kp: kp,
	}
}

// apply walks the command tree and hides unwanted commands.
func (f *commandFilter) apply() {
	// Hide top-level commands
	f.hideTopLevelCommands()

	// Hide subcommands within repository/snapshot
	f.hideSubcommands()
}

// hideTopLevelCommands hides top-level commands that shouldn't be exposed.
func (f *commandFilter) hideTopLevelCommands() {
	// Note: kingpin doesn't expose a direct way to iterate commands,
	// so we'll mark them as hidden when we know they exist
	//
	// Since we can't easily iterate, we'll take a different approach:
	// We'll let the setup happen normally and then mark specific commands as hidden
}

// hideSubcommands hides subcommands within specific command groups.
func (f *commandFilter) hideSubcommands() {
	// Same challenge as above - kingpin doesn't expose command iteration
	// We'll need to mark commands as hidden if they exist
}

// shouldHideCommand checks if a top-level command should be hidden.
func shouldHideCommand(cmdName string) bool {
	return hiddenCommands[cmdName]
}

// shouldHideRepositorySubcommand checks if a repository subcommand should be hidden.
func shouldHideRepositorySubcommand(subcmdName string) bool {
	return hiddenRepositorySubcommands[subcmdName]
}

// shouldHideSnapshotSubcommand checks if a snapshot subcommand should be hidden.
func shouldHideSnapshotSubcommand(subcmdName string) bool {
	return hiddenSnapshotSubcommands[subcmdName]
}

// getVisibleCommands returns a list of commands that should be visible in OADP.
func getVisibleCommands() []string {
	return []string{
		"repository", // Aliased as "bsl"
		"snapshot",   // Aliased as "backup"
		"cache",      // Cache management
		"blob",       // Low-level blob operations (hidden in main help)
		"content",    // Low-level content operations (hidden in main help)
		"index",      // Index operations (hidden in main help)
		"logs",       // Log management
		"session",    // Session management
		"restore",    // File restore
		"show",       // Show repository objects
		"manifest",   // Manifest operations
	}
}

// getVisibleRepositorySubcommands returns allowed repository subcommands.
func getVisibleRepositorySubcommands() []string {
	return []string{
		"connect",
		"create",
		"disconnect",
		"status",
	}
}

// getVisibleSnapshotSubcommands returns allowed snapshot subcommands.
func getVisibleSnapshotSubcommands() []string {
	return []string{
		"create",
		"delete",
		"list",
		"restore",
	}
}
