// Package commands provides CLI command implementations.
// Date: 2025-09-10
// Description: Version command implementation that displays application version information.

package commands

import (
	"fmt"

	"github.com/jaepetto/nbjrnlctl/internal/version"
	"github.com/spf13/cobra"
)

// VersionCmd returns the version command that displays application version information.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display detailed version information including build metadata.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.GetVersion())
		},
	}
}
