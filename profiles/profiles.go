package profiles

import (
	"errors"
	"fmt"
	"io"
	"menv/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	profileFile string = ".menv_profile"
	template           = `<?xml version="1.0" encoding="UTF-8"?>
<settings xmlns="http://maven.apache.org/SETTINGS/1.0.0"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.0.0 http://maven.apache.org/xsd/settings-1.0.0.xsd">
</settings>
`
)

var cfg config.Config

type ShellCommand interface {
	Run() error
	Stdin(io.Reader)
	Stdout(io.Writer)
	Stderr(io.Writer)
	Output() ([]byte, error)
}

type execShellCommand struct {
	*exec.Cmd
}

func (e execShellCommand) Run() error {
	return e.Cmd.Run()
}

func (e execShellCommand) Stdin(stdin io.Reader) {
	e.Cmd.Stdin = stdin
}

func (e execShellCommand) Stdout(stdout io.Writer) {
	e.Cmd.Stdout = stdout
}

func (e execShellCommand) Stderr(stderr io.Writer) {
	e.Cmd.Stderr = stderr
}

func (e execShellCommand) Output() ([]byte, error) {
	return e.Cmd.Output()
}

func Create(profile string) error {
	if Exists(profile) {
		return errors.New(fmt.Sprintf("profile %v already exists", profile))
	}

	path := cfg.MenvRoot + "/settings.xml." + profile
	_ = os.WriteFile(path, []byte(template), 0644)
	return nil
}

func Profiles() []string {
	dir, _ := os.ReadDir(cfg.MenvRoot)

	result := make([]string, 0)
	for _, file := range dir {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "settings.xml.") {
			result = append(result, strings.ReplaceAll(file.Name(), "settings.xml.", ""))
		}
	}
	return result
}

func Clear() {
	_ = os.Remove(profileFile)
}

func Remove(profile string) error {
	if !Exists(profile) {
		return errors.New(fmt.Sprintf("profile %v does not exist", profile))
	}

	_ = os.Remove(cfg.MenvRoot + "/settings.xml." + profile)
	_ = os.Remove(cfg.MenvRoot + "/" + profile + ".maven_opts")
	return nil
}

func Set(profile string) error {
	if !Exists(profile) {
		return errors.New("profile does not exist")
	}
	err := os.WriteFile(profileFile, []byte(profile+"\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Exists(profile string) bool {
	_, err := os.Stat(cfg.MenvRoot + "/settings.xml." + profile)
	return !os.IsNotExist(err)
}

func Active() (profile string, path string) {
	currentDirectory, _ := os.Getwd()

	for {
		if !strings.HasSuffix(currentDirectory, "/") {
			currentDirectory += "/"
		}

		profileFilePath := currentDirectory + profileFile
		if _, err := os.Stat(filepath.Join(currentDirectory, profileFile)); !os.IsNotExist(err) {
			return extractActiveVersionFromFile(profileFilePath), profileFilePath
		}

		if currentDirectory == "/" {
			return "", ""
		}

		currentDirectory = filepath.Clean(filepath.Join(currentDirectory, ".."))
	}

}

func extractActiveVersionFromFile(filePath string) (version string) {
	fileContent, _ := os.ReadFile(filePath)
	version = string(fileContent)
	version = removeNewLineFromString(version)
	return version
}

func removeNewLineFromString(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return input
}

func Init(config config.Config) {
	cfg = config
}

func genericEdit(profile string, shell func(string, ...string) ShellCommand, fileFn func(string) string) error {
	editor := config.Editor()

	cmd := shell(editor, fileFn(profile))

	cmd.Stdin(os.Stdin)
	cmd.Stdout(os.Stdout)
	cmd.Stderr(os.Stderr)
	_ = cmd.Run()
	return nil
}

func Edit(profile string, shell func(string, ...string) ShellCommand) error {
	if !Exists(profile) {
		return errors.New(fmt.Sprintf("profile %v does not exist", profile))
	}
	return genericEdit(profile, shell, File)
}

func EditOpts(profile string, shell func(string, ...string) ShellCommand) error {
	if !Exists(profile) {
		return errors.New(fmt.Sprintf("profile %v does not exist", profile))
	}
	return genericEdit(profile, shell, OptsFile)
}

func MvnOptsExists(profile string) bool {
	_, err := os.Stat(OptsFile(profile))
	return !os.IsNotExist(err)
}

func MvnOpts(profile string) string {
	data, _ := os.ReadFile(OptsFile(profile))
	opts := string(data)
	opts = strings.ReplaceAll(opts, "\n", "")
	opts = strings.ReplaceAll(opts, "\r", "")
	return opts
}
func File(profile string) string {
	return cfg.MenvRoot + "/settings.xml." + profile
}

func OptsFile(profile string) string {
	return cfg.MenvRoot + "/" + profile + ".maven_opts"
}

func ExecCmdProvider(command string, args ...string) ShellCommand {
	return execShellCommand{
		exec.Command(command, args...),
	}
}
