package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/config"
	"menv/profiles"
	"os"
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

	},
}

func handleMavenPropertyAlreadySet(profile string) {
	// menv set, but not same profile
	if isMenvProperty() {
		// replace profile with current

	}
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
	return strings.Contains(workspace, "settings.xml."+profile)
}

func init() {
	rootCmd.AddCommand(ideaCmd)
}
