package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// ReadCmd returns the read command
func ReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [journal-id]",
		Short: "Read a specific journal entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Parse journal ID
			journalID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid journal ID: %s\n", err)
				os.Exit(1)
			}

			// Get config
			config, err := utils.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
				os.Exit(1)
			}

			// Create client
			nbClient := client.NewNetboxClient(config.NetboxURL, config.APIToken)

			// Get journal entry
			entry, err := nbClient.GetJournalEntry(journalID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting journal entry: %s\n", err)
				os.Exit(1)
			}

			// Display journal entry
			fmt.Printf("ID: %d\n", entry.ID)
			fmt.Printf("Created: %s\n", entry.Created.Format("2006-01-02 15:04:05"))
			fmt.Printf("Created By: %s\n", entry.CreatedBy)
			fmt.Printf("Kind: %s\n", entry.Kind)
			fmt.Printf("Comments: %s\n", entry.Comments)
		},
	}

	return cmd
}
