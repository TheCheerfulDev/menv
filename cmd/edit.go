package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [profile]",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `With this command you can edit a profile. By default it will open the profile in vi.

You can change the editor by setting the MENV_EDITOR environment variable.

Example:
export MENV_EDITOR=nano`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]

		err := profiles.Edit(profile)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
