# API Documentation: nb-jrnl-ctl

**Date:** 2025-09-10
**Author:** Cline (AI Assistant)
**Version:** 1.0.0

## Overview

The `nb-jrnl-ctl` tool interacts with the NetBox API to review journal entries. It uses both REST API endpoints and GraphQL queries to provide efficient data retrieval with enhanced formatting.

## NetBox API Endpoints

### Authentication

All API calls require a valid NetBox API token set via environment variables (`nbjrnlctl_api_key`).

### REST API Endpoints

#### Device Management

**GET `/api/dcim/devices/`**
- **Purpose**: Retrieve devices from NetBox for device lookup
- **Method**: GET
- **Pagination**: Supports pagination with `next` and `previous` fields
- **Response**: List of devices with `id` and `name` fields
- **Used by**: Device lookup in the list command

### GraphQL Endpoint

**POST `/graphql/`**
- **Purpose**: Fetch journal entries with rich relationship data
- **Method**: POST
- **Request Body**:
  ```json
  {
    "query": "GraphQL query string"
  }
  ```
- **Query Used**:
  ```graphql
  {
    device_list(filters: {id: {exact: $deviceID}}) {
      id
      name
      journal_entries {
        id
        created
        display
        comments
        created_by {
          username
        }
      }
    }
  }
  ```
- **Response**: Device with nested journal entries including creator information
- **Used by**: `list` command

## Data Models

### JournalEntry

```go
type JournalEntry struct {
    ID          int       `json:"id,omitempty"`
    AssignedObj string    `json:"assigned_object_type,omitempty"`
    AssignedID  int       `json:"assigned_object_id,omitempty"`
    Created     time.Time `json:"created,omitempty"`
    CreatedBy   string    `json:"created_by,omitempty"`
    Kind        string    `json:"kind,omitempty"`
    Comments    string    `json:"comments"`
}
```

### Configuration

```go
type Config struct {
    NetboxURL string `json:"netbox_url"`
    APIToken  string `json:"api_token"`
}
```

## API Response Processing

### Kind Extraction Logic

The tool extracts journal entry kinds from the GraphQL `display` field using multiple patterns:

1. **Primary Pattern**: `"date time (Kind)"` - extracts content between parentheses
2. **Fallback Pattern 1**: `"Kind: description"` or `"Kind - description"`
3. **Fallback Pattern 2**: Date pattern detection and kind extraction from remainder
4. **Final Fallback**: First word as kind

### Data Transformation

The GraphQL response is transformed from the API format to the internal `JournalEntry` model:

- **ID**: String to integer conversion
- **Created**: RFC3339 string to `time.Time`
- **CreatedBy**: Nested object to flat string
- **Kind**: Extracted from `display` field

## Error Handling

### HTTP Status Codes

- **200 OK**: Successful GET operations
- **400 Bad Request**: Invalid request data
- **401 Unauthorized**: Invalid or missing API token
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server-side errors

### Error Responses

All API errors are caught and presented to users with descriptive messages including the HTTP status code.

## Rate Limiting and Best Practices

### Client Configuration

- **Timeout**: 30 seconds for all HTTP requests
- **Connection Reuse**: Standard HTTP client connection pooling
- **Headers**: Proper authorization and content-type headers

### API Usage Patterns

1. **Device Lookup**: Paginated requests to find devices by name
2. **Data Retrieval**: Single GraphQL request for efficient data fetching
3. **Error Recovery**: Graceful handling of network and API errors

## Security Considerations

### Token Management

- API tokens are provided via environment variables only
- Tokens are never logged or displayed in output
- No file-based credential storage eliminates security risks

### Data Protection

- HTTPS-only connections enforced
- No sensitive data cached in memory longer than necessary
- Input validation prevents injection attacks
- Error messages sanitized to prevent information disclosure

## Performance Optimization

### Request Optimization

- Minimal field selection in GraphQL queries
- Single request per user operation
- Efficient pagination for large device datasets
- Connection reuse through HTTP client pooling

## Current Implementation Status

The tool currently implements read-only functionality focused on journal entry review:

- **Device Lookup**: REST API calls for efficient device ID resolution
- **Journal Entry Listing**: GraphQL queries for rich data retrieval
- **Enhanced Display**: Local processing for formatted output with colors and emojis
- **No Write Operations**: Focus on stable, read-only functionality for reliability
