package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type LayoutConfig struct {
	Name    string         `yaml:"name"`
	Windows []WindowConfig `yaml:"windows"`
}

type WindowConfig struct {
	ID        string `yaml:"id"`
	Title     string `yaml:"title"`
	Component string `yaml:"component"`
	X         int    `yaml:"x"`
	Y         int    `yaml:"y"`
	Width     int    `yaml:"width"`
	Height    int    `yaml:"height"`
}

func LoadLayout(path string) (*LayoutConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config LayoutConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
