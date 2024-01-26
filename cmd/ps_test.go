package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestPrintActiveProfileNone(t *testing.T) {
	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printActiveProfile("", "")
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout

	assert.Contains(t, output, "Active profile:")
	assert.Contains(t, output, "none (default)")
}

func TestPrintActiveProfile(t *testing.T) {
	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	profile := "active_profile"
	path := "/path/to/active/menv_file"
	printActiveProfile(profile, path)
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout

	assert.Contains(t, output, "Active profile:")
	assert.Contains(t, output, fmt.Sprintf("  %v (set by %v)\n", profile, path))
}
