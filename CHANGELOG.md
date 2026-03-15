# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Unit tests for ProjectService and WebShellService (2026-03-15)
- Project standard checker skill for automated compliance validation (2026-03-15)
- PowerShell script for project standards checking (2026-03-15)

### Changed
- Optimized project structure and documentation organization (2026-03-15)
- Updated .gitignore with comprehensive file exclusions (2026-03-15)

### Fixed
- WebShell restore functionality (GORM model specification) (2026-03-15)
- Repository layer soft delete operations (2026-03-15)

## [1.0.0] - 2026-03-15

### Added
- Initial release of FG-ABYSS
- Project management and WebShell management features
- Clean architecture implementation
- Vue 3 + TypeScript frontend
- Wails v3 desktop application
- Comprehensive documentation system
- Automated testing framework
- Project standards and guidelines

### Features
- **Project Management**
  - Create, read, update, delete projects
  - Soft delete support
  - Project validation
  - Unique name enforcement
  
- **WebShell Management**
  - Full CRUD operations
  - Soft delete and recovery
  - Pagination and search
  - Batch operations
  - URL and ProjectID validation
  
- **Architecture**
  - Clean architecture (Handler → Service → Repository)
  - Dependency injection
  - GORM database layer
  - Event system
  
- **Frontend**
  - Vue 3 + TypeScript
  - Naive UI components
  - i18n support (EN/CN)
  - Responsive design
  
- **Development Tools**
  - Taskfile.yml for task automation
  - Vitest for frontend testing
  - testify for Go testing
  - Project standard checker

### Technical Stack
- **Backend**: Go 1.21+
- **Frontend**: Vue 3.4+, TypeScript 5+
- **UI Framework**: Naive UI 2.37+
- **Database**: SQLite (via GORM)
- **Desktop**: Wails v3
- **Build Tools**: Vite 5+

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 1.0.0 | 2026-03-15 | Initial release with full features |

---

## Notes

- All changes are documented following the Keep a Changelog format
- Version numbers follow Semantic Versioning (MAJOR.MINOR.PATCH)
- Breaking changes are highlighted in the release notes
- Bug fixes reference issue numbers when applicable
