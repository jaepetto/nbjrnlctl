# Developer Guide: nb-jrnl-ctl

**Date:** 2025-09-04
**Author:** Cline (AI Assistant)
**Version:** 1.0.0

## Project Overview

`nb-jrnl-ctl` is a Go-based command-line tool for managing NetBox journal entries. It provides a streamlined interface for creating, reading, updating, and deleting journal entries through NetBox's API.

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
│   ├── commands/       # Command implementations
│   └── models/         # Data structures
├── pkg/                # Public libraries
│   └── utils/          # Shared utilities
├── docs/               # Documentation
├── go.mod             # Go module definition
└── go.sum             # Dependency checksums
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Access to a NetBox instance
- NetBox API token with appropriate permissions

### Building from Source

```bash
# Clone the repository
git clone https://github.com/jaepetto/nbjrnlctl.git
cd nbjrnlctl

# Build the binary
go build -o nbjrnlctl ./cmd/nbjrnlctl

# Or install directly
go install github.com/jaepetto/nbjrnlctl/cmd/nbjrnlctl@latest
```

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests (when available)
5. Submit a pull request

## Code Structure

### Main Application (`cmd/nbjrnlctl/main.go`)

The entry point that initializes the Cobra command structure and registers all commands.

### Commands (`internal/commands/`)

Each command is implemented as a separate file with a corresponding function that returns a `*cobra.Command`.

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

- REST API calls for CRUD operations
- GraphQL queries for complex data fetching
- Error handling and response parsing
- Authentication management

### Models (`internal/models/journal.go`)

Defines the data structures used throughout the application with proper JSON tags for serialization.

### Utilities (`pkg/utils/config.go`)

Provides shared functionality:

- Configuration loading and saving
- Hostname detection
- File system operations

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

### Config File Location

- **Unix/Linux/macOS**: `~/.nbjrnlctl/config.json`
- **Windows**: `%USERPROFILE%\.nbjrnlctl\config.json`

### Config Structure

```json
{
  "netbox_url": "https://your-netbox-instance.com",
  "api_token": "your-api-token"
}
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

1. **Authentication Errors**: Check API token and NetBox URL in config
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
4. **Store credentials securely** with appropriate file permissions

### Credential Storage

Configuration files are stored with user-only read/write permissions (0600) to protect API tokens.

## Future Development

### Roadmap

1. **Testing Infrastructure**: Add comprehensive test suite
2. **Advanced Features**: Bulk operations, advanced filtering
3. **Performance**: Caching, connection pooling optimization
4. **Usability**: Shell completion, better help system
5. **Reliability**: Retry logic, better error recovery

### Extension Points

1. **New Object Types**: Support for other NetBox objects
2. **Output Formats**: JSON, CSV export options
3. **Integration**: Webhook support, notification systems
4. **Automation**: Scripting support, batch processing

## Troubleshooting

### Common Error Messages

- **"Error loading config"**: Config file corrupted or inaccessible
- **"Error finding device"**: Device doesn't exist in NetBox
- **"Unexpected status code"**: API returned unexpected response
- **"Invalid journal ID"**: Malformed or non-existent journal entry ID

### Debugging Steps

1. Verify configuration file contents
2. Test NetBox API connectivity manually
3. Check NetBox logs for server-side errors
4. Enable verbose output (future feature)

## Dependencies

### Direct Dependencies

- `github.com/spf13/cobra` - CLI framework
- Go standard library

### Indirect Dependencies

Managed automatically by Go modules.

### Updating Dependencies

```bash
go get -u github.com/spf13/cobra
go mod tidy
```

## Build and Distribution

### Cross-Platform Builds

```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o nbjrnlctl-linux ./cmd/nbjrnlctl
GOOS=windows GOARCH=amd64 go build -o nbjrnlctl-windows.exe ./cmd/nbjrnlctl
GOOS=darwin GOARCH=amd64 go build -o nbjrnlctl-macos ./cmd/nbjrnlctl
```

### Release Artifacts

Single binary distribution with no external dependencies makes deployment simple.

## Documentation Updates

Keep documentation synchronized with code changes:

1. Update README.md for user-facing changes
2. Update developer guide for architectural changes
3. Update API documentation for endpoint changes
4. Update changelog for all notable changes
