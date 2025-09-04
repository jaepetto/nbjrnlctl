# System Patterns: nb-jrnl-ctl

## Architecture Overview
The system follows a layered architecture pattern with clear separation of concerns:
- **Command Layer**: CLI interface and command routing (cmd/)
- **Business Logic Layer**: Core operations and workflows (internal/commands/)
- **Data Access Layer**: API client and model representations (internal/client/, internal/models/)
- **Utility Layer**: Shared utilities and configuration (pkg/utils/)

## Key Design Patterns

### Command Pattern
Commands are organized by operation (create, read, update, delete, list) rather than by entity, providing a consistent user experience where verbs remain constant and objects vary.

### Client Pattern
API interactions are encapsulated in a dedicated client (netbox.go) that handles HTTP communication, authentication, and basic error handling, abstracting these concerns from business logic.

### Model Pattern
Data structures are defined as Go structs in models package, providing type safety and clear data contracts between layers.

### Configuration Pattern
Utilities for configuration management are centralized, ensuring consistent handling of settings across all components.

## Component Relationships

### Command Flow
```
main.go → command handlers → business logic → client → NetBox API
```

### Data Flow
```
NetBox API → client → models → business logic → command output
```

## Error Handling Strategy
- Centralized error handling in client layer for API-related issues
- Command-specific error handling for user input validation
- Consistent error message formatting for user experience
- Graceful degradation when API connectivity is limited

## Configuration Management
- Environment variable support for sensitive data
- Configuration file support for persistent settings
- Clear precedence hierarchy (env vars > config file > defaults)
- Validation at startup to prevent runtime configuration errors

## Extensibility Points
- New command handlers can be added without modifying core routing
- Client methods can be extended for additional NetBox API endpoints
- Utility functions can be shared across components
- Model definitions can be expanded for additional data fields

## Performance Considerations
- Minimal API calls per operation
- Efficient data serialization/deserialization
- Connection reuse where possible
- Streaming output for large result sets
