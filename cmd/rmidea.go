package cmd

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"menv/profiles"
	"os"
	"strings"
)

// rmideaCmd represents the rmidea command
var rmideaCmd = &cobra.Command{
	Use:   "rmidea",
	Short: "Try to remove the IntelliJ overridden settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if IsNotMavenProject() {
			fmt.Println("Not a maven project or maven project root")
			return
		}

		if IsNotIntellijProject() {
			fmt.Println("Not an IntelliJ project")
			return
		}

		profile, _ := profiles.Active()
		if profile == "" {
			fmt.Println("No active profile")
			return
		}

		if !workspaceExists() {
			fmt.Println("No .idea/workspace.xml found")
			return
		}

		if !isProfileUsedInWorkspace(profile) {
			fmt.Println("Active profile is not used in workspace.xml")
			return
		}

		err := removeProfileFromWorkspace(profile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Active profile removed from .idea/workspace.xml")
	},
}

func removeProfileFromWorkspace(profile string) error {
	workspace, _ := os.ReadFile(".idea/workspace.xml")

	doc := etree.NewDocument()
	err := doc.ReadFromBytes(workspace)
	if err != nil {
		return err
	}

	root := doc.SelectElement("project")
	optionElements := root.SelectElement("component").SelectElements("option")
	var mavenGeneralSettings *etree.Element

	for _, option := range optionElements {
		if option.SelectAttrValue("name", "") == "generalSettings" {
			mavenGeneralSettings = option.SelectElement("MavenGeneralSettings")
			break
		}
	}

	optionElements = mavenGeneralSettings.ChildElements()

	for _, option := range optionElements {
		name := option.SelectAttrValue("name", "")
		value := option.SelectAttrValue("value", "")
		if name == "userSettingsFile" && value == profiles.File(profile) {
			mavenGeneralSettings.RemoveChild(option)
		}
	}

	doc.Indent(2)
	err = doc.WriteToFile(".idea/workspace.xml")
	if err != nil {
		return err
	}
	return nil
}

func isProfileUsedInWorkspace(profile string) bool {
	template := "<option name=\"userSettingsFile\" value=\"{{profile}}\" />"
	template = strings.Replace(template, "{{profile}}", profiles.File(profile), 1)
	file, _ := os.ReadFile(".idea/workspace.xml")
	workspace := string(file)

	return strings.Contains(workspace, template)
}

func init() {
	rootCmd.AddCommand(rmideaCmd)
}
