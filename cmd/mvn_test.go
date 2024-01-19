package cmd

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"menv/config"
	"menv/profiles"
	"os"
	"path/filepath"
	"testing"
)

type MockShellCommand struct {
	*mock.Mock
}

func (m MockShellCommand) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m MockShellCommand) Stdin(stdin io.Reader) {
	m.Called(stdin)
}

func (m MockShellCommand) Stdout(stdout io.Writer) {
	m.Called(stdout)
}

func (m MockShellCommand) Stderr(stderr io.Writer) {
	m.Called(stderr)
}

func (m MockShellCommand) Output() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestFindMvnNoCellar(t *testing.T) {
	initMvnTest(t)

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	mockShell.On("Output").Return([]byte(""), errors.New("could not find maven in (home)brew cellar"))

	_, err := findMaven(mockProvider)
	assert.EqualError(t, err, "could not find maven in (home)brew cellar")
	mockShell.AssertExpectations(t)
}

func TestFindMvnCellarNoMvn(t *testing.T) {
	initMvnTest(t)

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	tempDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tempDir, "maven"), 0755)
	mockShell.On("Output").Return([]byte(tempDir), nil)

	_, err := findMaven(mockProvider)
	assert.EqualError(t, err, "could not find maven in (home)brew cellar")
	mockShell.AssertExpectations(t)
}

func TestFindMvn(t *testing.T) {
	initMvnTest(t)
	tempDir := t.TempDir()
	mvnDir := filepath.Join(tempDir, "maven", "3.9.6", "bin")
	_ = os.MkdirAll(mvnDir, 0755)
	os.Create(filepath.Join(mvnDir, "mvn"))

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	mockShell.On("Output").Return([]byte(tempDir), nil)

	actual, err := findMaven(mockProvider)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(mvnDir, "mvn"), actual)
	mockShell.AssertExpectations(t)
}

func TestFindMvnWrapperDisabled(t *testing.T) {
	initMvnTest(t)

	_ = os.Setenv("MENV_DISABLE_WRAPPER", "true")

	tempDir := t.TempDir()
	mvnDir := filepath.Join(tempDir, "maven", "3.9.6", "bin")
	_ = os.MkdirAll(mvnDir, 0755)
	os.Create(filepath.Join(mvnDir, "mvn"))
	wrapper := "mvnw"
	_, _ = os.Create(wrapper)

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	mockShell.On("Output").Return([]byte(tempDir), nil)

	actual, err := findMaven(mockProvider)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(mvnDir, "mvn"), actual)
	mockShell.AssertExpectations(t)
}

func TestFindMvnWithWrapper(t *testing.T) {
	_ = os.Unsetenv("MENV_DISABLE_WRAPPER")
	initMvnTest(t)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	wrapper := "mvnw"
	_, _ = os.Create(wrapper)

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	actual, err := findMaven(mockProvider)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(".", "mvnw"), actual)
	mockShell.AssertExpectations(t)
}

func TestFindMavenInvalidEnv(t *testing.T) {
	initMvnTest(t)

	_ = os.Setenv("MENV_DISABLE_WRAPPER", "invalid")

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	_, err := findMaven(mockProvider)
	assert.EqualError(t, err, "MENV_DISABLE_WRAPPER is not a boolean value")
	mockShell.AssertExpectations(t)
}

func TestSetMavenOptsNonExistentWithEnv(t *testing.T) {
	initMvnTest(t)

	expected := "-Xmx2g"
	_ = os.Setenv("MAVEN_OPTS", expected)

	actual := setMavenOpts("non_existent")
	assert.Equal(t, expected, actual)
}

func TestSetMavenOptsNonExistentWithoutEnv(t *testing.T) {
	initMvnTest(t)

	_ = os.Unsetenv("MAVEN_OPTS")
	actual := setMavenOpts("non_existent")
	assert.Empty(t, actual)
}

func TestSetMavenOpts(t *testing.T) {
	initMvnTest(t)

	_ = os.Unsetenv("MAVEN_OPTS")
	expected := "-Xmx2g"
	profile := "test"
	_ = profiles.Create(profile)
	_ = os.WriteFile(profiles.OptsFile(profile), []byte(expected), 0644)

	actual := setMavenOpts(profile)
	assert.Equal(t, expected, actual)
}

func TestSetMavenOptsEmpty(t *testing.T) {
	initMvnTest(t)

	_ = os.Unsetenv("MAVEN_OPTS")
	profile := "test"
	_ = profiles.Create(profile)
	_ = os.WriteFile(profiles.OptsFile(profile), []byte(""), 0644)

	actual := setMavenOpts(profile)
	assert.Empty(t, actual)
}

func TestExecMvnNoMvn(t *testing.T) {
	initMvnTest(t)

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	mockShell.On("Output").Return([]byte{}, errors.New("could not find maven in (home)brew cellar"))
	execMvn([]string{}, mockProvider)
	mockShell.AssertExpectations(t)
}

func TestExecMvn(t *testing.T) {
	initMvnTest(t)

	_ = profiles.Create("test")
	_ = profiles.Set("test")

	mockShell := MockShellCommand{
		Mock: &mock.Mock{},
	}

	mockProvider := func(string, ...string) profiles.ShellCommand {
		return &mockShell
	}

	tempDir := t.TempDir()
	mvnDir := filepath.Join(tempDir, "maven", "3.9.6", "bin")
	_ = os.MkdirAll(mvnDir, 0755)

	os.Create(filepath.Join(mvnDir, "mvn"))

	mockShell.On("Output").Return([]byte(tempDir), nil)
	mockShell.On("Stdin", os.Stdin).Return()
	mockShell.On("Stdout", os.Stdout).Return()
	mockShell.On("Stderr", os.Stderr).Return()
	mockShell.On("Run").Return(nil)
	execMvn([]string{}, mockProvider)
	mockShell.AssertExpectations(t)
}

func initMvnTest(t *testing.T) {
	_ = os.Unsetenv("MENV_DISABLE_WRAPPER")
	tempDir := t.TempDir()
	testConfig := config.Config{
		MenvRoot: tempDir,
		Editor:   "vi",
	}
	profiles.Init(testConfig)
	_ = os.Chdir(t.TempDir())
}
