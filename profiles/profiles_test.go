package profiles

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"menv/config"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	dir := t.TempDir()
	editor := "editor"
	cfg = &config.Config{
		MenvRoot: dir,
		Editor:   editor,
	}
	Init(cfg)

	assert.Equal(t, dir, cfg.MenvRoot, "cfg.MenvRoot should be %v, got %v", dir, cfg.MenvRoot)
	assert.Equal(t, editor, cfg.Editor, "cfg.Editor should be %v, got %v", editor, cfg.Editor)
}

func TestCreateNew(t *testing.T) {
	initTest(t)
	err := Create("test")

	assert.NoErrorf(t, err, "Create should not return an error, got %v", err)
	assert.FileExists(t, cfg.MenvRoot+"/settings.xml.test", "Create should create a profile")
}

func TestCreateExisting(t *testing.T) {
	initTest(t)
	_ = Create("test")
	err := Create("test")

	assert.Error(t, err, "Creating a duplicate should return an error")
}

func TestProfilesEmpty(t *testing.T) {
	initTest(t)
	profileList := Profiles()

	assert.Empty(t, profileList, "Profiles should return an empty list")
}

func TestProfiles(t *testing.T) {
	tempDir := t.TempDir()
	cfg = &config.Config{
		MenvRoot: tempDir,
		Editor:   "vi",
	}

	Init(cfg)
	_ = Create("test")
	_ = Create("test2")
	_ = Create("test3")
	profileList := Profiles()

	assert.Len(t, profileList, 3, "Profiles should return a slice with 3 items")
}

func TestClear(t *testing.T) {
	initTest(t)
	dir := t.TempDir()
	profile := "test"
	_ = Create(profile)
	_ = os.Chdir(dir)
	_ = Set(profile)

	assert.FileExists(t, dir+"/.menv_profile", "Profile should be set")
	Clear()
	assert.NoFileExists(t, dir+"/.menv_profile", "Clear should remove the profile")

}

func TestRemove(t *testing.T) {
	tests := []struct {
		profile string
		err     error
	}{
		{"existing_profile", nil},
		{"non_existing", errors.New("profile non_existing does not exist")},
	}

	initTest(t)
	_ = Create("existing_profile")
	for _, test := range tests {
		err := Remove(test.profile)
		if test.err != nil {
			assert.EqualErrorf(t, err, test.err.Error(), "Create(%v) should return %v, got %v", test.profile, test.err, err)
		} else {
			assert.NoErrorf(t, err, "Create(%v) should not return an error, got %v", test.profile, err)
		}
	}
}

func TestSetNonExisting(t *testing.T) {
	initTest(t)
	err := Set("non_existing")
	assert.Error(t, err, "Set should return an error")
}

func TestSet(t *testing.T) {
	initTest(t)
	_ = Create("test")
	err := Set("test")
	assert.NoError(t, err, "Set should not return an error")
}

func TestExists(t *testing.T) {
	initTest(t)
	_ = Create("test")

	assert.True(t, Exists("test"), "Exists should return true")
	assert.False(t, Exists("non_existing"), "Exists should return false")
}

func TestActiveSameDir(t *testing.T) {
	initTest(t)
	dir, _ := os.Getwd()
	tests := []struct {
		profile string
		want    string
		fn      func(string) error
	}{
		{"", "", func(s string) error {
			return nil
		}},
		{"test", dir + "/" + profileFile, Set},
	}

	_ = Create("test")

	for _, test := range tests {
		_ = test.fn(test.profile)
		profile, path := Active()

		assert.Equalf(t, test.profile, profile, "Active should return %v, got %v", test.profile, profile)
		assert.Equalf(t, test.want, path, "Active should return %v, got %v", test.want, path)
	}
}

func TestActiveParentDir(t *testing.T) {
	initTest(t)
	dir, _ := os.Getwd()
	parent := filepath.Dir(dir)

	expectedProfile := "test"
	_ = Create(expectedProfile)
	_ = os.Chdir(parent)
	_ = Set(expectedProfile)
	_ = os.Chdir(dir)
	profile, path := Active()

	assert.Equalf(t, expectedProfile, profile, "Active should return %v, got %v", expectedProfile, profile)
	assert.Equalf(t, parent+"/"+profileFile, path, "Active should return %v, got %v", parent+"/"+profileFile, path)
}

func TestExtractActiveVersionFromFile(t *testing.T) {
	initTest(t)
	expected := "test"
	_ = Create(expected)
	_ = Set(expected)
	actual, _ := Active()

	assert.Equal(t, expected, actual, "Active should return %v, got %v", expected, actual)
}

func TestRemoveNewLineFromString(t *testing.T) {
	input := "test\n\r"
	actual := removeNewLineFromString(input)
	expected := "test"
	assert.Equal(t, expected, actual, "removeNewLineFromString should return %v, got %v", expected, actual)
}

func TestActiveNonExistent(t *testing.T) {
	initTest(t)
	profile, path := Active()
	assert.Empty(t, profile, "Active should return empty profile")
	assert.Empty(t, path, "Active should return empty path")
}

func TestEdit(t *testing.T) {

}

func TestEditNonExistent(t *testing.T) {
	initTest(t)
	actual := Edit("non_existent")
	expected := errors.New("profile non_existent does not exist")

	assert.Error(t, actual)
	assert.Equal(t, actual, expected)
}

func TestEditOpts(t *testing.T) {
	// TODO
}

func TestEditOptsNonExistent(t *testing.T) {
	initTest(t)
	actual := EditOpts("non_existent")
	expected := errors.New("profile non_existent does not exist").Error()

	assert.EqualError(t, actual, expected)
}

func TestMvnOptsExists(t *testing.T) {
	tests := []struct {
		profile string
		exists  bool
	}{
		{"existing_profile", true},
		{"non_existing", false},
	}
	initTest(t)
	_ = Create("existing_profile")
	_ = os.WriteFile(cfg.MenvRoot+"/existing_profile.maven_opts", []byte("test"), 0644)
	for _, test := range tests {
		actual := MvnOptsExists(test.profile)
		assert.Equalf(t, test.exists, actual, "MvnOptsExists(%v) should return %v, actual %v", test.profile, test.exists, actual)
	}
}

func TestMvnOpts(t *testing.T) {
	initTest(t)
	expected := "test"
	_ = Create(expected)
	_ = os.WriteFile(cfg.MenvRoot+"/test.maven_opts", []byte(expected), 0644)
	opts := MvnOpts(expected)
	assert.Equalf(t, expected, opts, "MvnOpts(%v) should return %v, got %v", expected, expected, opts)
}

func TestFile(t *testing.T) {
	initTest(t)
	actual := File("test")

	assert.Equalf(t, cfg.MenvRoot+"/settings.xml.test", actual, "File should return %v, got %v", cfg.MenvRoot+"/settings.xml.test", actual)
}

func TestOptsFile(t *testing.T) {
	initTest(t)
	actual := OptsFile("test")

	assert.Equalf(t, cfg.MenvRoot+"/test.maven_opts", actual, "OptsFile should return %v, got %v", cfg.MenvRoot+"/test.maven_opts", actual)
}

func initTest(t *testing.T) {
	tempDir := t.TempDir()
	cfg = &config.Config{
		MenvRoot: tempDir,
		Editor:   "vi",
	}
	Init(cfg)
	_ = os.Chdir(t.TempDir())
}
