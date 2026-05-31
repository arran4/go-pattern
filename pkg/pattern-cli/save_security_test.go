package pattern_cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arran4/go-pattern/dsl"
)

func TestSaveSecurity(t *testing.T) {
	// Setup temporary directory to run tests in (as CWD)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir, err := os.MkdirTemp("", "pattern-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Change directory to tempDir for test execution context
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Errorf("failed to restore CWD: %v", err)
		}
	}()

	// Create subdirectory for legitimate use
	if err := os.Mkdir("subdir", 0755); err != nil {
		t.Fatal(err)
	}

	// Prepare targets
	traversalTarget := filepath.Join(tempDir, "..", "traversal_pwned.png") // outside tempDir
	absTarget := filepath.Join(os.TempDir(), "abs_pwned.png")

	// Clean up potential leftovers (if test runs multiple times or failed before)
	os.Remove(traversalTarget)
	os.Remove(absTarget)
	defer os.Remove(traversalTarget)
	defer os.Remove(absTarget)

	testCases := []struct {
		name        string
		filename    string
		expectError bool // We expect error for security violations
		pathCreated string // Path expected to be created if successful (relative or absolute)
	}{
		{
			name:        "Traversal Attack",
			filename:    "../traversal_pwned.png",
			expectError: true,
			pathCreated: traversalTarget,
		},
		{
			name:        "Absolute Path Attack",
			filename:    absTarget,
			expectError: true,
			pathCreated: absTarget,
		},
		{
			name:        "Valid File",
			filename:    "valid.png",
			expectError: false,
			pathCreated: "valid.png",
		},
		{
			name:        "Valid Subdirectory File",
			filename:    filepath.Join("subdir", "valid_sub.png"),
			expectError: false,
			pathCreated: filepath.Join("subdir", "valid_sub.png"),
		},
	}

	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pipeline := "checkers red blue | save " + tc.filename
			err := process(pipeline, funcMap)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got nil", tc.name)
				}
				// Verify file was NOT created
				if _, statErr := os.Stat(tc.pathCreated); statErr == nil {
					t.Errorf("Security check failed: File %s was created!", tc.pathCreated)
					os.Remove(tc.pathCreated)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tc.name, err)
				}
				// Verify file WAS created
				if _, statErr := os.Stat(tc.pathCreated); os.IsNotExist(statErr) {
					t.Errorf("Expected file %s to be created, but it does not exist", tc.pathCreated)
				}
			}
		})
	}
}
