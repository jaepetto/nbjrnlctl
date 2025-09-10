# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation review and synchronization
- Review report documenting the alignment process
- Updated all documentation to reflect current implementation reality

### Changed
- All documentation files synchronized with current codebase state
- Memory-bank files updated to reflect focused implementation
- Removed all references to removed command functionality
- Updated configuration documentation to environment variables only
- Enhanced user documentation with current feature descriptions

## [1.0.0] - 2025-09-04

### Added
- Initial release of nb-jrnl-ctl
- Basic CRUD operations for NetBox journal entries
- Command-line interface for journal management
- Device-specific journal operations
- Configuration management utilities
- Comprehensive code assessment documentation
- Detailed API endpoint documentation
- Complete developer guide and best practices
- Documentation directory structure
- Justfile for build automation and development workflow
- Built-in version tracking system with detailed build metadata
- Version command to display application version information
- Compile-time version injection using ldflags

### Changed
- Updated justfile build commands to create statically linked binaries for better portability
- Implemented go-pretty library for enhanced table formatting in list command output
- Replaced manual ANSI color codes with go-pretty styling for better maintainability
- Improved output formatting with rounded table borders and professional appearance
- **Configuration loading now uses environment variables instead of config file for better server deployment support**
- **Removed file-based configuration in favor of environment variables `nbjrnlctl_base_url` and `nbjrnlctl_api_key`**

### Fixed
- Journal entry IDs now display correctly in the list command output
- Fixed ID mapping in GraphQL response parsing for journal entries
- Populated the 'kind' column in journal entry listings by extracting it from the GraphQL 'display' field
- Improved 'kind' extraction logic to properly handle NetBox display format "YYYY-MM-DD HH:MM (Kind)"
- **Constrained table width to 80 characters for better terminal compatibility** - Table output now automatically wraps long text and fits within standard terminal widths

### Changed
- Enhanced list command output formatting:
  - Removed ID column from output
  - Shortened date format to MM/DD HH:MM
  - Replaced kind text with emojis (ℹ️ Info, ✅ Success, ⚠️ Warning, 🚨 Danger)
  - Removed comment truncation to show full comments
  - Added colored output using go-pretty styling
  - Sort entries by creation date (newest first)
  - **Added automatic text wrapping and column width constraints for better readability on smaller terminals**
  - **Added horizontal line separators between all rows for improved visual separation and readability**

### Removed
- Unused `getColorForKind` function that was replaced with go-pretty styling
- Unused command implementations (create, read, update, delete, device) - only list command remains as it's the only implemented verb
- **File-based configuration system in favor of environment variable-only configuration**
