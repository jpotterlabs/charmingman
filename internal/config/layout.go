package config

import (
	"fmt"
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

	if err := layout.Validate(); err != nil {
		return nil, err
	}

	return &layout, nil
}

// Validate checks the layout for common configuration errors.
func (l *Layout) Validate() error {
	// Check for empty layout
	if len(l.Windows) == 0 {
		return fmt.Errorf("layout validation failed: no windows defined")
	}

	// Check for duplicate window IDs
	idMap := make(map[string]bool)
	for i, w := range l.Windows {
		if w.ID == "" {
			return fmt.Errorf("layout validation failed: window at index %d has empty ID", i)
		}
		if idMap[w.ID] {
			return fmt.Errorf("layout validation failed: duplicate window ID %q", w.ID)
		}
		idMap[w.ID] = true

		// Check for missing or empty Type field
		if w.Type == "" {
			return fmt.Errorf("layout validation failed: window %q has empty Type field", w.ID)
		}

		// Check for invalid dimensions
		if w.Width <= 0 {
			return fmt.Errorf("layout validation failed: window %q has invalid width %d (must be > 0)", w.ID, w.Width)
		}
		if w.Height <= 0 {
			return fmt.Errorf("layout validation failed: window %q has invalid height %d (must be > 0)", w.ID, w.Height)
		}
	}

	return nil
}