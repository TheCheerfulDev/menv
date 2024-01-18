package cmd

import (
	"github.com/stretchr/testify/assert"
	"menv/config"
	"menv/profiles"
	"os"
	"testing"
)

func TestSetNonExistent(t *testing.T) {
	initSetTest(t)

	profile := "test"
	err := setProfile(profile)
	assert.EqualError(t, err, "profile does not exist")
}

func TestSetExisting(t *testing.T) {
	initSetTest(t)

	profile := "test"
	_ = profiles.Create(profile)
	err := setProfile(profile)
	assert.NoError(t, err)
}

func initSetTest(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.Config{
		MenvRoot: tempDir,
		Editor:   "vi",
	}
	_ = os.Chdir(t.TempDir())
	profiles.Init(cfg)
}
