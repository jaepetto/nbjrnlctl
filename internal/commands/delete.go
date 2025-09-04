package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// DeleteCmd returns the delete command
func DeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete [journal-id]",
		Short: "Delete a journal entry",
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

			// If not forced, show the journal entry details and ask for confirmation
			if !force {
				entry, err := nbClient.GetJournalEntry(journalID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting journal entry: %s\n", err)
					os.Exit(1)
				}

				fmt.Printf("Are you sure you want to delete the following journal entry?\n")
				fmt.Printf("ID: %d\n", entry.ID)
				fmt.Printf("Created: %s\n", entry.Created.Format("2006-01-02 15:04:05"))
				fmt.Printf("Created By: %s\n", entry.CreatedBy)
				fmt.Printf("Kind: %s\n", entry.Kind)
				fmt.Printf("Comments: %s\n", entry.Comments)

				fmt.Print("Type 'yes' to confirm: ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" {
					fmt.Println("Deletion cancelled")
					return
				}
			}

			// Delete journal entry
			err = nbClient.DeleteJournalEntry(journalID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting journal entry: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Journal entry deleted successfully")
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete without confirmation")

	return cmd
}
