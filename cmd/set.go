package cmd

import (
	"fmt"
	"menv/profiles"
	"strconv"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:               "set [profile]",
	Args:              cobra.MaximumNArgs(1),
	Aliases:           []string{"profile"},
	ValidArgsFunction: profiles.CustomProfileCompletion,
	Short:             "Set given profile as active profile for this folder and children.",
	Long:              `Set given profile as active profile for this folder and children. If no profile is given, you will be prompted to select one`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			profile := PromptForProfile()
			if profile == "" {
				return
			}
			err := setProfile(profile)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
		profile := args[0]
		err := setProfile(profile)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func setProfile(profile string) error {
	err := profiles.Set(profile)
	if err != nil {
		return err
	}
	fmt.Printf("Set profile %v\n", profile)
	return nil
}

func PromptForProfile() string {
	profileList := profiles.Profiles()

	if len(profileList) == 0 {
		fmt.Println("No profiles found")
		return ""
	}

	var i int64 = 1

	for _, profile := range profileList {
		fmt.Printf("%v) %v\n", i, profile)
		i++
	}

	var choice string
	fmt.Print("Select profile: ")
	_, _ = fmt.Scanln(&choice)
	parseInt, err := strconv.ParseInt(choice, 10, 0)

	if err != nil || parseInt < 1 || parseInt > int64(len(profileList)) {
		fmt.Println("Invalid choice")
		return ""
	}

	return profileList[parseInt-1]

}

func init() {
	rootCmd.AddCommand(setCmd)
}
