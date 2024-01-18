package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
	"os"
)

// newCmd represents the add command
var newCmd = &cobra.Command{
	Use:     "new [profile]",
	Aliases: []string{"add"},
	Args:    cobra.ExactArgs(1),
	Short:   "Create a new profile",
	Long: `With this command you can create a new profile.
The profile name must be unique, cannot be empty and cannot contain spaces.

The following characters are allowed (not including the comma's): a-z, A-Z, 0-9, -, and _
`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		err := profiles.Create(profile)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Created profile %v\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
