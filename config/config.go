package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	MenvRoot string
	Editor   string
	Verbose  bool
}

var cfg Config

func Default() Config {
	home, _ := os.UserHomeDir()
	return Config{
		MenvRoot: filepath.Join(home, ".config", "menv"),
		Editor:   "vi",
		Verbose:  false,
	}
}

func Editor() string {
	editor, b := os.LookupEnv("MENV_EDITOR")
	if b {
		return editor
	}
	return cfg.Editor
}

func Verbose() bool {
	verbose, b := os.LookupEnv("MENV_VERBOSE")
	if b {
		parseBool, err := strconv.ParseBool(verbose)
		if err != nil {
			return false
		}
		return parseBool
	}
	return cfg.Verbose
}

func Set(config Config) {
	cfg = config
}

func Get() Config {
	return cfg
}

func Init() error {
	return os.MkdirAll(cfg.MenvRoot, 0755)
}
