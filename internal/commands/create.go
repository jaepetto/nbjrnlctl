package commands

import (
	"fmt"
	"os"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/internal/models"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// CreateCmd returns the create command
func CreateCmd() *cobra.Command {
	var comments string
	var kind string

	cmd := &cobra.Command{
		Use:   "create [device-name]",
		Short: "Create a journal entry for a device",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Determine device name - use local hostname if not provided
			deviceName := ""
			if len(args) == 1 {
				deviceName = args[0]
			} else {
				deviceName = utils.GetHostname()
				fmt.Printf("No device name provided. Using local hostname: %s\n", deviceName)
			}

			// Get config
			config, err := utils.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
				os.Exit(1)
			}

			// Create client
			nbClient := client.NewNetboxClient(config.NetboxURL, config.APIToken)

			// Look up device ID from name
			deviceID, err := nbClient.GetDeviceIDByName(deviceName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding device '%s': %s\n", deviceName, err)
				os.Exit(1)
			}

			// Prepare journal entry
			entry := models.JournalEntry{
				Comments: comments,
				Kind:     kind,
			}

			// Create journal entry
			result, err := nbClient.CreateJournalEntry(deviceID, entry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating journal entry: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Journal entry created with ID: %d\n", result.ID)
		},
	}

	cmd.Flags().StringVarP(&comments, "comments", "c", "", "Comments for the journal entry (required)")
	cmd.Flags().StringVarP(&kind, "kind", "k", "info", "Kind of journal entry (info, success, warning, danger)")
	cmd.MarkFlagRequired("comments")

	return cmd
}
