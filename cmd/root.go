package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "menv",
	Short:   "Maven Environment Manager",
	Long:    `menv is a tool to manage maven profiles for a given folder and its children.`,
	Version: "0.9.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("menv version %s - © Mark Hendriks <thecheerfuldev>\n", rootCmd.Version))
}
