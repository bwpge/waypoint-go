package main

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"

	log "github.com/bwpge/systemdlog"
)

type yamlItem struct {
	URL   string   `json:"url"`
	Redir bool     `json:"redir"`
	Alias []string `json:"alias"`
}

type configItem struct {
	url   string
	redir bool
}

type config map[string]configItem

var confPaths = []string{
	"/etc/waypoint.yml",
	"/etc/waypoint.yaml",
	"~/.config/waypoint.yml",
	"~/.config/waypoint.yaml",
	"waypoint.yml",
	"waypoint.yaml",
}

func loadConfig() (config, error) {
	filePath := ""
	for _, f := range confPaths {
		abs, err := filepath.Abs(expandTilde(f))
		if err != nil {
			log.Warnf("could not determine absolute path for '%s'", f)
			continue
		}
		f = abs

		log.Debugf("looking for config: %s", f)
		if checkFile(f) == nil {
			filePath = f
			break
		}
	}

	result := make(config)
	if filePath == "" {
		log.Warn("no server config found")
		return result, nil
	}

	log.Infof("loading server config: %s", filePath)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("error reading config: %s", err)
	}

	c := map[string]yamlItem{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return result, fmt.Errorf("error deserializing config: %s", err)
	}

	for k, v := range c {
		item := configItem{url: v.URL, redir: v.Redir}

		if item.url == "" {
			return result, fmt.Errorf("error in config key '%s': a url is required", k)
		}

		for _, alias := range append(v.Alias, k) {
			redir := ""
			if item.redir {
				redir = " (redirect)"
			}
			log.Debugf("setting key: %s => %s%s", alias, item.url, redir)
			result[alias] = item
		}
	}

	return result, nil
}
