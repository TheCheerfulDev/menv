package main

import (
	"menv/cmd"
	"menv/config"
	"menv/profiles"
)

func main() {
	config.Set(config.Default())
	_ = config.Init()
	profiles.Init(config.Get())
	cmd.Execute()
}
