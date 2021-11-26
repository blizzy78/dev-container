package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Cron *cronConfiguration `yaml:"cron,omitempty"`
}

type cronConfiguration struct {
	Jobs []*cronJobConfiguration `yaml:"jobs,omitempty"`
}

type cronJobConfiguration struct {
	Every   string   `yaml:"every"`
	Delay   string   `yaml:"delay,omitempty"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}

func loadConfig(path string) (*configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	var config configuration

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &config, nil
}
