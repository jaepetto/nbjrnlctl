package commands

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	"github.com/spf13/cobra"
)

// ListCmd returns the list command
func ListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list [device-name]",
		Short: "List journal entries for a device",
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

			// Get journal entries
			entries, err := nbClient.ListJournalEntries(deviceID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing journal entries: %s\n", err)
				os.Exit(1)
			}

			if len(entries) == 0 {
				fmt.Println("No journal entries found for this device")
				return
			}

			// Sort entries by creation date (newest first)
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].Created.After(entries[j].Created)
			})

			// Limit the number of entries if specified
			if limit > 0 && limit < len(entries) {
				entries = entries[:limit]
			}

			// Display entries in a table with enhanced formatting
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "CREATED\tCREATED BY\tKIND\tCOMMENTS")

			for _, entry := range entries {
				// Format date as MM/DD HH:MM
				formattedDate := entry.Created.Format("01/02 15:04")

				// Convert kind to emoji
				kindEmoji := getKindEmoji(entry.Kind)

				// Apply colors
				dateColor := getColorForKind(entry.Kind)
				coloredDate := fmt.Sprintf("%s%s\x1b[0m", dateColor, formattedDate)
				coloredUser := fmt.Sprintf("\x1b[36m%s\x1b[0m", entry.CreatedBy)
				coloredKind := fmt.Sprintf("%s%s\x1b[0m", dateColor, kindEmoji)

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					coloredDate,
					coloredUser,
					coloredKind,
					entry.Comments,
				)
			}
			w.Flush()
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of entries to display")

	return cmd
}

// getKindEmoji converts journal entry kinds to emojis
func getKindEmoji(kind string) string {
	switch kind {
	case "Info":
		return "ℹ️"
	case "Success":
		return "✅"
	case "Warning":
		return "⚠️"
	case "Danger":
		return "🚨"
	default:
		return "📝"
	}
}

// getColorForKind returns ANSI color codes based on entry kind
func getColorForKind(kind string) string {
	switch kind {
	case "Info":
		return "\x1b[34m" // Blue
	case "Success":
		return "\x1b[32m" // Green
	case "Warning":
		return "\x1b[33m" // Yellow
	case "Danger":
		return "\x1b[31m" // Red
	default:
		return "\x1b[37m" // White
	}
}
