package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

//go:embed config.yml
var embeddedConfig []byte

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

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Dir(cfgPath), 0o755)
		_ = os.WriteFile(cfgPath, embeddedConfig, 0o644)
	}

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Printf("❌ kunde inte läsa %s:\n%v\n", cfgPath, err)
		os.Exit(1)
	}

	aliasMap := make(map[string]string)
	for key, entry := range cfg.Links {
		aliasMap[key] = entry.URL
		for _, a := range entry.Aliases {
			aliasMap[a] = entry.URL
		}
	}

	if len(os.Args) < 2 || os.Args[1] == "--help" {
		fmt.Println("docs <språk>  – öppna officiell dokumentation i webbläsaren")
		fmt.Println("docs --help   – denna hjälp")
		fmt.Println("Tillgängliga språk:")
		keys := make([]string, 0, len(cfg.Links))
		for k := range cfg.Links {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Println("  ", k)
		}
		return
	}

	lang := os.Args[1]
	url, ok := aliasMap[lang]
	if !ok {
		fmt.Printf("🚫 okänt språk: %s\nKör  docs --help  för lista.\n", lang)
		os.Exit(1)
	}

	_ = exec.Command("xdg-open", url).Start()
	_ = exec.Command("i3-msg", "workspace", cfg.Workspace.Web).Run()
}
