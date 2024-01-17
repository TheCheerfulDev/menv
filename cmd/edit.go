package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [profile]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Edit the provided profile, or the active profile if none is provided, or prompt for a profile if none is active",
	Long: `With this command you can edit a profile. By default it will open the profile in vi.

You can change the editor by setting the MENV_EDITOR environment variable.

Example:
export MENV_EDITOR=nano`,
	Run: func(cmd *cobra.Command, args []string) {

		var profile string
		if len(args) == 0 {
			profile, _ = profiles.Active()
		} else {
			profile = args[0]
		}

		if profile == "" {
			profile = PromptForProfile()
		}

		if profile == "" {
			return
		}

		err := profiles.Edit(profile, profiles.ExecCmdProvider)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
