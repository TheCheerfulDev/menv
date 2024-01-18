package cmd

import (
	"github.com/stretchr/testify/assert"
	"menv/config"
	"menv/profiles"
	"os"
	"testing"
)

func TestIsActiveProfileUsedInWorkspace(t *testing.T) {
	initRmIdeaTest(t)
	activeProfile := "test"
	_ = profiles.Create(activeProfile)
	_ = profiles.Set(activeProfile)

	_ = os.Mkdir(".idea", 0755)
	_ = writeWorkspaceTemplate(activeProfile)

	assert.True(t, isProfileUsedInWorkspace(activeProfile))
	assert.False(t, isProfileUsedInWorkspace("other"))
}

func TestRemoveProfileFromWorkspace(t *testing.T) {
	initRmIdeaTest(t)
	activeProfile := "test"
	_ = profiles.Create(activeProfile)
	_ = profiles.Set(activeProfile)

	_ = os.Mkdir(".idea", 0755)
	_ = writeWorkspaceTemplate(activeProfile)

	assert.True(t, isProfileUsedInWorkspace(activeProfile))
	err := removeProfileFromWorkspace(activeProfile)
	assert.NoError(t, err)
	assert.False(t, isProfileUsedInWorkspace(activeProfile))
}

func initRmIdeaTest(t *testing.T) {
	testDir := t.TempDir()
	testCfg := config.Config{
		MenvRoot: testDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)
	_ = os.Chdir(t.TempDir())
}
