# Technical Context: nb-jrnl-ctl

## Technology Stack
- **Language**: Go (Golang) 1.24.5
- **Build System**: Go modules
- **CLI Framework**: Cobra v1.8.0
- **HTTP Client**: Standard library net/http
- **Table Formatting**: go-pretty v6.6.8
- **Configuration**: Environment variables only
- **Testing**: Go testing package
- **Build Automation**: Justfile

## Development Environment
- Go 1.24.5 (based on go.mod)
- Cross-platform development (macOS, Linux, Windows)
- Standard Go workspace structure
- Git for version control

## Project Structure
```
├── cmd/                 # Main applications
│   └── nbjrnlctl/      # Primary CLI entry point
├── internal/           # Private application code
│   ├── client/         # NetBox API client
│   ├── commands/       # Command implementations (only list.go and version.go)
│   └── models/         # Data structures
├── pkg/                # Public libraries
│   └── utils/          # Shared utilities
├── docs/               # Documentation files
├── memory-bank/        # Memory bank documentation
├── go.mod             # Go module definition
├── go.sum             # Dependency checksums
├── justfile           # Build automation
├── README.md          # Project documentation
└── CHANGELOG.md       # Version history
```

## Dependencies
- `github.com/spf13/cobra` v1.8.0 - CLI framework
- `github.com/jedib0t/go-pretty/v6` v6.6.8 - Table formatting
- Standard Go library for core functionality
- NetBox API v2.0+ compatibility

## Build and Deployment
- `just build` for compilation with version injection
- Single binary distribution (statically linked)
- Cross-compilation support via `just build-all`
- No external runtime dependencies
- Version tracking with build metadata injection

## Development Practices
- Go idiomatic code style
- Comprehensive error handling
- Structured logging (if implemented)
- Unit testing for critical components
- Code documentation through comments
- Build automation with Justfile

## API Integration
- RESTful API communication with NetBox
- GraphQL queries for efficient data fetching
- JSON serialization/deserialization
- Authentication via API tokens
- HTTPS-only connections
- Proper HTTP status code handling

## Configuration Management
- Environment variables only for sensitive data (`nbjrnlctl_base_url` and `nbjrnlctl_api_key`)
- No configuration files for security and simplicity
- Command-line flags for runtime options (limit flag in list command)
- Secure storage practices for credentials

## Security Considerations
- API token security through environment variables only
- HTTPS enforcement
- Input validation
- Error message sanitization
- Secure credential handling
- No file-based credential storage
