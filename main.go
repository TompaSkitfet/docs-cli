package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type linkEntry struct {
	URL     string   `yaml:"url"`
	Aliases []string `yaml:"aliases"`
}

type config struct {
	Workspace struct {
		Wm  string `yaml:"wm"`
		Web string `yaml:"web"`
	} `yaml:"workspace"`
	Links map[string]linkEntry `yaml:"links"`
}

func defaultConfigPath() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdg, "docs-cli", "config.yml")
}

func loadConfig(path string) (*config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfgPath := defaultConfigPath()
	if len(os.Args) > 1 && os.Args[1] == "--config" {
		fmt.Println(cfgPath)
		return
	}

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Println("❌ kunde inte läsa %s:\n%v\n", cfgPath, err)
		os.Exit(1)
	}

	fmt.Println("Följande språk finns:")
	for l := range cfg.Links {
		fmt.Println("•", l)
	}
}
