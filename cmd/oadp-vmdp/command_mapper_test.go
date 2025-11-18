package main

import (
	"reflect"
	"testing"
)

func TestMapCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "map bsl to repository",
			input:    []string{"oadp-vmdp", "bsl", "connect", "s3"},
			expected: []string{"oadp-vmdp", "repository", "connect", "s3"},
		},
		{
			name:     "map bsls to repository",
			input:    []string{"oadp-vmdp", "bsls", "status"},
			expected: []string{"oadp-vmdp", "repository", "status"},
		},
		{
			name:     "map backup to snapshot",
			input:    []string{"oadp-vmdp", "backup", "create", "/path"},
			expected: []string{"oadp-vmdp", "snapshot", "create", "/path"},
		},
		{
			name:     "map bkp to snapshot",
			input:    []string{"oadp-vmdp", "bkp", "list"},
			expected: []string{"oadp-vmdp", "snapshot", "list"},
		},
		{
			name:     "no mapping for standard kopia commands",
			input:    []string{"oadp-vmdp", "cache", "info"},
			expected: []string{"oadp-vmdp", "cache", "info"},
		},
		{
			name:     "mapping with flags before command",
			input:    []string{"oadp-vmdp", "--log-level=debug", "bsl", "connect"},
			expected: []string{"oadp-vmdp", "--log-level=debug", "repository", "connect"},
		},
		{
			name:     "mapping with flags and short flags",
			input:    []string{"oadp-vmdp", "-v", "--config=/path", "backup", "create"},
			expected: []string{"oadp-vmdp", "-v", "--config=/path", "snapshot", "create"},
		},
		{
			name:     "no command (just binary name)",
			input:    []string{"oadp-vmdp"},
			expected: []string{"oadp-vmdp"},
		},
		{
			name:     "help command not mapped",
			input:    []string{"oadp-vmdp", "help"},
			expected: []string{"oadp-vmdp", "help"},
		},
		{
			name:     "only flags, no command",
			input:    []string{"oadp-vmdp", "--version"},
			expected: []string{"oadp-vmdp", "--version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test case
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			result := mapCommandArgs(input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("mapCommandArgs(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCommandMapping(t *testing.T) {
	mapping := getCommandMapping()

	// Check expected mappings
	expectedMappings := map[string]string{
		"bsl":    "repository",
		"bsls":   "repository",
		"backup": "snapshot",
		"bkp":    "snapshot",
	}

	for oadpCmd, expectedKopiaCmd := range expectedMappings {
		if kopiaCmd, exists := mapping[oadpCmd]; !exists {
			t.Errorf("Expected mapping for %q to exist", oadpCmd)
		} else if kopiaCmd != expectedKopiaCmd {
			t.Errorf("Expected %q -> %q, got %q -> %q", oadpCmd, expectedKopiaCmd, oadpCmd, kopiaCmd)
		}
	}

	// Test that the returned map is a copy (modifying it shouldn't affect the original)
	mapping["test"] = "modified"
	newMapping := getCommandMapping()
	if _, exists := newMapping["test"]; exists {
		t.Error("Modifying returned mapping affected the original")
	}
}

func TestIsOADPCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "bsl is OADP command",
			command:  "bsl",
			expected: true,
		},
		{
			name:     "bsls is OADP command",
			command:  "bsls",
			expected: true,
		},
		{
			name:     "backup is OADP command",
			command:  "backup",
			expected: true,
		},
		{
			name:     "bkp is OADP command",
			command:  "bkp",
			expected: true,
		},
		{
			name:     "repository is not OADP command",
			command:  "repository",
			expected: false,
		},
		{
			name:     "snapshot is not OADP command",
			command:  "snapshot",
			expected: false,
		},
		{
			name:     "cache is not OADP command",
			command:  "cache",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOADPCommand(tt.command)
			if result != tt.expected {
				t.Errorf("isOADPCommand(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestGetKopiaCommand(t *testing.T) {
	tests := []struct {
		name        string
		oadpCommand string
		expected    string
	}{
		{
			name:        "bsl maps to repository",
			oadpCommand: "bsl",
			expected:    "repository",
		},
		{
			name:        "backup maps to snapshot",
			oadpCommand: "backup",
			expected:    "snapshot",
		},
		{
			name:        "unmapped command returns original",
			oadpCommand: "cache",
			expected:    "cache",
		},
		{
			name:        "empty string returns empty",
			oadpCommand: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getKopiaCommand(tt.oadpCommand)
			if result != tt.expected {
				t.Errorf("getKopiaCommand(%q) = %q, want %q", tt.oadpCommand, result, tt.expected)
			}
		})
	}
}
