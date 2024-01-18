package cmd

import (
	"github.com/stretchr/testify/assert"
	"menv/config"
	"menv/profiles"
	"os"
	"path"
	"testing"
)

func TestCreateMavenDir(t *testing.T) {
	initMvnLocalTest(t)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_, err := os.Stat(path.Join(tempDir, ".mvn"))
	assert.True(t, os.IsNotExist(err), ".mvn folder should not exist")
	createMavenDir()
	_, err = os.Stat(path.Join(tempDir, ".mvn"))
	assert.False(t, os.IsNotExist(err), ".mvn folder should exist")
}

func TestWriteMavenConfig(t *testing.T) {
	initMvnLocalTest(t)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	profile := "test"
	_ = profiles.Create(profile)
	_ = profiles.Set(profile)

	createMavenDir()
	writeMavenConfig(profiles.File(profile))
	file, _ := os.ReadFile(".mvn/maven.config")
	actualConfig := string(file)
	expectedConfig := "--settings\n" + profiles.File(profile)
	assert.Equal(t, expectedConfig, actualConfig, "maven.config should be equal")
}

func TestWriteMavenOpts(t *testing.T) {
	initMvnLocalTest(t)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	profile := "test"
	_ = profiles.Create(profile)
	_ = profiles.Set(profile)

	expectedOpts := "-Xmx2g -Xms1g"
	_ = os.WriteFile(profiles.OptsFile(profile), []byte(expectedOpts), 0644)
	createMavenDir()
	writeMavenOpts(profiles.OptsFile(profile))
	file, _ := os.ReadFile(".mvn/jvm.config")
	actualOpts := string(file)
	assert.Equal(t, expectedOpts, actualOpts, "jvm.config should be equal")
}

func initMvnLocalTest(t *testing.T) {
	configDir := t.TempDir()
	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)
	_ = os.Chdir(t.TempDir())
}
