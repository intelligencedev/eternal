// internal/python/executor.go

package python

import (
	"os/exec"
)

// ExecuteScript runs a Python script with the given arguments.
func ExecuteScript(scriptPath string, args ...string) (string, error) {
	cmd := exec.Command("python3", append([]string{scriptPath}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
