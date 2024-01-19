package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"menv/profiles"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// mvnCmd represents the mvn command
var mvnCmd = &cobra.Command{
	Use:                "mvn",
	Hidden:             true,
	DisableFlagParsing: true,
	Short:              "Execute a command with maven",
	Long:               `This command will execute a command with maven.`,
	Run: func(cmd *cobra.Command, args []string) {
		execMvn(args, profiles.ExecCmdProvider)
	},
}

func execMvn(args []string, shell func(string, ...string) profiles.ShellCommand) {
	mvnArgs := make([]string, 0)
	profile, _ := profiles.Active()
	_ = setMavenOpts(profile)
	if profiles.Exists(profile) {
		file := profiles.File(profile)
		mvnArgs = []string{"--settings", file, "--global-settings", file}
		mvnArgs = append(mvnArgs, args...)
	} else {
		mvnArgs = args
	}

	mvn, err := findMaven(shell)

	if err != nil {
		fmt.Println(err)
		return
	}
	cmd := shell(mvn, mvnArgs...)

	cmd.Stdin(os.Stdin)
	cmd.Stdout(os.Stdout)
	cmd.Stderr(os.Stderr)
	_ = cmd.Run()
}

func findMaven(shell func(string, ...string) profiles.ShellCommand) (string, error) {
	env, b := os.LookupEnv("MENV_DISABLE_WRAPPER")
	var disabled bool
	if b {
		parsed, err := strconv.ParseBool(env)
		if err != nil {
			return "", errors.New("MENV_DISABLE_WRAPPER is not a boolean value")
		}
		disabled = parsed
	}

	_, err := os.Stat("mvnw")
	wrapperFound := !os.IsNotExist(err)

	if wrapperFound && !disabled {
		return findMvnWrapper()
	}

	return findMvnInCellar(shell)
}

func findMvnWrapper() (string, error) {
	return filepath.Join(".", "mvnw"), nil
}

func findMvnInCellar(shell func(string, ...string) profiles.ShellCommand) (string, error) {
	cmd, _ := shell("brew", "--cellar").Output()
	cellar := string(cmd)
	cellar = strings.ReplaceAll(cellar, "\n", "")
	cellar = strings.ReplaceAll(cellar, "\r", "")
	cellar = filepath.Join(cellar, "maven")

	if _, err := os.Stat(cellar); os.IsNotExist(err) {
		return "", errors.New("could not find maven in (home)brew cellar")
	}

	mvn := ""

	_ = filepath.WalkDir(cellar, func(path string, d fs.DirEntry, err error) error {
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
			_ = os.Unsetenv("MAVEN_OPTS")
			return opts
		}
		_ = os.Setenv("MAVEN_OPTS", opts)
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
