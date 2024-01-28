package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"menv/profiles"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:               "rm [profile]",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: profiles.CustomProfileCompletion,
	Short:             "Removes the provided profile",
	Long:              `This command will remove the provided profile. This includes the MAVEN_OPTS`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		profile := args[0]
		err := profiles.Remove(profile)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Removed profile %v\n", profile)

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
