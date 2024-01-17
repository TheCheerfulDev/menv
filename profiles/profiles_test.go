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
	dir := "/tmp"
	editor := "editor"
	cfg = &config.Config{
		MenvRoot: dir,
		Editor:   editor,
	}
	Init(cfg)
	if cfg.MenvRoot != dir {
		t.Errorf("cfg.MenvRoot should be /tmp, got %v", cfg.MenvRoot)
	}
	if cfg.Editor != editor {
		t.Errorf("cfg.Editor should be non_default_editor, got %v", cfg.Editor)
	}
}

func TestCreateNew(t *testing.T) {
	initTest(t)
	err := Create("test")
	if err != nil {
		t.Errorf("Create should not return an error, got %v", err)
	}

	if _, err := os.Stat(cfg.MenvRoot + "/settings.xml.test"); os.IsNotExist(err) {
		t.Errorf("Create should create a profile")
	}

}

func TestCreateExisting(t *testing.T) {
	initTest(t)
	_ = Create("test")
	err := Create("test")
	if err == nil {
		t.Errorf("Create should return an error, got %v", err)
	}

}

func TestProfilesEmpty(t *testing.T) {
	initTest(t)
	profileList := Profiles()

	if len(profileList) != 0 {
		t.Errorf("Profiles should return an empty list, got %v", profileList)
	}
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

	if len(profileList) != 3 {
		t.Errorf("Profiles should return a slice with 3 items, got %v", profileList)
	}
}

func TestClear(t *testing.T) {
	initTest(t)
	profile := "test"
	dir := t.TempDir()
	_ = Create(profile)
	_ = os.Chdir(dir)
	_ = Set(profile)
	Clear(profile)
	if _, err := os.Stat(dir + "/settings.xml.test"); !os.IsNotExist(err) {
		t.Errorf("Clear should remove the profile")
	}
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
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Remove(%v) should return %v, got %v", test.profile, test.err, err)
		}
	}
}

func TestSetNonExisting(t *testing.T) {
	initTest(t)
	err := Set("non_existing")
	if err == nil {
		t.Errorf("Set should return an error, got %v", err)
	}
}

func TestSet(t *testing.T) {
	initTest(t)
	_ = Create("test")
	err := Set("test")
	if err != nil {
		t.Errorf("Set should not return an error, got %v", err)
	}
}

func TestExists(t *testing.T) {
	initTest(t)
	_ = Create("test")
	if !Exists("test") {
		t.Errorf("Exists should return true")
	}
	if Exists("non_existing") {
		t.Errorf("Exists should return false")
	}
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
		if profile != test.profile {
			t.Errorf("Active should return %v, got %v", test.profile, profile)
		}
		if path != test.want {
			t.Errorf("Active should return %v, got %v", test.want, path)
		}
	}
}

func TestActiveParentDir(t *testing.T) {
	initTest(t)
	dir, _ := os.Getwd()

	parent := filepath.Dir(dir)

	_ = Create("test")
	_ = os.Chdir(parent)
	_ = Set("test")
	_ = os.Chdir(dir)
	profile, path := Active()
	if profile != "test" {
		t.Errorf("Active should return test, got %v", profile)
	}
	if path != parent+"/"+profileFile {
		t.Errorf("Active should return %v, got %v", parent+"/"+profileFile, path)
	}
}

func TestExtractActiveVersionFromFile(t *testing.T) {
	initTest(t)
	_ = Create("test")
	_ = Set("test")
	profile, _ := Active()
	if profile != "test" {
		t.Errorf("Active should return test, got %v", profile)
	}
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
		got := MvnOptsExists(test.profile)
		if got != test.exists {
			t.Errorf("MvnOptsExists(%v) should return %v, got %v", test.profile, test.exists, got)
		}
	}
}

func TestMvnOpts(t *testing.T) {
	initTest(t)
	_ = Create("test")
	_ = os.WriteFile(cfg.MenvRoot+"/test.maven_opts", []byte("test"), 0644)
	opts := MvnOpts("test")
	if opts != "test" {
		t.Errorf("MvnOpts should return test, got %v", opts)
	}
}

func TestFile(t *testing.T) {
	initTest(t)
	file := File("test")
	if file != cfg.MenvRoot+"/settings.xml.test" {
		t.Errorf("File should return %v, got %v", cfg.MenvRoot+"/settings.xml.test", file)
	}
}

func TestOptsFile(t *testing.T) {
	initTest(t)
	file := OptsFile("test")
	if file != cfg.MenvRoot+"/test.maven_opts" {
		t.Errorf("File should return %v, got %v", cfg.MenvRoot+"/test.maven_opts", file)
	}
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
