package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"os/exec"
	"runtime"
	"strings"
)

// Config holds the application configuration
type Config struct {
	NetboxURL string `json:"netbox_url"`
	APIToken  string `json:"api_token"`
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".nbjrnlctl")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Config file path
	configFile := filepath.Join(configDir, "config.json")

	// Check if config file exists, create if not
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return createDefaultConfig(configFile)
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configFile string) (*Config, error) {
	config := Config{
		NetboxURL: "https://netbox.example.com",
		APIToken:  "your-api-token",
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return nil, err
	}

	fmt.Printf("Created default config file at %s. Please edit it with your NetBox details.\n", configFile)
	return &config, nil
}

// GetHostname returns the hostname of the current machine
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		// If os.Hostname fails, try platform specific commands
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("hostname")
		default: // Unix-like systems
			cmd = exec.Command("hostname")
		}
		
		output, err := cmd.Output()
		if err != nil {
			return "unknown-host"
		}
		hostname = strings.TrimSpace(string(output))
	}
	return hostname
}
