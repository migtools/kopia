package main

import (
	"os"
	"testing"
)

func TestTranslateEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name           string
		setup          map[string]string // env vars to set before test
		expectedKopia  map[string]string // expected KOPIA_* vars after translation
		cleanup        []string          // env vars to unset after test
	}{
		{
			name: "translate OADP_PASSWORD to KOPIA_PASSWORD",
			setup: map[string]string{
				"OADP_PASSWORD": "secret123",
			},
			expectedKopia: map[string]string{
				"KOPIA_PASSWORD": "secret123",
			},
			cleanup: []string{"OADP_PASSWORD", "KOPIA_PASSWORD"},
		},
		{
			name: "translate OADP_CONFIG_PATH to KOPIA_CONFIG_PATH",
			setup: map[string]string{
				"OADP_CONFIG_PATH": "/path/to/config",
			},
			expectedKopia: map[string]string{
				"KOPIA_CONFIG_PATH": "/path/to/config",
			},
			cleanup: []string{"OADP_CONFIG_PATH", "KOPIA_CONFIG_PATH"},
		},
		{
			name: "OADP var overrides existing KOPIA var",
			setup: map[string]string{
				"OADP_PASSWORD":  "oadp-secret",
				"KOPIA_PASSWORD": "kopia-secret",
			},
			expectedKopia: map[string]string{
				"KOPIA_PASSWORD": "oadp-secret",
			},
			cleanup: []string{"OADP_PASSWORD", "KOPIA_PASSWORD"},
		},
		{
			name: "translate multiple environment variables",
			setup: map[string]string{
				"OADP_PASSWORD":        "secret",
				"OADP_CACHE_DIRECTORY": "/cache",
				"OADP_LOG_DIR":         "/logs",
			},
			expectedKopia: map[string]string{
				"KOPIA_PASSWORD":        "secret",
				"KOPIA_CACHE_DIRECTORY": "/cache",
				"KOPIA_LOG_DIR":         "/logs",
			},
			cleanup: []string{
				"OADP_PASSWORD", "KOPIA_PASSWORD",
				"OADP_CACHE_DIRECTORY", "KOPIA_CACHE_DIRECTORY",
				"OADP_LOG_DIR", "KOPIA_LOG_DIR",
			},
		},
		{
			name:  "no OADP vars set does not clear existing KOPIA vars",
			setup: map[string]string{
				"KOPIA_PASSWORD": "existing",
			},
			expectedKopia: map[string]string{
				"KOPIA_PASSWORD": "existing",
			},
			cleanup: []string{"KOPIA_PASSWORD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			for k, v := range tt.setup {
				os.Setenv(k, v)
			}

			// Execute translation
			translateEnvironmentVariables()

			// Verify expected KOPIA vars
			for kopiaVar, expectedValue := range tt.expectedKopia {
				actualValue := os.Getenv(kopiaVar)
				if actualValue != expectedValue {
					t.Errorf("Expected %s=%q, got %q", kopiaVar, expectedValue, actualValue)
				}
			}

			// Cleanup
			for _, v := range tt.cleanup {
				os.Unsetenv(v)
			}
		})
	}
}

func TestTranslateSingleEnvVar(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectedExists bool
	}{
		{
			name:           "translate OADP_PASSWORD",
			input:          "OADP_PASSWORD",
			expectedOutput: "KOPIA_PASSWORD",
			expectedExists: true,
		},
		{
			name:           "translate OADP_CACHE_DIRECTORY",
			input:          "OADP_CACHE_DIRECTORY",
			expectedOutput: "KOPIA_CACHE_DIRECTORY",
			expectedExists: true,
		},
		{
			name:           "automatic translation for unmapped OADP var",
			input:          "OADP_CUSTOM_VAR",
			expectedOutput: "KOPIA_CUSTOM_VAR",
			expectedExists: true,
		},
		{
			name:           "no translation for non-OADP var",
			input:          "SOME_OTHER_VAR",
			expectedOutput: "SOME_OTHER_VAR",
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, exists := translateSingleEnvVar(tt.input)
			if output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
			if exists != tt.expectedExists {
				t.Errorf("Expected exists=%v, got %v", tt.expectedExists, exists)
			}
		})
	}
}

func TestGetEnvPrefix(t *testing.T) {
	prefix := getEnvPrefix()
	if prefix != "OADP" {
		t.Errorf("Expected prefix 'OADP', got %q", prefix)
	}
}

func TestGetAllOADPEnvVars(t *testing.T) {
	vars := getAllOADPEnvVars()

	// Check that we have some variables
	if len(vars) == 0 {
		t.Error("Expected at least one OADP environment variable")
	}

	// Check that common vars are present
	expectedVars := []string{
		"OADP_PASSWORD",
		"OADP_CONFIG_PATH",
		"OADP_CACHE_DIRECTORY",
	}

	varMap := make(map[string]bool)
	for _, v := range vars {
		varMap[v] = true
	}

	for _, expected := range expectedVars {
		if !varMap[expected] {
			t.Errorf("Expected %q to be in the list of OADP vars", expected)
		}
	}
}
