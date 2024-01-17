package cmd

import (
	"github.com/spf13/cobra"
	"menv/profiles"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the local profile, if one is set in the current directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		profiles.Clear()
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
