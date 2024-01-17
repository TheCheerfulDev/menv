package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

// editCmd represents the edit command
var editoptsCmd = &cobra.Command{
	Use:   "editopts [profile]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Edit MAVEN_OPTS of the provided profile, or the active profile if none is provided, or prompt for a profile if none is active",
	Long: `With this command you can edit the MAVEN_OPTS of a profile. By default it will open the profile in vi.

If the profile has no MAVEN_OPTS set, maven will use the MAVEN_OPTS from the environment.

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

		err := profiles.EditOpts(profile)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(editoptsCmd)
}
