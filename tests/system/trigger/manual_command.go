//go:build system

package trigger

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// ManualCommand executes a CLI command and returns the output.
// Fails the test if the command returns a non-zero exit code.
func ManualCommand(t *testing.T, args ...string) string {
	t.Helper()

	output, err := ManualCommandWithError(t, args...)
	require.NoError(t, err, "command should succeed: %s", output)

	return output
}

// ManualCommandWithError executes a CLI command and returns the output and error.
// Use this when testing error scenarios.
func ManualCommandWithError(t *testing.T, args ...string) (string, error) {
	t.Helper()

	root, err := projectRoot()
	if err != nil {
		t.Fatalf("ManualCommand: failed to find project root: %v", err)
	}

	binaryPath := filepath.Join(root, "bin", "system-test-app")

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "ENVIRONMENT=test")
	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	return string(output), err
}
