# Developer Guide: nb-jrnl-ctl

**Date:** 2025-09-10
**Author:** Cline (AI Assistant)
**Version:** 1.0.0

## Project Overview

`nb-jrnl-ctl` is a Go-based command-line tool for reviewing NetBox journal entries. It provides a streamlined interface for listing and viewing journal entries through NetBox's API with enhanced formatting and usability.

## Architecture

The project follows a layered architecture pattern:

```
CLI Layer (cmd/) → Business Logic (internal/commands/) → Data Access (internal/client/) → Models (internal/models/) → Utilities (pkg/utils/)
```

### Directory Structure

```
├── cmd/                 # Main applications
│   └── nbjrnlctl/      # Primary CLI entry point
├── internal/           # Private application code
│   ├── client/         # NetBox API client
│   ├── commands/       # Command implementations (list.go, version.go)
│   └── models/         # Data structures
├── pkg/                # Public libraries
│   └── utils/          # Shared utilities
├── docs/               # Documentation
├── memory-bank/        # Memory bank documentation
├── go.mod             # Go module definition
└── go.sum             # Dependency checksums
```

## Getting Started

### Prerequisites

- Go 1.24.5 or higher
- Access to a NetBox instance
- NetBox API token with appropriate permissions

### Building from Source

Using the Justfile (recommended):
```bash
# Clone the repository
git clone https://github.com/jaepetto/nbjrnlctl.git
cd nbjrnlctl

# Build the binary with version info
just build

# Or install directly with version info
just install
```

Manual build:
```bash
go build -o nbjrnlctl ./cmd/nbjrnlctl
```

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests (when available)
5. Submit a pull request

## Code Structure

### Main Application (`cmd/nbjrnlctl/main.go`)

The entry point that initializes the Cobra command structure and registers the available commands (currently list and version).

### Commands (`internal/commands/`)

Each command is implemented as a separate file with a corresponding function that returns a `*cobra.Command`.

**Current commands:**
1. `list.go` - Lists journal entries with enhanced formatting
2. `version.go` - Displays version information

**Example pattern:**
```go
func ListCmd() *cobra.Command {
    var limit int

    cmd := &cobra.Command{
        Use:   "list [device-name]",
        Short: "List journal entries for a device",
        // ... command implementation
    }

    cmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of entries to display")
    // ... flag definitions

    return cmd
}
```

### Client (`internal/client/netbox.go`)

Handles all HTTP communication with the NetBox API:

- REST API calls for device lookup
- GraphQL queries for efficient journal entry fetching
- Error handling and response parsing
- Authentication management
- Smart kind extraction from display fields

### Models (`internal/models/journal.go`)

Defines the data structures used throughout the application with proper JSON tags for serialization.

### Utilities (`pkg/utils/config.go`)

Provides shared functionality:

- Environment variable-based configuration loading
- Hostname detection with fallback mechanisms

## Adding New Features

### Adding a New Command

1. Create a new file in `internal/commands/` (e.g., `search.go`)
2. Implement the command function following the existing pattern
3. Register the command in `cmd/nbjrnlctl/main.go`
4. Add appropriate flags and validation
5. Implement the business logic
6. Use the existing client for API calls

### Extending the Client

1. Add new methods to `internal/client/netbox.go`
2. Follow the existing error handling patterns
3. Use proper HTTP methods and status code checking
4. Parse responses into the appropriate models
5. Add unit tests (when implemented)

### Modifying Data Models

1. Update `internal/models/journal.go` with new fields
2. Ensure proper JSON tags are added
3. Update any affected client methods
4. Update command implementations if needed

## Testing

### Current State

The project currently lacks unit tests. This is a priority area for improvement.

### Recommended Testing Strategy

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test API interactions with a mock server
3. **End-to-End Tests**: Test complete command flows

### Test Structure

Recommended test directory structure:
```
├── internal/
│   ├── client/
│   │   └── netbox_test.go
│   ├── commands/
│   │   └── list_test.go
│   │   └── version_test.go
│   └── models/
│       └── journal_test.go
└── pkg/
    └── utils/
        └── config_test.go
```

## Error Handling

### Principles

1. **Fail Fast**: Exit early with clear error messages
2. **User-Friendly**: Provide actionable error messages
3. **Consistent**: Use the same error handling patterns throughout
4. **Graceful**: Handle network and API errors appropriately

### Implementation

```go
// Example error handling pattern
config, err := utils.LoadConfig()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
    os.Exit(1)
}
```

## Configuration Management

### Environment Variables

The application uses environment variables exclusively for configuration (no config files):

- `nbjrnlctl_base_url`: The base URL of your NetBox instance (e.g., `https://your-netbox-instance.com`)
- `nbjrnlctl_api_key`: Your NetBox API token with appropriate permissions

Example usage:
```bash
export nbjrnlctl_base_url="https://your-netbox-instance.com"
export nbjrnlctl_api_key="your-api-token"
nbjrnlctl list
```

## API Integration

### HTTP Client

Uses Go's standard `net/http` package with a 30-second timeout.

### Authentication

API tokens are passed in the `Authorization: Token {token}` header.

### Response Handling

- Check HTTP status codes before parsing responses
- Use proper JSON unmarshaling with struct tags
- Handle pagination for list endpoints
- Parse GraphQL responses with nested data structures

## Contributing

### Code Style

Follow Go idioms and conventions:

1. Use `gofmt` for code formatting
2. Write clear, descriptive comments
3. Use meaningful variable and function names
4. Keep functions focused and small
5. Handle errors explicitly

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make changes with clear commit messages
4. Add tests if applicable
5. Update documentation
6. Submit pull request with description

### Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters or less
- Reference issues and pull requests liberally

## Release Process

### Versioning

Follow Semantic Versioning (SemVer):

- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Changelog

Maintain `CHANGELOG.md` with all notable changes following the Keep a Changelog format.

### Release Steps

1. Update version in documentation
2. Update `CHANGELOG.md`
3. Create Git tag
4. Create GitHub release
5. Update documentation if needed

## Debugging

### Common Issues

1. **Authentication Errors**: Check API token and NetBox URL in environment variables
2. **Device Not Found**: Verify device exists in NetBox
3. **Network Issues**: Check connectivity to NetBox instance
4. **Permission Errors**: Verify API token permissions

### Logging

Currently, the application uses `fmt` for output. Future improvements could include structured logging.

## Performance Considerations

### Memory Usage

- Keep memory footprint low for CLI tool
- Avoid loading unnecessary data into memory
- Use streaming where possible for large responses

### Network Efficiency

- Minimize API calls per operation
- Use appropriate pagination limits
- Cache frequently accessed data (future enhancement)

## Security

### Best Practices

1. **Never log sensitive data** (API tokens, passwords)
2. **Validate all inputs** to prevent injection attacks
3. **Use HTTPS** for all API communications
4. **Store credentials securely** using environment variables only

### Credential Storage

Environment variables are the only method for storing credentials, eliminating file-based security risks.

## Future Development

### Roadmap

1. **Testing Infrastructure**: Add comprehensive test suite
2. **Advanced Features**: Additional command implementations
3. **Performance**: Caching, connection pooling optimization
4. **Usability**: Shell completion, better help system
5. **Reliability**: Retry logic, better error recovery

### Extension Points

1. **New Commands**: Additional journal entry operations
2. **Output Formats**: JSON, CSV export options
3. **Integration**: Webhook support, notification systems
4. **Automation**: Scripting support, batch processing

## Troubleshooting

### Common Error Messages

- **"Error loading config"**: Missing required environment variables
- **"Error finding device"**: Device doesn't exist in NetBox
- **"Unexpected status code"**: API returned unexpected response
- **"Invalid journal ID"**: Malformed or non-existent journal entry ID

### Debugging Steps

1. Verify environment variable configuration
2. Test NetBox API connectivity manually
3. Check NetBox logs for server-side errors
4. Enable verbose output (future feature)

## Dependencies

### Direct Dependencies

- `github.com/spf13/cobra` v1.8.0 - CLI framework
- `github.com/jedib0t/go-pretty/v6` v6.6.8 - Table formatting
- Go standard library

### Indirect Dependencies

Managed automatically by Go modules.

### Updating Dependencies

```bash
go get -u ./...
go mod tidy
```

## Build and Distribution

### Build Automation

The project uses a `justfile` for common development tasks. Here are the available commands:

```bash
just build          # Build the application with version info (statically linked)
just install        # Install the application globally with version info (statically linked)
just run            # Run the application
just run-with-args  # Run with custom arguments
just build-all      # Build for all platforms with version info (statically linked)
just clean          # Clean build artifacts
just test           # Run tests
just test-cover     # Run tests with coverage
just fmt            # Format code
just vet            # Vet code for issues
just tidy           # Tidy go modules
just deps           # Download dependencies
just update-deps    # Update dependencies
just audit          # Check for security vulnerabilities
just docs           # Generate documentation
just version        # Show Go version information
just help           # Show all available commands
```

### Version Management

The application includes built-in version tracking with detailed build metadata:

- **Semantic Version**: Hardcoded to "1.0.0" during build
- **Git Commit**: Embedded commit hash at build time (dynamically injected)
- **Build Date**: Timestamp when binary was compiled (dynamically injected)
- **Go Version**: Go compiler version used (dynamically detected at runtime)
- **Platform**: Target OS and architecture (dynamically detected at runtime)

Version information is injected at compile time using ldflags and can be accessed via:
- `nbjrnlctl version` command (shows detailed build information)
- `nbjrnlctl --version` flag (shows semantic version only)

Both the `build` and `install` commands use the same version injection mechanism to ensure consistency.

### Cross-Platform Builds

```bash
# Build for different platforms (with version info)
just build-all
```

### Release Artifacts

Single binary distribution with no external dependencies makes deployment simple. All binaries are statically linked for maximum portability.

## Documentation Updates

Keep documentation synchronized with code changes:

1. Update README.md for user-facing changes
2. Update developer guide for architectural changes
3. Update API documentation for endpoint changes
4. Update changelog for all notable changes
5. Update memory-bank documentation to reflect current state
