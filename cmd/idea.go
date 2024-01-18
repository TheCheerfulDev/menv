package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/config"
	"menv/profiles"
	"os"
	"regexp"
	"strings"
)

const workspaceTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="MavenImportPreferences">
    <option name="generalSettings">
      <MavenGeneralSettings>
        <option name="userSettingsFile" value="{{profile}}" />
      </MavenGeneralSettings>
    </option>
    <option name="enabledProfiles">
      <list>
        <option value="release" />
      </list>
    </option>
  </component>
</project>
`

const componentTemplate = `  <component name="MavenImportPreferences">
    <option name="generalSettings">
      <MavenGeneralSettings>
        <option name="userSettingsFile" value="{{menv_home}}/settings.xml.{{profile}}" />
      </MavenGeneralSettings>
    </option>
    <option name="enabledProfiles">
      <list>
        <option value="release" />
      </list>
    </option>
  </component>
</project>`

// ideaCmd represents the idea command
var ideaCmd = &cobra.Command{
	Use:   "idea",
	Args:  cobra.NoArgs,
	Short: "Override IntelliJ IDEA maven settings.xml to the active profile one.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if isNotMavenProject() {
			fmt.Println("Not a maven project or maven project root")
			return
		}

		if isNotIntellijProject() {
			fmt.Println("Not an IntelliJ project")
			return
		}

		profile, _ := profiles.Active()
		if profile == "" {
			fmt.Println("No active profile")
			return
		}

		if !workspaceExists() {
			// inject template
			_ = writeWorkspaceTemplate(profile)
			fmt.Printf("Maven settings set to profile %v\n", profile)
			return
		}

		if isProfileAlreadySet(profile) {
			fmt.Println("Profile already set")
			return
		}

		if isMavenPropertyAlreadySet() {
			handleMavenPropertyAlreadySet(profile)
			return
		}

		handleMavenPropertyNotSet(profile)

	},
}

func handleMavenPropertyNotSet(profile string) {
	// get current workspace.xml
	currFile, _ := os.ReadFile(".idea/workspace.xml")
	currWorkspace := string(currFile)

	// replace current workspace.xml with new profile
	exp := regexp.MustCompile("</project>")
	newWorkspace := exp.ReplaceAllString(currWorkspace, componentTemplate)
	newWorkspace = strings.ReplaceAll(newWorkspace, "{{menv_home}}", config.Get().MenvRoot)
	newWorkspace = strings.ReplaceAll(newWorkspace, "{{profile}}", profile)

	// write new workspace.xml
	_ = os.WriteFile(".idea/workspace.xml", []byte(newWorkspace), 0644)
}

func handleMavenPropertyAlreadySet(profile string) {
	// menv set, but not same profile
	if isMenvProperty() {
		replaceExistingMenvProfile(profile) // replace profile with current
		fmt.Printf("Maven settings set to profile %v\n", profile)
		return
	}

	instructions := `The IntelliJ workspace already has some custom settings.
Please override the maven 'User setting file:' property manually
in IntelliJ to the following value:

	{{menv_home}}/settings.xml.${profile}

`
	instructions = strings.ReplaceAll(instructions, "{{menv_home}}", config.Get().MenvRoot)
	fmt.Print(instructions)
}

func replaceExistingMenvProfile(profile string) {

	// get current workspace.xml
	currFile, _ := os.ReadFile(".idea/workspace.xml")
	currWorkspace := string(currFile)

	// replace current workspace.xml with new profile
	exp := regexp.MustCompile("settings[.]xml[.].[a-zA-Z0-9-_]+")
	newWorkspace := exp.ReplaceAllString(currWorkspace, "settings.xml."+profile)

	// write new workspace.xml
	_ = os.WriteFile(".idea/workspace.xml", []byte(newWorkspace), 0644)

}

func isMenvProperty() bool {
	file, _ := os.ReadFile(".idea/workspace.xml")
	workspace := string(file)
	return strings.Contains(workspace, "MavenImportPreferences") && strings.Contains(workspace, config.Get().MenvRoot)
}

func isMavenPropertyAlreadySet() bool {
	file, _ := os.ReadFile(".idea/workspace.xml")
	workspace := string(file)
	return strings.Contains(workspace, "MavenImportPreferences")
}

func isNotMavenProject() bool {
	_, err := os.Stat("pom.xml")
	return os.IsNotExist(err)
}

func isNotIntellijProject() bool {
	_, err := os.Stat(".idea")
	return os.IsNotExist(err)
}

func workspaceExists() bool {
	_, err := os.Stat(".idea/workspace.xml")
	return !os.IsNotExist(err)
}

func writeWorkspaceTemplate(profile string) error {
	profiles.File(profile)
	template := strings.ReplaceAll(workspaceTemplate, "{{profile}}", profiles.File(profile))
	return os.WriteFile(".idea/workspace.xml", []byte(template), 0644)
}

func isProfileAlreadySet(profile string) bool {
	file, _ := os.ReadFile(".idea/workspace.xml")
	workspace := string(file)
	return strings.Contains(workspace, config.Get().MenvRoot+"/settings.xml."+profile)
}

func init() {
	rootCmd.AddCommand(ideaCmd)
}
