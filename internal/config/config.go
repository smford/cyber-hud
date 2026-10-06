package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration for cyber-hud.
type Config struct {
	Operator      OperatorConfig      `json:"operator" yaml:"operator"`
	Telemetry     TelemetryConfig     `json:"telemetry" yaml:"telemetry"`
	Settings      SettingsConfig      `json:"settings" yaml:"settings"`
	VFS           map[string]string   `json:"vfs" yaml:"vfs"`
	Repositories  []Repository        `json:"repositories" yaml:"repositories"`
	BuildLogs     []BuildLog          `json:"build_logs" yaml:"build_logs"`
	QuickChips    []QuickChip         `json:"quick_chips" yaml:"quick_chips"`
	CustomCommands map[string]Command `json:"custom_commands" yaml:"custom_commands"`
	Skills        []Skill             `json:"skills" yaml:"skills"`
	Benches       []Bench             `json:"benches" yaml:"benches"`
}

type OperatorConfig struct {
	Handle   string `json:"handle" yaml:"handle"`
	Host     string `json:"host" yaml:"host"`
	Title    string `json:"title" yaml:"title"`
	Subtitle string `json:"subtitle" yaml:"subtitle"`
	Bio      string `json:"bio" yaml:"bio"`
	SecNet   string `json:"sec_net" yaml:"sec_net"`
}

type TelemetryConfig struct {
	Kernel       string `json:"kernel" yaml:"kernel"`
	Cyberware    string `json:"cyberware" yaml:"cyberware"`
	SecurityICE  string `json:"security_ice" yaml:"security_ice"`
	Repositories string `json:"repositories" yaml:"repositories"`
}

type SettingsConfig struct {
	DefaultTheme string `json:"default_theme" yaml:"default_theme"` // "cyberpunk" or "altered-carbon"
	MouseEnabled bool   `json:"mouse_enabled" yaml:"mouse_enabled"`
	SkipBoot     bool   `json:"skip_boot" yaml:"skip_boot"`
}

type Repository struct {
	Name     string   `json:"name" yaml:"name"`
	Category string   `json:"category" yaml:"category"`
	Stars    int      `json:"stars" yaml:"stars"`
	Forks    int      `json:"forks" yaml:"forks"`
	Desc     string   `json:"desc" yaml:"desc"`
	Tech     []string `json:"tech" yaml:"tech"`
	URL      string   `json:"url" yaml:"url"`
	Demo     string   `json:"demo,omitempty" yaml:"demo,omitempty"`
}

type BOMItem struct {
	Item string `json:"item" yaml:"item"`
	Qty  string `json:"qty" yaml:"qty"`
	Role string `json:"role" yaml:"role"`
}

type PinoutItem struct {
	GPIO   string `json:"gpio" yaml:"gpio"`
	Device string `json:"device" yaml:"device"`
	Note   string `json:"note" yaml:"note"`
}

type BuildLog struct {
	ID        string       `json:"id" yaml:"id"`
	Category  string       `json:"category" yaml:"category"`
	Title     string       `json:"title" yaml:"title"`
	Date      string       `json:"date" yaml:"date"`
	Status    string       `json:"status" yaml:"status"`
	Summary   string       `json:"summary" yaml:"summary"`
	BOM       []BOMItem    `json:"bom" yaml:"bom"`
	Pinout    []PinoutItem `json:"pinout" yaml:"pinout"`
	Schematic string       `json:"schematic" yaml:"schematic"`
	Notes     string       `json:"notes" yaml:"notes"`
}

type QuickChip struct {
	Label   string `json:"label" yaml:"label"`
	Command string `json:"command" yaml:"command"`
}

type Command struct {
	Description string `json:"description" yaml:"description"`
	Output      string `json:"output" yaml:"output"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
	Executable  string `json:"executable,omitempty" yaml:"executable,omitempty"`
}

type Skill struct {
	Name       string `json:"name" yaml:"name"`
	Percentage int    `json:"percentage" yaml:"percentage"`
}

type Bench struct {
	ID    string   `json:"id" yaml:"id"`
	Title string   `json:"title" yaml:"title"`
	Items []string `json:"items" yaml:"items"`
}

// Load reads configuration from a YAML or JSON file. If path is empty, it returns the embedded default.
func Load(path string) (*Config, error) {
	if path == "" {
		cfg := DefaultConfig()
		return &cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := DefaultConfig()
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	default:
		// Attempt YAML first, then JSON
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			if jsonErr := json.Unmarshal(data, &cfg); jsonErr != nil {
				return nil, fmt.Errorf("failed to parse config as YAML (%v) or JSON (%v)", err, jsonErr)
			}
		}
	}

	return &cfg, nil
}

// Save writes the configuration to a file (in YAML or JSON format).
func Save(path string, cfg *Config) error {
	ext := strings.ToLower(filepath.Ext(path))
	var data []byte
	var err error

	if ext == ".json" {
		data, err = json.MarshalIndent(cfg, "", "  ")
	} else {
		data, err = yaml.Marshal(cfg)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
