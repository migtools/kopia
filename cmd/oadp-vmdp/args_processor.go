package main

import (
	"fmt"
	"os"
	"strings"
)

// argsProcessor handles preprocessing of command-line arguments before they reach kingpin.
type argsProcessor struct{}

// newArgsProcessor creates a new arguments processor.
func newArgsProcessor() *argsProcessor {
	return &argsProcessor{}
}

// processS3Prefix modifies os.Args to automatically prepend "oadp-vmdp/" to S3 --prefix flags.
// This allows transparent prefix normalization without modifying core Kopia code.
func (ap *argsProcessor) processS3Prefix() error {
	// Look for S3-related commands and --prefix flag
	inS3Context := false

	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]

		// Check if we're in an S3 context (repository/bsl create/connect s3)
		if arg == "s3" {
			inS3Context = true
			continue
		}

		// If we're in S3 context, look for --prefix flag
		if inS3Context {
			// Handle --prefix=value format
			if strings.HasPrefix(arg, "--prefix=") {
				parts := strings.SplitN(arg, "=", 2)
				if len(parts) == 2 {
					userPrefix := parts[1]
					normalized, err := NormalizeOADPPrefix(userPrefix)
					if err != nil {
						return err
					}
					os.Args[i] = fmt.Sprintf("--prefix=%s", normalized)
				}
				continue
			}

			// Handle --prefix value format (two separate args)
			if arg == "--prefix" {
				if i+1 < len(os.Args) {
					userPrefix := os.Args[i+1]
					normalized, err := NormalizeOADPPrefix(userPrefix)
					if err != nil {
						return err
					}
					os.Args[i+1] = normalized
					i++ // Skip the next arg since we just processed it
				}
				continue
			}
		}
	}

	return nil
}

// processAllArgs runs all argument preprocessing steps.
func (ap *argsProcessor) processAllArgs() error {
	// Process S3 prefix normalization
	if err := ap.processS3Prefix(); err != nil {
		return err
	}

	return nil
}
