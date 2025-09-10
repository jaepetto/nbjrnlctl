# nbjrnlctl - Netbox Journal CLI Tool

A command-line interface for reviewing Netbox device journal entries with enhanced formatting.

[![Go Report Card](https://goreportcard.com/badge/github.com/jaepetto/nbjrnlctl)](https://goreportcard.com/report/github.com/jaepetto/nbjrnlctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- List journal entries for devices with enhanced colored output and emojis
- Automatic hostname detection for local machine journal review
- Environment variable-based configuration for security
- Built-in version tracking with detailed build metadata
- Cross-platform compatibility (Linux, macOS, Windows)

## Installation

Using Go (recommended):
```bash
go install github.com/jaepetto/nbjrnlctl/cmd/nbjrnlctl@latest
```

Or clone and build manually:
```bash
git clone https://github.com/jaepetto/nbjrnlctl.git
cd nbjrnlctl
go build -o nbjrnlctl ./cmd/nbjrnlctl
```

Or use the provided Justfile for build automation:
```bash
just build
```

## Configuration

The tool uses environment variables for configuration (no config files):

```bash
export nbjrnlctl_base_url="https://your-netbox-instance.com"
export nbjrnlctl_api_key="your-api-token"
```

Or set them inline when running the command:

```bash
nbjrnlctl_base_url="https://your-netbox-instance.com" nbjrnlctl_api_key="your-api-token" nbjrnlctl list
```

## Usage

### List journal entries for a device

```bash
# For a specific device
nbjrnlctl list switch-01

# For the current machine (automatically detects local hostname)
nbjrnlctl list

# Limit the number of entries (show only recent entries)
nbjrnlctl list switch-01 --limit 10
```

### Check version information

```bash
# Display detailed version information including build metadata
nbjrnlctl version

# Display version using the built-in flag
nbjrnlctl --version
```

## Enhanced Output Features

The list command provides enhanced output formatting:
- Colored output based on journal entry type
- Emojis for quick visual identification (ℹ️ Info, ✅ Success, ⚠️ Warning, 🚨 Danger)
- Automatic text wrapping for better terminal compatibility
- Horizontal line separators for improved readability
- Sorted by creation date (newest first)
- 80-character width constraint for standard terminal compatibility

## Documentation

### User Documentation
- [README.md](README.md) - This file, containing user guide and installation instructions
- [CHANGELOG.md](CHANGELOG.md) - Version history and changes

### Developer Documentation
- [Developer Guide](docs/developer-guide.md) - Comprehensive guide for contributors
- [API Documentation](docs/api-documentation.md) - Detailed API endpoints and integration
- [Code Assessment](docs/code-assessment.md) - Technical analysis and code quality review

## Journal Entry Types

- `info` - General information (default) ℹ️
- `success` - Success message ✅
- `warning` - Warning message ⚠️
- `danger` - Danger/error message 🚨

## Build Automation

The project includes a Justfile for common development tasks:
```bash
just build          # Build with version info
just install        # Install globally with version info
just build-all      # Build for all platforms
just run            # Run the application
just run-with-args  # Run with custom arguments
just clean          # Clean build artifacts
```
