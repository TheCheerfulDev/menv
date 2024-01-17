package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"menv/profiles"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// mvnCmd represents the mvn command
var mvnCmd = &cobra.Command{
	Use:                "mvn",
	Hidden:             true,
	DisableFlagParsing: true,
	Short:              "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		execMvn(args)
	},
}

func execMvn(args []string) {
	mvnArgs := make([]string, 0)
	profile, _ := profiles.Active()
	opts := setMavenOpts(profile)
	if profiles.Exists(profile) {
		file := profiles.File(profile)
		mvnArgs = []string{"--settings", file, "--global-settings", file}
		mvnArgs = append(mvnArgs, args...)
	} else {
		mvnArgs = args
	}

	mvn, err := findMvnExecutable()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("[INFO] MAVEN_OPTS: %v\n", opts)
	cmd := exec.Command(mvn, mvnArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func findMvnExecutable() (string, error) {

	// look for mvn in homebrew
	cmd, _ := exec.Command("brew", "--cellar").Output()
	cellar := string(cmd)
	cellar = strings.ReplaceAll(cellar, "\n", "")
	cellar = strings.ReplaceAll(cellar, "\r", "")
	cellar = filepath.Join(cellar, "maven")

	if _, err := os.Stat(cellar); os.IsNotExist(err) {
		return "", errors.New("could not find maven in (home)brew cellar")
	}

	mvn := ""

	filepath.WalkDir(cellar, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == "mvn" && !strings.Contains(path, "libexec") {
			mvn = path
			return nil

		}
		return nil
	})

	if mvn == "" {
		return "", errors.New("could not find maven in (home)brew cellar")
	}

	return mvn, nil
}

func setMavenOpts(profile string) string {
	if profiles.Exists(profile) && profiles.MvnOptsExists(profile) {
		opts := profiles.MvnOpts(profile)
		if opts == "" {
			os.Unsetenv("MAVEN_OPTS")
			return opts
		}
		os.Setenv("MAVEN_OPTS", opts)
		return opts
	}

	env, exists := os.LookupEnv("MAVEN_OPTS")
	if exists {
		return env
	}

	return ""
}

func init() {
	rootCmd.AddCommand(mvnCmd)
}
