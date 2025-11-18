package main

import (
	"fmt"
	"os"
)

// helpHandler provides custom help output for oadp-vmdp.
type helpHandler struct{}

// newHelpHandler creates a new help handler.
func newHelpHandler() *helpHandler {
	return &helpHandler{}
}

// shouldShowCustomHelp checks if we should intercept and show custom help.
func (h *helpHandler) shouldShowCustomHelp() bool {
	// Show custom help if:
	// 1. No arguments (just binary name)
	// 2. --help or -h flag at top level
	// 3. "help" command with no subcommand
	// 4. --help-full (intercept and show custom help instead of all commands)
	// 5. --help-long or --help-man (also intercept)

	if len(os.Args) == 1 {
		// Just "oadp-vmdp" with no args
		return true
	}

	for i, arg := range os.Args {
		// Intercept --help-full, --help-long, --help-man and show our custom help
		if arg == "--help-full" || arg == "--help-long" || arg == "--help-man" {
			return true
		}

		if arg == "--help" || arg == "-h" {
			// Only show custom help at the top level (no command specified before --help)
			// Check if there's a command before the help flag
			hasCommand := false
			for j := 1; j < i; j++ {
				if !isFlag(os.Args[j]) {
					hasCommand = true
					break
				}
			}
			if !hasCommand {
				return true
			}
		}

		if arg == "help" && i == 1 {
			// "oadp-vmdp help" (with no subcommand)
			if len(os.Args) == 2 {
				return true
			}
		}
	}
	return false
}

// isFlag checks if an argument is a flag (starts with -).
func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

// showCustomHelp displays custom, clean help output.
func (h *helpHandler) showCustomHelp() {
	helpText := `OADP Virtual Machine Data Protection for OpenShift Virtualization

Usage:
  oadp-vmdp <command> [flags] [arguments]

Available Commands:
  BSL (Backup Storage Location) Management:
    bsl create s3      Create S3 backup storage location
    bsl connect s3     Connect to existing S3 backup storage location
    bsl disconnect     Disconnect from backup storage location
    bsl status         Show backup storage location status

  Backup Management:
    backup create      Create a backup of files/directories
    backup list        List available backups
    backup delete      Delete a backup
    backup restore     Restore files from a backup (alias: restore)

  Cache Management:
    cache info         Display cache information
    cache set          Configure cache settings
    cache clear        Clear the cache

  Other Commands:
    help               Show this help message
    --version          Show version information

Common Flags:
  -p, --password=PASSWORD         Repository password (or use OADP_PASSWORD)
      --config-file=FILE          Config file path (default: repository.config)
      --log-level=LEVEL           Log level (default: info)

Environment Variables:
  OADP_PASSWORD          Repository password
  OADP_CONFIG_PATH       Configuration file path
  OADP_CACHE_DIRECTORY   Cache directory location
  OADP_LOG_DIR           Log file directory

Quick Start:
  1. Connect to backup storage:
     oadp-vmdp bsl create s3 --bucket=my-bucket --endpoint=s3.amazonaws.com \
       --access-key=KEY --secret-access-key=SECRET

  2. Create a backup:
     oadp-vmdp backup create /path/to/data

  3. List backups:
     oadp-vmdp backup list

  4. Restore from backup:
     oadp-vmdp backup restore /path/to/data

For detailed help on a specific command:
  oadp-vmdp <command> --help

Examples:
  oadp-vmdp bsl create s3 --help
  oadp-vmdp backup create --help
  oadp-vmdp backup list --help
`
	fmt.Fprint(os.Stdout, helpText)
}

// checkAndShowHelp checks if custom help should be shown and displays it if needed.
// Returns true if help was shown (and program should exit).
func (h *helpHandler) checkAndShowHelp() bool {
	if h.shouldShowCustomHelp() {
		h.showCustomHelp()
		return true
	}
	return false
}
