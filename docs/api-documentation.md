# API Documentation: nb-jrnl-ctl

**Date:** 2025-09-04
**Author:** Cline (AI Assistant)
**Version:** 1.0.0

## Overview

The `nb-jrnl-ctl` tool interacts with the NetBox API to manage journal entries. It uses both REST API endpoints and GraphQL queries to provide comprehensive functionality.

## NetBox API Endpoints

### Authentication

All API calls require a valid NetBox API token set in the configuration file (`~/.nbjrnlctl/config.json`).

### REST API Endpoints

#### Device Management

**GET `/api/dcim/devices/`**
- **Purpose**: Retrieve devices from NetBox
- **Method**: GET
- **Pagination**: Supports pagination with `next` and `previous` fields
- **Response**: List of devices with `id` and `name` fields
- **Used by**: Device lookup in all commands

#### Journal Entry Management

**POST `/api/dcim/devices/{id}/journal/`**
- **Purpose**: Create a new journal entry for a specific device
- **Method**: POST
- **Path Parameters**:
  - `id` (integer): Device ID
- **Request Body**:
  ```json
  {
    "comments": "string",
    "kind": "string" // info, success, warning, danger
  }
  ```
- **Response**: Created journal entry object
- **Used by**: `create` command

**GET `/api/extras/journal-entries/{id}/`**
- **Purpose**: Retrieve a specific journal entry by ID
- **Method**: GET
- **Path Parameters**:
  - `id` (integer): Journal entry ID
- **Response**: Journal entry object with full details
- **Used by**: `read` and `update` commands

**PATCH `/api/extras/journal-entries/{id}/`**
- **Purpose**: Update an existing journal entry
- **Method**: PATCH
- **Path Parameters**:
  - `id` (integer): Journal entry ID
- **Request Body**:
  ```json
  {
    "comments": "string",
    "kind": "string"
  }
  ```
- **Response**: Updated journal entry object
- **Used by**: `update` command

**DELETE `/api/extras/journal-entries/{id}/`**
- **Purpose**: Delete a journal entry
- **Method**: DELETE
- **Path Parameters**:
  - `id` (integer): Journal entry ID
- **Response**: 204 No Content on success
- **Used by**: `delete` command

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

- **200 OK**: Successful GET/PATCH operations
- **201 Created**: Successful POST operations
- **204 No Content**: Successful DELETE operations
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
2. **Batch Operations**: Single request per operation (no bulk operations yet)
3. **Data Consistency**: Always fetch fresh data before updates
4. **Error Recovery**: Graceful handling of network and API errors

## Future API Extensions

### Planned Enhancements

1. **Bulk Operations**: Support for creating/updating multiple entries
2. **Advanced Filtering**: Query parameters for filtering journal entries
3. **Export Formats**: Additional response formats (CSV, JSON arrays)
4. **Webhook Support**: Real-time notifications for journal changes

### Potential New Endpoints

- **Search**: `/api/extras/journal-entries/?search=term`
- **Filtering**: `/api/extras/journal-entries/?device_id=id&kind=type`
- **Bulk Operations**: `/api/extras/journal-entries/bulk/`

## Security Considerations

### Token Management

- API tokens are stored in user's home directory with restricted permissions
- Tokens are never logged or displayed in output
- Configuration file creation uses secure file permissions (0644)

### Data Protection

- HTTPS-only connections enforced
- No sensitive data cached in memory longer than necessary
- Input validation prevents injection attacks
- Error messages sanitized to prevent information disclosure

## Performance Optimization

### Caching Strategy

Currently no caching is implemented, but future versions could cache:

- Device ID lookups for frequently accessed devices
- Configuration data
- Recently accessed journal entries

### Request Optimization

- Minimal field selection in GraphQL queries
- Single request per user operation
- Efficient pagination for large datasets
- Connection reuse through HTTP client pooling
