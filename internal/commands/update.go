package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/internal/models"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// UpdateCmd returns the update command
func UpdateCmd() *cobra.Command {
	var comments string
	var kind string

	cmd := &cobra.Command{
		Use:   "update [journal-id]",
		Short: "Update an existing journal entry",
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

			// First get the current entry
			currentEntry, err := nbClient.GetJournalEntry(journalID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting journal entry: %s\n", err)
				os.Exit(1)
			}

			// Update fields if provided
			updatedEntry := models.JournalEntry{
				ID:       journalID,
				Comments: currentEntry.Comments,
				Kind:     currentEntry.Kind,
			}

			if comments != "" {
				updatedEntry.Comments = comments
			}

			if kind != "" {
				updatedEntry.Kind = kind
			}

			// Update journal entry
			result, err := nbClient.UpdateJournalEntry(journalID, updatedEntry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating journal entry: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Journal entry %d updated successfully\n", result.ID)
		},
	}

	cmd.Flags().StringVarP(&comments, "comments", "c", "", "New comments for the journal entry")
	cmd.Flags().StringVarP(&kind, "kind", "k", "", "New kind of journal entry (info, success, warning, danger)")

	return cmd
}
