package cmd

import (
	"github.com/spf13/cobra"
	"menv/profiles"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Args:  cobra.NoArgs,
	Short: "Clears the active profile, if it is set in the current directory",
	Long: `This commands checks if a profile is set via a .menv_profile file and removes it if it is set.
This results in that the default profile is used, or the profile that is set in one of the parent directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles.Clear()
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
