package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show active profile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printActiveProfile(profiles.Active())
	},
}

func printActiveProfile(profile string, path string) {
	fmt.Println("Active profile: ")

	if profile == "" {
		fmt.Println("none")
		return
	}

	fmt.Printf("  %v (set by %v)\n", profile, path)
}

func init() {
	rootCmd.AddCommand(psCmd)
}
