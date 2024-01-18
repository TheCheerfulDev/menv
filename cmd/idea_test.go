package cmd

import (
	"github.com/stretchr/testify/assert"
	"menv/config"
	"menv/profiles"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsNotMavenProject(t *testing.T) {
	_ = os.Chdir(t.TempDir())
	assert.True(t, isNotMavenProject())
	_ = os.WriteFile("pom.xml", []byte("test"), 0644)
	assert.False(t, isNotMavenProject())
}

func TestIsNotIntellijProject(t *testing.T) {
	_ = os.Chdir(t.TempDir())
	assert.True(t, isNotIntellijProject())
	_ = os.Mkdir(".idea", 0755)
	assert.False(t, isNotIntellijProject())
}

func TestWorkspaceExists(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)

	assert.False(t, workspaceExists(), "workspace should not exist")
	_ = os.Mkdir(filepath.Join(tempDir, ".idea"), 0755)
	assert.False(t, workspaceExists(), "workspace should not exist, because .idea is empty")
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte("test"), 0644)
	assert.True(t, workspaceExists(), "workspace should exist")
}

func TestWriteWorkspaceTemplate(t *testing.T) {
	testCfg := config.Config{
		MenvRoot: t.TempDir(),
		Editor:   "vi",
	}
	profiles.Init(testCfg)
	_ = profiles.Create("test")
	_ = profiles.Set("test")
	profile, _ := profiles.Active()

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)

	assert.False(t, workspaceExists(), "workspace should not exist")
	assert.NoError(t, writeWorkspaceTemplate(profile), "workspace should be created")
	assert.True(t, workspaceExists(), "workspace should exist")

	actual, _ := os.ReadFile(filepath.Join(tempDir, ".idea", "workspace.xml"))
	expected := strings.Replace(workspaceTemplate, "{{profile}}", profiles.File(profile), 1)

	assert.Equal(t, expected, string(actual), "workspace template should be equal")
}

func TestIsProfileAlreadySet(t *testing.T) {
	template := `<component name="MavenImportPreferences">
    <option name="generalSettings">
      <MavenGeneralSettings>
        <option name="userSettingsFile" value="$USER_HOME$/.config/menv/settings.xml.test" />
      </MavenGeneralSettings>
    </option>
    <option name="enabledProfiles">
      <list>
        <option value="release" />
      </list>
    </option>
  </component>`
	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(template), 0644)

	assert.True(t, isProfileAlreadySet("test"), "profile should be set")
	assert.False(t, isProfileAlreadySet("non_existent"), "profile should not be set")
}

func TestIsMavenPropertyAlreadySet(t *testing.T) {
	template := `<component name="MavenImportPreferences">
	<option name="generalSettings">
	  <MavenGeneralSettings>
		<option name="userSettingsFile" value="$USER_HOME$/.config/menv/settings.xml.test" />
	  </MavenGeneralSettings>
	</option>
	<option name="enabledProfiles">
	  <list>
		<option value="release" />
	  </list>
	</option>
  </component>`
	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte("template"), 0644)
	assert.False(t, isMavenPropertyAlreadySet(), "maven property should not be set")
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(template), 0644)

	assert.True(t, isMavenPropertyAlreadySet(), "maven property should be set")
}

func TestIsMenvProperty(t *testing.T) {

	configDir := t.TempDir()

	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)

	_ = profiles.Create("test")

	template := strings.Replace(workspaceTemplate, "{{profile}}", profiles.File("test"), 1)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte("template"), 0644)
	assert.False(t, isMenvProperty(), "menv property should not be set")
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(template), 0644)

	assert.True(t, isMenvProperty(), "menv property should be set")
}

func TestHandleMavenPropertyAlreadySet(t *testing.T) {
	template := `<component name="MavenImportPreferences">
	<option name="generalSettings">
	  <MavenGeneralSettings>
		<option name="userSettingsFile" value="/other/path/.config/menv/settings.xml.test" />
	  </MavenGeneralSettings>
	</option>
	<option name="enabledProfiles">
	  <list>
		<option value="release" />
	  </list>
	</option>
  </component>`

	configDir := t.TempDir()

	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(template), 0644)

	handleMavenPropertyAlreadySet("new_profile")

	file, _ := os.ReadFile(filepath.Join(tempDir, ".idea", "workspace.xml"))
	actual := string(file)
	expected := template

	assert.Equal(t, expected, actual, "workspace template should be equal")
}

func TestHandleMavenPropertyAlreadySetByMenv(t *testing.T) {
	template := `<component name="MavenImportPreferences">
	<option name="generalSettings">
	  <MavenGeneralSettings>
		<option name="userSettingsFile" value="{{menv_home}}/settings.xml.test" />
	  </MavenGeneralSettings>
	</option>
	<option name="enabledProfiles">
	  <list>
		<option value="release" />
	  </list>
	</option>
  </component>`
	configDir := t.TempDir()
	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	input := strings.Replace(template, "{{menv_home}}", configDir, 1)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(input), 0644)

	handleMavenPropertyAlreadySet("new_profile")

	file, _ := os.ReadFile(filepath.Join(tempDir, ".idea", "workspace.xml"))
	actual := string(file)
	expected := strings.Replace(template, "{{menv_home}}", configDir, 1)
	expected = strings.Replace(expected, "settings.xml.test", "settings.xml.new_profile", 1)
	assert.Equalf(t, expected, actual, "workspace template should be equal")

}

func TestReplaceExistingMenvProperty(t *testing.T) {
	configDir := t.TempDir()

	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)

	_ = profiles.Create("test")

	template := strings.Replace(workspaceTemplate, "{{profile}}", profiles.File("test"), 1)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(template), 0644)

	replaceExistingMenvProfile("new_profile")

	actual, _ := os.ReadFile(filepath.Join(tempDir, ".idea", "workspace.xml"))
	expected := strings.Replace(template, "settings.xml.test", "settings.xml.new_profile", 1)

	assert.Equal(t, expected, string(actual), "workspace template should be equal")
}

func TestHandleMavenPropertyNotSet(t *testing.T) {
	emptyTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
</project>
`
	configDir := t.TempDir()

	testCfg := config.Config{
		MenvRoot: configDir,
		Editor:   "vi",
	}
	config.Set(testCfg)
	profiles.Init(testCfg)

	profile := "test"
	_ = profiles.Create(profile)

	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)
	_ = os.Mkdir(".idea", 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".idea", "workspace.xml"), []byte(emptyTemplate), 0644)

	handleMavenPropertyNotSet(profile)

	actual, _ := os.ReadFile(filepath.Join(tempDir, ".idea", "workspace.xml"))
	expected := strings.Replace(workspaceTemplate, "{{profile}}", profiles.File(profile), 1)

	assert.Equal(t, expected, string(actual), "workspace emptyTemplate should be equal")
}
