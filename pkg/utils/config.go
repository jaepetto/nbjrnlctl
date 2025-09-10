package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Config holds the application configuration
type Config struct {
	NetboxURL string `json:"netbox_url"`
	APIToken  string `json:"api_token"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load configuration from environment variables
	netboxURL := os.Getenv("nbjrnlctl_base_url")
	apiToken := os.Getenv("nbjrnlctl_api_key")

	// Check if required environment variables are set
	var missingVars []string
	if netboxURL == "" {
		missingVars = append(missingVars, "nbjrnlctl_base_url")
	}
	if apiToken == "" {
		missingVars = append(missingVars, "nbjrnlctl_api_key")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s. Please set these variables and try again", strings.Join(missingVars, ", "))
	}

	return &Config{
		NetboxURL: netboxURL,
		APIToken:  apiToken,
	}, nil
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
