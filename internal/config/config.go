package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CompanyName     string   `yaml:"company_name"`
	NoticeFormat    string   `yaml:"notice_format"`
	AutoFix         bool     `yaml:"auto_fix"`
	FilePatterns    []string `yaml:"file_patterns"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		CompanyName:  "Your Company",
		NoticeFormat: "Copyright (C) $year $company_name. All rights reserved.",
		AutoFix:      false,
		FilePatterns: []string{
			"*.go",
			"*.py",
			"*.js",
			"*.ts",
			"*.java",
			"*.cpp",
			"*.c",
			"*.h",
		},
		ExcludePatterns: []string{
			"vendor/",
			"node_modules/",
			".git/",
			"*.pb.go",
			"*_generated.go",
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	if configPath == "" {
		// Try to find config file automatically
		configPath = findConfigFile()
	}

	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	}

	return config, nil
}

// findConfigFile looks for config files in common locations
func findConfigFile() string {
	candidates := []string{
		".copyright.yaml",
		".copyright.yml",
		"copyright.yaml",
		"copyright.yml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

// loadFromFile loads config from YAML file
func loadFromFile(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// OverrideFromFlags allows CLI flags to override config values
func (c *Config) OverrideFromFlags(company, format string, autoFix bool) {
	if company != "" {
		c.CompanyName = company
	}
	if format != "" {
		c.NoticeFormat = format
	}
	if autoFix {
		c.AutoFix = autoFix
	}
}

// GenerateNotice creates a copyright notice using the template
func (c *Config) GenerateNotice() string {
	notice := c.NoticeFormat
	notice = strings.ReplaceAll(notice, "$company_name", c.CompanyName)
	notice = strings.ReplaceAll(notice, "$year", fmt.Sprintf("%d", time.Now().Year()))
	notice = strings.ReplaceAll(notice, "$current_year", fmt.Sprintf("%d", time.Now().Year()))
	return notice
}

// ShouldProcessFile checks if a file should be processed based on patterns
func (c *Config) ShouldProcessFile(filePath string) bool {
	// Check exclude patterns first
	for _, pattern := range c.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return false
		}
		// Also check if any part of the path matches
		if strings.Contains(filePath, strings.TrimSuffix(pattern, "/")) {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range c.FilePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
	}

	return false
}
