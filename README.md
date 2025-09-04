# nbjrnlctl - Netbox Journal CLI Tool

A command-line interface for managing Netbox device journal entries.

[![Go Report Card](https://goreportcard.com/badge/github.com/jaepetto/nbjrnlctl)](https://goreportcard.com/report/github.com/jaepetto/nbjrnlctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Create journal entries for devices
- Read journal entries
- Update existing journal entries
- Delete journal entries
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

On first run, the tool will create a config file at `~/.nbjrnlctl/config.json`. Edit this file to add your Netbox URL and API token:

```json
{
  "netbox_url": "https://your-netbox-instance.com",
  "api_token": "your-api-token"
}
```

## Usage

### Device-based Commands

These commands work with device names and automatically use your local machine's hostname if no device name is provided.

#### Check current device

```bash
# Show information about the current default device (local hostname)
nbjrnlctl device
```

#### Create a journal entry

```bash
# For a specific device
nbjrnlctl create switch-01 --comments "Replaced failed power supply" --kind "info"

# For the current machine (uses local hostname)
nbjrnlctl create --comments "Scheduled maintenance complete" --kind "success"
```

#### List journal entries for a device

```bash
# For a specific device
nbjrnlctl list switch-01

# For the current machine (uses local hostname)
nbjrnlctl list

# Limit the number of entries
nbjrnlctl list switch-01 --limit 10
```

### Journal Entry Commands

These commands work directly with journal entry IDs (which you can find using the `list` command).

#### Read a journal entry

```bash
nbjrnlctl read 456
```

#### Update a journal entry

```bash
nbjrnlctl update 456 --comments "Updated comment" --kind "success"
```

#### Delete a journal entry

```bash
# With confirmation prompt
nbjrnlctl delete 456

# With force option (no confirmation)
nbjrnlctl delete 456 --force
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
