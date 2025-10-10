# Requirements Document

## Introduction

This feature involves back-filling necessary changes and improvements from the golang-starter template repository (https://github.com/toozej/golang-starter) from commit 17cfa3a6e25170b6dea6b779ebde0f2b934bd8c2 (May 30, 2025) until the most recent commit into the go-find-liquor repository. The golang-starter template serves as a foundation for Go projects and contains best practices, tooling configurations, and infrastructure improvements that should be incorporated into go-find-liquor to maintain consistency and leverage the latest improvements.

## Requirements


### Requirement 1

**User Story:** As a developer, I want to incorporate GitHub Actions workflow improvements from golang-starter, so that the CI/CD pipeline benefits from the latest security and functionality enhancements.

#### Acceptance Criteria

1. WHEN reviewing GitHub Actions changes THEN the system SHALL identify workflow updates including action version bumps and new features
2. WHEN updating workflows THEN the system SHALL incorporate CodeQL action updates from v3 to v4
3. WHEN updating workflows THEN the system SHALL incorporate setup-go action updates from v5 to v6
4. WHEN updating workflows THEN the system SHALL add attestation support in the release workflow
5. WHEN workflow updates are complete THEN the system SHALL verify all workflows pass successfully

### Requirement 2

**User Story:** As a developer, I want to incorporate Goreleaser configuration improvements from golang-starter, so that the release process benefits from modern Docker and Homebrew distribution methods.

#### Acceptance Criteria

1. WHEN reviewing Goreleaser changes THEN the system SHALL identify updates to .goreleaser.yml configuration
2. WHEN updating Goreleaser config THEN the system SHALL migrate to dockers_v2 configuration format
3. WHEN updating Goreleaser config THEN the system SHALL add homebrew_casks support for macOS distribution
4. WHEN updating Goreleaser config THEN the system SHALL configure tap token usage for Homebrew taps
5. WHEN Goreleaser updates are complete THEN the system SHALL verify release builds work correctly

### Requirement 3

**User Story:** As a developer, I want to incorporate Docker configuration improvements from golang-starter, so that containerization benefits from the latest best practices and optimizations.

#### Acceptance Criteria

1. WHEN reviewing Docker changes THEN the system SHALL identify updates to Dockerfile configurations
2. WHEN updating Docker configs THEN the system SHALL remove obsolete Dockerfile.docs if present
3. WHEN updating Docker configs THEN the system SHALL incorporate any new distroless or goreleaser-specific Dockerfiles
4. WHEN Docker updates are complete THEN the system SHALL verify all Docker builds work correctly

### Requirement 4

**User Story:** As a developer, I want to incorporate tooling and configuration improvements from golang-starter, so that development workflow benefits from the latest best practices.

#### Acceptance Criteria

1. WHEN reviewing tooling changes THEN the system SHALL identify updates to pre-commit configuration, Makefile, and other development tools
2. WHEN updating tooling THEN the system SHALL incorporate pre-commit hook improvements
3. WHEN updating tooling THEN the system SHALL update Makefile targets and functionality as appropriate
4. WHEN updating tooling THEN the system SHALL ensure checkmake.ini and other configuration files are current
5. WHEN updating tooling THEN the system SHALL incorporate GitHub Actions workflow improvements in .github/workflows/cicd.yaml
6. WHEN updating tooling THEN the system SHALL incorporate DevContainer functionality in .devcontainer/devcontainer.json
7. WHEN tooling updates are complete THEN the system SHALL verify all development commands work correctly

### Requirement 5

**User Story:** As a developer, I want to maintain project-specific customizations while incorporating template improvements, so that go-find-liquor retains its unique functionality while benefiting from template updates.

#### Acceptance Criteria

1. WHEN applying template changes THEN the system SHALL preserve go-find-liquor specific code and configurations
2. WHEN conflicts arise THEN the system SHALL prioritize go-find-liquor specific requirements over template defaults
3. WHEN merging changes THEN the system SHALL maintain existing notification functionality and liquor search features
4. WHEN updates are complete THEN the system SHALL verify all go-find-liquor specific features continue to work
5. WHEN updates are complete THEN the system SHALL document any breaking changes or required manual interventions