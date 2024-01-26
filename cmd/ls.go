package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available profiles",
	Long:  `This command lists all available profiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		printProfiles(profiles.Profiles())
	},
}

func printProfiles(profileList []string) {
	if len(profileList) == 0 {
		fmt.Println("No profiles found")
		return
	}

	active, _ := profiles.Active()

	fmt.Println("Available profiles:")
	for _, profile := range profileList {
		if profile == active {
			fmt.Print("* ")
		} else {
			fmt.Print("  ")
		}
		fmt.Println(profile)
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
