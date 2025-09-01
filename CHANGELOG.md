# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Enhanced Schema Registry command structure with intuitive syntax
- Automatic latest version detection for schema operations
- Smart error handling with user-friendly messages
- Positional arguments for better command-line experience
- **NEW**: Schema compatibility information display
- **NEW**: `--detail` flag for subjects command to show compatibility and version information
- Parallel API processing for improved performance
- Schema compatibility levels display (BACKWARD, FORWARD, FULL, etc.)

### Changed

- **BREAKING**: Simplified schema command syntax:
  - `kafta schema subjects-list` → `kafta schema subjects`
  - `kafta schema subjects-version --subject NAME` → `kafta schema versions NAME`
  - `kafta schema get --subject NAME` → `kafta schema get NAME`
  - `kafta schema diff --subject NAME` → `kafta schema diff NAME`
- Schema version parameter is now optional (defaults to latest)
- Improved error messages for schema registry operations
- Enhanced command validation with clear feedback
- **PERFORMANCE**: `subjects` command now uses single API call for faster execution
- **NEW**: `--detail` flag provides compatibility and version count information

### Fixed

- Duplicate help text display when validation fails
- Cryptic error messages replaced with user-friendly alternatives
- Better handling of Schema Registry API errors
- Performance issues with multiple API calls resolved

### Improved

- User experience with more intuitive command structure
- Error handling follows project standards using `cmdutil.CheckErr()`
- Consistent argument validation across all schema commands
- Better integration with existing project patterns
- **PERFORMANCE**: Optimized API calls with goroutines for detailed commands
- **UX**: Two-tier approach: fast basic listing vs detailed information

## Previous Versions

See git history for changes prior to this changelog.
