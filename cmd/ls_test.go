package cmd

import (
	"github.com/stretchr/testify/assert"
	"io"
	"menv/config"
	"menv/profiles"
	"os"
	"testing"
)

func TestPrintProfilesNone(t *testing.T) {

	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printProfiles(make([]string, 0))
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout

	assert.Contains(t, output, "No profiles found")

}

func TestPrintProfilesOneActive(t *testing.T) {
	configDir := t.TempDir()
	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
		Verbose:  false,
	}
	config.Set(testCfg)
	profiles.Init(testCfg)
	_ = os.Chdir(t.TempDir())

	_ = profiles.Create("active")
	_ = profiles.Create("non_active")
	_ = profiles.Set("active")

	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printProfiles(profiles.Profiles())
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout

	assert.Contains(t, output, "* active")
	assert.Contains(t, output, "  non_active")

}
