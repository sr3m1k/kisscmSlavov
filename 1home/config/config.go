package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ShellConfig struct {
	Username     string `yaml:"username"`
	ComputerName string `yaml:"computer_name"`
	VFSArchive   string `yaml:"vfs_archive"`
	LogFile      string `yaml:"log_file"`
}

func LoadConfig(path string) (*ShellConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ShellConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
