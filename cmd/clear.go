package cmd

import (
	"github.com/spf13/cobra"
	"menv/profiles"
	"os"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the local profile, if one is set in the current directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		profiles.Clear(dir)
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
