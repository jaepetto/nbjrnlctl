package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/jaepetto/nbjrnlctl/internal/client"
	"github.com/jaepetto/nbjrnlctl/pkg/utils"
	gptable "github.com/jedib0t/go-pretty/v6/table"
	gptext "github.com/jedib0t/go-pretty/v6/text"
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

			// Display entries in a table using go-pretty
			t := gptable.NewWriter()
			t.SetOutputMirror(os.Stdout)

			// Set table style
			style := gptable.StyleRounded
			style.Color.Header = gptext.Colors{gptext.Bold, gptext.FgHiBlue}
			t.SetStyle(style)

			// Add header
			t.AppendHeader(gptable.Row{"Created", "Created By", "Kind", "Comments"})

			// Add rows
			for _, entry := range entries {
				// Format date as MM/DD HH:MM
				formattedDate := entry.Created.Format("01/02 15:04")

				// Convert kind to emoji
				kindEmoji := getKindEmoji(entry.Kind)

				// Apply colors based on kind
				var dateColor gptext.Colors
				switch entry.Kind {
				case "Info":
					dateColor = gptext.Colors{gptext.FgBlue}
				case "Success":
					dateColor = gptext.Colors{gptext.FgGreen}
				case "Warning":
					dateColor = gptext.Colors{gptext.FgYellow}
				case "Danger":
					dateColor = gptext.Colors{gptext.FgRed}
				default:
					dateColor = gptext.Colors{gptext.FgWhite}
				}

				// Create styled cells
				styledDate := gptext.Colors(dateColor).Sprintf("%s", formattedDate)
				styledUser := gptext.Colors{gptext.FgCyan}.Sprintf("%s", entry.CreatedBy)
				styledKind := gptext.Colors(dateColor).Sprintf("%s", kindEmoji)

				t.AppendRow(gptable.Row{styledDate, styledUser, styledKind, entry.Comments})
			}

			t.Render()
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
