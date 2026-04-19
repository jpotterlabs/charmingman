package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Layout struct {
	Windows []WindowConfig `yaml:"windows"`
}

type WindowConfig struct {
	ID      string            `yaml:"id"`
	Title   string            `yaml:"title"`
	Type    string            `yaml:"type"` // chat, document, tool
	X       int               `yaml:"x"`
	Y       int               `yaml:"y"`
	Width   int               `yaml:"width"`
	Height  int               `yaml:"height"`
	Focused bool              `yaml:"focused"`
	Config  ComponentConfig   `yaml:"config"`
}

type ComponentConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	Persona  string `yaml:"persona"`
	UseRAG   bool   `yaml:"use_rag"`
	Content  string `yaml:"content"` // For document type
}

// LoadLayout loads the layout configuration from a YAML file.
func LoadLayout(path string) (*Layout, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var layout Layout
	if err := yaml.Unmarshal(data, &layout); err != nil {
		return nil, err
	}

	return &layout, nil
}
