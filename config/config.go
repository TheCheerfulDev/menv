package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	MenvRoot string
	Editor   string
}

var cfg *Config

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		MenvRoot: filepath.Join(home, ".config", "menv"),
		Editor:   "vi",
	}
}

func Editor() string {
	editor, b := os.LookupEnv("MENV_EDITOR")
	if b {
		return editor
	}
	return cfg.Editor
}

func Set(config *Config) {
	cfg = config
}

func Get() *Config {
	return cfg
}

func Init() {
	err := os.MkdirAll(cfg.MenvRoot, 0755)
	if err != nil {
		fmt.Println("Could not create config dir")
		os.Exit(1)
	}

}
