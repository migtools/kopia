package main

import (
	"fmt"
	"os"
	"strings"
)

// commandInterceptor checks if a command should be blocked before execution.
type commandInterceptor struct{}

// newCommandInterceptor creates a new command interceptor.
func newCommandInterceptor() *commandInterceptor {
	return &commandInterceptor{}
}

// checkAndBlock examines os.Args and blocks execution of hidden commands.
// Returns true if the command should be blocked, false otherwise.
func (ci *commandInterceptor) checkAndBlock() bool {
	if len(os.Args) < 2 {
		return false // No command, let it proceed (will show help)
	}

	// Skip flags to find the actual command
	var command string
	var subcommand string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		// Skip flags
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if command == "" {
			command = arg
		} else if subcommand == "" {
			subcommand = arg
			break // We have both command and subcommand
		}
	}

	// Check if this is a hidden top-level command
	if shouldHideCommand(command) {
		ci.printBlockedMessage(command, "")
		return true
	}

	// Check repository subcommands
	if (command == "repository" || command == "repo" || command == "bsl" || command == "bsls") && subcommand != "" {
		if shouldHideRepositorySubcommand(subcommand) {
			ci.printBlockedMessage(command, subcommand)
			return true
		}
	}

	// Check snapshot subcommands
	if (command == "snapshot" || command == "snap" || command == "backup" || command == "bkp") && subcommand != "" {
		if shouldHideSnapshotSubcommand(subcommand) {
			ci.printBlockedMessage(command, subcommand)
			return true
		}
	}

	return false // Command is allowed
}

// printBlockedMessage displays an error message for blocked commands.
func (ci *commandInterceptor) printBlockedMessage(command, subcommand string) {
	var fullCommand string
	if subcommand != "" {
		fullCommand = fmt.Sprintf("%s %s", command, subcommand)
	} else {
		fullCommand = command
	}

	fmt.Fprintf(os.Stderr, "Error: The command '%s' is not available in oadp-vmdp.\n", fullCommand)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "OADP VM Data Protection provides a simplified interface for VM backup and restore operations.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Available commands:\n")
	fmt.Fprintf(os.Stderr, "  oadp-vmdp bsl create|connect|disconnect|status   - Manage Backup Storage Location\n")
	fmt.Fprintf(os.Stderr, "  oadp-vmdp backup create|list|restore|delete      - Manage backups\n")
	fmt.Fprintf(os.Stderr, "  oadp-vmdp cache info|set                         - Manage cache\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "For more information: oadp-vmdp --help\n")
	fmt.Fprintf(os.Stderr, "\n")
}
