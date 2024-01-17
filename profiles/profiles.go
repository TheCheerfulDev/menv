package profiles

import (
	"errors"
	"menv/config"
	"os"
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

var cfg *config.Config

func Add(name string) error {
	if Exists(name) {
		return errors.New("profile already exists")
	}

	path := cfg.MenvRoot + "/settings.xml." + name
	err := os.WriteFile(path, []byte(template), 0644)

	if err != nil {
		return err
	}
	return nil
}

func Profiles() []string {
	dir, err := os.ReadDir(cfg.MenvRoot)

	if err != nil {
		return []string{}
	}

	result := make([]string, 0)
	for _, file := range dir {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "settings.xml.") {
			result = append(result, strings.ReplaceAll(file.Name(), "settings.xml.", ""))
		}
	}
	return result
}

func Clear(path string) {
	os.Remove(path + "/" + profileFile)
}

func Delete(profile string) error {
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

func Active() (string, string) {
	currentDirectory, err := os.Getwd()

	if err != nil {
		return "", ""
	}

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

func Init(config *config.Config) {
	cfg = config
}
