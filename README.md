# nbjrnlctl - Netbox Journal CLI Tool

A command-line interface for managing Netbox device journal entries.

[![Go Report Card](https://goreportcard.com/badge/github.com/jaepetto/nbjrnlctl)](https://goreportcard.com/report/github.com/jaepetto/nbjrnlctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- List journal entries for a device

## Installation

```bash
go install github.com/jaepetto/nbjrnlctl/cmd/nbjrnlctl@latest
```

Or clone and build:

```bash
git clone https://github.com/jaepetto/nbjrnlctl.git
cd nbjrnlctl
go build -o nbjrnlctl ./cmd/nbjrnlctl
```

## Configuration

The tool uses environment variables for configuration:

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

# For the current machine (uses local hostname)
nbjrnlctl list

# Limit the number of entries
nbjrnlctl list switch-01 --limit 10
```

### Check version information

```bash
# Display version information
nbjrnlctl version

# Display version using the built-in flag
nbjrnlctl --version
```

## Documentation

### User Documentation
- [README.md](README.md) - This file, containing user guide and installation instructions
- [CHANGELOG.md](CHANGELOG.md) - Version history and changes

### Developer Documentation
- [Developer Guide](docs/developer-guide.md) - Comprehensive guide for contributors
- [API Documentation](docs/api-documentation.md) - Detailed API endpoints and integration
- [Code Assessment](docs/code-assessment.md) - Technical analysis and code quality review

## Journal Entry Types

- `info` - General information (default)
- `success` - Success message
- `warning` - Warning message
