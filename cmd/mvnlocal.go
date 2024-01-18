package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
	"os"
)

// mvnlocalCmd represents the mvnlocal command
var mvnlocalCmd = &cobra.Command{
	Use:   "mvnlocal",
	Args:  cobra.NoArgs,
	Short: "Override current maven project/.mvn folder to the active profile settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if IsNotMavenProject() {
			fmt.Println("Not a maven project or maven project root")
			return
		}

		profile, _ := profiles.Active()
		if profile == "" {
			fmt.Println("No active profile")
			return
		}

		createMavenDir()
		file := profiles.File(profile)
		writeMavenConfig(file)

		if profiles.MvnOptsExists(profile) {
			optsFile := profiles.OptsFile(profile)
			writeMavenOpts(optsFile)
		}

		fmt.Printf("Maven project .mvn folder set to profile %v settings\n", profile)

	},
}

func writeMavenOpts(file string) {
	readFile, _ := os.ReadFile(file)
	opts := string(readFile)
	_ = os.WriteFile(".mvn/jvm.config", []byte(opts), 0644)
}

func writeMavenConfig(file string) {
	_ = os.WriteFile(".mvn/maven.config", []byte("--settings "+file), 0644)
}

func createMavenDir() {
	if _, err := os.Stat(".mvn"); os.IsNotExist(err) {
		_ = os.Mkdir(".mvn", 0755)
	}
}

func init() {
	rootCmd.AddCommand(mvnlocalCmd)
}
