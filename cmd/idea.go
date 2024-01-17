package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
	"os"
)

// ideaCmd represents the idea command
var ideaCmd = &cobra.Command{
	Use:   "idea",
	Args:  cobra.NoArgs,
	Short: "Override IntelliJ IDEA maven settings.xml to the active profile one.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if isNotMavenproject() {
			fmt.Println("Not a maven project or maven project root")
			return
		}
		if profile, _ := profiles.Active(); profile == "" {
			fmt.Println("No active profile")
			return
		}
	},
}

func isNotMavenproject() bool {
	_, err := os.Stat("pom.xml")
	return os.IsNotExist(err)
}

func isNotIntellijProject() bool {
	_, err := os.Stat(".idea")
	return os.IsNotExist(err)
}

func init() {
	rootCmd.AddCommand(ideaCmd)
}
