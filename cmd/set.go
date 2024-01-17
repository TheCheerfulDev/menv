package cmd

import (
	"fmt"
	"menv/profiles"
	"strconv"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Args:  cobra.MaximumNArgs(1),
	Short: "Set given profile as active profile for this folder and children. If no profile is given, you will be prompted to select one",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			profile := PromptForProfile()
			if profile == "" {
				return
			}
			setProfile(profile)
			return
		}
		profile := args[0]
		setProfile(profile)
	},
}

func setProfile(profile string) {
	err := profiles.Set(profile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Set profile %v\n", profile)
	return
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
	fmt.Scanln(&choice)
	parseInt, err := strconv.ParseInt(choice, 10, 0)
	if err != nil || parseInt < 1 || parseInt > int64(len(profileList)) {
		fmt.Println("Invalid choice")
		return ""
	}

	return profileList[parseInt-1]

}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
