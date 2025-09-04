package commands

import (
	"fmt"
	"os"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// DeviceCmd returns the device command
func DeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Show information about the current device",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Get hostname
			hostname := utils.GetHostname()

			// Get config
			config, err := utils.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
				os.Exit(1)
			}

			// Create client
			nbClient := client.NewNetboxClient(config.NetboxURL, config.APIToken)

			// Look up device ID from name
			deviceID, err := nbClient.GetDeviceIDByName(hostname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding device '%s' in Netbox: %s\n", hostname, err)
				fmt.Fprintf(os.Stderr, "Please ensure this device exists in Netbox.\n")
				os.Exit(1)
			}

			fmt.Printf("Default Device Information:\n")
			fmt.Printf("  Hostname: %s\n", hostname)
			fmt.Printf("  Netbox Device ID: %d\n", deviceID)
			fmt.Printf("\nThis device will be used when no device name is specified for commands.\n")
		},
	}

	return cmd
}
