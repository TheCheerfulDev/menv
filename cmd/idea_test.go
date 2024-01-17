package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsNotMavenProject(t *testing.T) {
	_ = os.Chdir(t.TempDir())
	assert.True(t, isNotMavenproject())
	_ = os.WriteFile("pom.xml", []byte("test"), 0644)
	assert.False(t, isNotMavenproject())
}

func TestIsNotIntellijProject(t *testing.T) {
	_ = os.Chdir(t.TempDir())
	assert.True(t, isNotIntellijProject())
	_ = os.Mkdir(".idea", 0644)
	assert.False(t, isNotIntellijProject())
}
