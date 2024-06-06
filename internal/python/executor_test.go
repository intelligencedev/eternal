package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteScript(t *testing.T) {
	output, err := ExecuteScript("test_script.py", "arg1", "arg2")
	assert.NoError(t, err)
	assert.Contains(t, output, "expected output")
}
