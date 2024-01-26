package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDefault(t *testing.T) {
	home, _ := os.UserHomeDir()
	expected := Config{
		MenvRoot: filepath.Join(home, ".config", "menv"),
		Editor:   "vi",
	}
	actual := Default()

	assert.Equal(t, expected, actual)
}

func TestEditor(t *testing.T) {
	tests := []struct {
		editor string
		fn     func(string)
	}{
		{"vi", func(s string) {}},
		{"nano", func(s string) {
			os.Setenv("MENV_EDITOR", s)
		}},
	}
	defer os.Unsetenv("MENV_EDITOR")
	cfg := Default()
	Set(cfg)

	for _, test := range tests {
		test.fn(test.editor)
		expected := test.editor
		actual := Editor()
		assert.Equal(t, expected, actual)
	}
}

func TestSetGet(t *testing.T) {
	expected := Default()
	Set(expected)
	actual := Get()

	assert.Equal(t, expected, actual)
}

func TestVerbose(t *testing.T) {
	tests := []struct {
		verbose string
		fn      func(string)
	}{
		{"false", func(s string) {}},
		{"true", func(s string) {
			os.Setenv("MENV_VERBOSE", s)
		}},
	}
	defer os.Unsetenv("MENV_VERBOSE")
	cfg := Default()
	Set(cfg)

	for _, test := range tests {
		test.fn(test.verbose)
		expected, _ := strconv.ParseBool(test.verbose)
		actual := Verbose()
		assert.Equal(t, expected, actual)
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	testCfg := Config{
		MenvRoot: filepath.Join(dir, ".config", "menv"),
		Editor:   "vi",
	}
	Set(testCfg)

	err := Init()
	assert.NoError(t, err)
}
