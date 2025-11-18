package main

import (
	"os"
	"strings"
)

// envMapping defines the translation from OADP environment variables to Kopia equivalents.
var envMapping = map[string]string{
	// Core configuration
	"OADP_PASSWORD":                "KOPIA_PASSWORD",
	"OADP_CONFIG_PATH":             "KOPIA_CONFIG_PATH",
	"OADP_CACHE_DIRECTORY":         "KOPIA_CACHE_DIRECTORY",
	"OADP_CHECK_FOR_UPDATES":       "KOPIA_CHECK_FOR_UPDATES",
	"OADP_PERSIST_CREDENTIALS":     "KOPIA_PERSIST_CREDENTIALS_ON_CONNECT",

	// Logging configuration
	"OADP_LOG_DIR":                 "KOPIA_LOG_DIR",
	"OADP_LOG_DIR_MAX_FILES":       "KOPIA_LOG_DIR_MAX_FILES",
	"OADP_LOG_DIR_MAX_AGE":         "KOPIA_LOG_DIR_MAX_AGE",
	"OADP_LOG_DIR_MAX_SIZE_MB":     "KOPIA_LOG_DIR_MAX_SIZE_MB",
	"OADP_LOG_FILE_MAX_SEGMENT_SIZE": "KOPIA_LOG_FILE_MAX_SEGMENT_SIZE",
	"OADP_CONTENT_LOG_DIR_MAX_FILES": "KOPIA_CONTENT_LOG_DIR_MAX_FILES",
	"OADP_CONTENT_LOG_DIR_MAX_AGE":   "KOPIA_CONTENT_LOG_DIR_MAX_AGE",
	"OADP_CONTENT_LOG_DIR_MAX_SIZE_MB": "KOPIA_CONTENT_LOG_DIR_MAX_SIZE_MB",
	"OADP_FILE_LOG_LOCAL_TZ":       "KOPIA_FILE_LOG_LOCAL_TZ",
	"OADP_FORCE_COLOR":             "KOPIA_FORCE_COLOR",
	"OADP_DISABLE_COLOR":           "KOPIA_DISABLE_COLOR",
	"OADP_CONSOLE_TIMESTAMPS":      "KOPIA_CONSOLE_TIMESTAMPS",
	"OADP_DISABLE_INTERNAL_LOG":    "KOPIA_DISABLE_INTERNAL_LOG",

	// Advanced configuration
	"OADP_INITIAL_UPDATE_CHECK_DELAY": "KOPIA_INITIAL_UPDATE_CHECK_DELAY",
	"OADP_UPDATE_CHECK_INTERVAL":      "KOPIA_UPDATE_CHECK_INTERVAL",
	"OADP_UPDATE_NOTIFY_INTERVAL":     "KOPIA_UPDATE_NOTIFY_INTERVAL",
	"OADP_TRACK_RELEASABLE":           "KOPIA_TRACK_RELEASABLE",
	"OADP_DUMP_ALLOCATOR_STATS":       "KOPIA_DUMP_ALLOCATOR_STATS",
	"OADP_REPO_UPGRADE_OWNER_ID":      "KOPIA_REPO_UPGRADE_OWNER_ID",
	"OADP_REPO_UPGRADE_NO_BLOCK":      "KOPIA_REPO_UPGRADE_NO_BLOCK",
	"OADP_SEND_ERROR_NOTIFICATIONS":   "KOPIA_SEND_ERROR_NOTIFICATIONS",
	"OADP_RESTORE_CONSISTENT_ATTRIBUTES": "KOPIA_RESTORE_CONSISTENT_ATTRIBUTES",
	"OADP_USE_KEYRING":                "KOPIA_USE_KEYRING",
}

// translateEnvironmentVariables translates OADP_* environment variables to KOPIA_*
// equivalents. This allows users to use OADP-branded environment variables while
// maintaining compatibility with Kopia's internal expectations.
//
// If both OADP_X and KOPIA_X are set, OADP_X takes precedence.
func translateEnvironmentVariables() {
	for oadpEnv, kopiaEnv := range envMapping {
		if oadpValue, exists := os.LookupEnv(oadpEnv); exists {
			// OADP env var is set, use it (overriding any existing KOPIA var)
			os.Setenv(kopiaEnv, oadpValue)
		}
	}
}

// getEnvPrefix returns "OADP" for use in environment variable naming.
// This is used by the CLI to generate help text and flag descriptions.
func getEnvPrefix() string {
	return "OADP"
}

// translateSingleEnvVar translates a single OADP environment variable name to its
// Kopia equivalent. Returns the Kopia name and true if a mapping exists,
// or the original name and false if no mapping exists.
func translateSingleEnvVar(oadpEnvName string) (string, bool) {
	if kopiaName, exists := envMapping[oadpEnvName]; exists {
		return kopiaName, true
	}

	// If it starts with OADP_, try automatic translation
	if strings.HasPrefix(oadpEnvName, "OADP_") {
		return "KOPIA_" + strings.TrimPrefix(oadpEnvName, "OADP_"), true
	}

	return oadpEnvName, false
}

// getAllOADPEnvVars returns a list of all supported OADP environment variable names.
// Useful for documentation and help text generation.
func getAllOADPEnvVars() []string {
	vars := make([]string, 0, len(envMapping))
	for k := range envMapping {
		vars = append(vars, k)
	}
	return vars
}
