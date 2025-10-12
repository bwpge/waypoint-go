package main

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"

	log "github.com/bwpge/systemdlog"
)

type yamlItem struct {
	Content string   `json:"content"`
	URL     string   `json:"url"`
	Redir   bool     `json:"redir"`
	Alias   []string `json:"alias"`
}

type cliOptions struct {
	Port uint   `json:"port"`
	Host string `json:"host"`
}

type yamlConfig struct {
	Options cliOptions          `json:"options"`
	Items   map[string]yamlItem `json:"items"`
}

var confPaths = []string{
	"/etc/waypoint/waypoint.yml",
	"/etc/waypoint/waypoint.yaml",
	"~/.config/waypoint.yml",
	"~/.config/waypoint.yaml",
	"waypoint.yml",
	"waypoint.yaml",
}

func loadConfig() (yamlConfig, error) {
	filePath := ""
	for _, f := range confPaths {
		abs, err := filepath.Abs(expandTilde(f))
		if err != nil {
			log.Warnf("could not determine absolute path for '%s'", f)
			continue
		}

		f = abs
		if checkFile(f) == nil {
			filePath = f
			break
		}
	}

	c := yamlConfig{}
	if filePath == "" {
		log.Warn("no server config found")
		return c, nil
	}

	log.Infof("loading server config: %s", filePath)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return c, fmt.Errorf("error reading config: %s", err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return c, fmt.Errorf("error deserializing config: %s", err)
	}

	if len(c.Items) == 0 {
		log.Warn("no items defined in server config")
	}

	return c, nil
}
