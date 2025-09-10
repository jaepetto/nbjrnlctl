package main

import (
	"fmt"
	"os"

	"github.com/jaepetto/nbjrnlctl/internal/commands"
	"github.com/jaepetto/nbjrnlctl/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nbjrnlctl",
		Short: "Netbox Journal CLI Tool",
		Long: `A command-line interface for managing Netbox device journal entries.
Complete documentation is available at https://github.com/jaepetto/nbjrnlctl`,
		Version: version.GetShortVersion(),
	}

	// Add commands
	rootCmd.AddCommand(commands.ListCmd())
	rootCmd.AddCommand(commands.VersionCmd())

	// Add version flag shorthand
	rootCmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
