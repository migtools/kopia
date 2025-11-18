package main

import (
	"os"
	"strings"
)

// commandAliases maps OADP-friendly command names to their Kopia equivalents.
var commandAliases = map[string]string{
	"bsl":    "repository", // BSL (Backup Storage Location) -> repository
	"bsls":   "repository", // Alternative plural form
	"backup": "snapshot",   // backup -> snapshot
	"bkp":    "snapshot",   // Alternative short form
}

// mapCommandArgs translates OADP command names in os.Args to Kopia equivalents.
// This allows users to type "oadp-vmdp bsl connect s3 ..." which gets translated
// to "kopia repository connect s3 ..." internally.
//
// The function modifies os.Args in-place and returns the modified slice.
func mapCommandArgs(args []string) []string {
	if len(args) < 2 {
		return args // No command to map
	}

	// Look for the first command (skip flags)
	for i := 1; i < len(args); i++ {
		arg := args[i]

		// Skip flags (anything starting with -)
		if strings.HasPrefix(arg, "-") {
			continue
		}

		// Check if this is an aliased command
		if kopiaCmd, exists := commandAliases[arg]; exists {
			// Replace the OADP command with the Kopia equivalent
			args[i] = kopiaCmd
			break
		}

		// Stop at the first non-flag argument (the command)
		break
	}

	return args
}

// getCommandMapping returns the full command alias map.
// Useful for documentation and testing.
func getCommandMapping() map[string]string {
	// Return a copy to prevent external modification
	mapping := make(map[string]string, len(commandAliases))
	for k, v := range commandAliases {
		mapping[k] = v
	}
	return mapping
}

// isOADPCommand checks if a given command is an OADP-specific alias.
func isOADPCommand(cmd string) bool {
	_, exists := commandAliases[cmd]
	return exists
}

// getKopiaCommand returns the Kopia command for a given OADP command.
// Returns the original command if no mapping exists.
func getKopiaCommand(oadpCmd string) string {
	if kopiaCmd, exists := commandAliases[oadpCmd]; exists {
		return kopiaCmd
	}
	return oadpCmd
}

// translateArgs translates the entire argument list, handling both
// command aliases and any future argument transformations.
func translateArgs() {
	os.Args = mapCommandArgs(os.Args)
}
