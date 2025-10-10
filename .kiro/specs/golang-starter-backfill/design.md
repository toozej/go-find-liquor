# Design Document

## Overview

This design outlines the systematic approach to back-fill improvements from the golang-starter template repository into go-find-liquor. The changes span from commit 17cfa3a6e25170b6dea6b779ebde0f2b934bd8c2 (May 30, 2025) to the most recent commits, focusing on dependency updates, CI/CD improvements, Docker configuration enhancements, and Goreleaser modernization.

The design prioritizes maintaining go-find-liquor's unique functionality while incorporating template improvements that enhance security, maintainability, and distribution capabilities.

## Architecture

### Change Categories

The back-fill process is organized into five main categories:

1. **CI/CD Pipeline**: GitHub Actions workflow enhancements and security updates
2. **Release Management**: Goreleaser configuration modernization
3. **Container Infrastructure**: Docker configuration improvements
4. **Development Tooling**: Build tools and development workflow enhancements

### Migration Strategy

The migration follows a phased approach:
- **Phase 1**: CI/CD pipeline improvements
- **Phase 2**: Goreleaser modernization
- **Phase 3**: Docker configuration updates
- **Phase 4**: Development tooling enhancements

## Components and Interfaces

### GitHub Actions Enhancement Component

**Purpose**: Modernize CI/CD pipeline with latest action versions and security features

**Key Changes**:
- Update `github/codeql-action` from v3 to v4
- Update `actions/setup-go` from v5 to v6 (already current)
- Add attestation support in release workflow
- Maintain existing security scanning (Trivy, Snyk, Gitleaks)

**Interface**: GitHub Actions workflow files in `.github/workflows/`

### Goreleaser Modernization Component

**Purpose**: Migrate from legacy Docker configuration to modern dockers_v2 format

**Key Changes**:
- Replace `dockers` section with `dockers_v2` configuration
- Add `homebrew_casks` support for macOS distribution
- Configure tap token usage with `TAP_GITHUB_TOKEN`
- Add comprehensive OCI annotations and labels
- Support both regular and distroless Docker images

**Interface**: `.goreleaser.yml` configuration file

### Docker Configuration Component

**Purpose**: Update Docker build configurations and remove obsolete files

**Key Changes**:
- Remove `Dockerfile.docs` if present (not applicable to go-find-liquor)
- Ensure `Dockerfile.goreleaser` and `Dockerfile.goreleaser.distroless` are current
- Update Docker labels and annotations for better metadata

**Interface**: Dockerfile configurations and build contexts

### Development Tooling Component

**Purpose**: Enhance development workflow with updated tooling configurations

**Key Changes**:
- Update pre-commit configuration if needed
- Ensure Makefile targets are current with template
- Verify checkmake.ini configuration
- Update any shell scripts in `scripts/` directory

**Interface**: Configuration files and build scripts

## Data Models

### Configuration Mapping

```yaml
# GitHub Actions Updates
github_actions:
  - name: "github/codeql-action"
    from: "v3"
    to: "v4"
    files: [".github/workflows/cicd.yaml"]

# Goreleaser Changes
goreleaser_changes:
  - section: "dockers"
    action: "replace_with_dockers_v2"
  - section: "homebrew_casks"
    action: "add"
```

### Project-Specific Customizations

```yaml
# go-find-liquor specific configurations that must be preserved
project_customizations:
  binary_name: "go-find-liquor"
  project_name: "go-find-liquor"
  description: "Find liquor stores and their inventory"
  repository: "github.com/toozej/go-find-liquor"
  docker_images:
    - "toozej/go-find-liquor"
    - "ghcr.io/toozej/go-find-liquor"
    - "quay.io/toozej/go-find-liquor"
```

## Error Handling

### Build Failures

**Strategy**: Incremental validation
- Test each configuration change independently
- Maintain backup of working configurations
- Use feature flags for gradual rollout

### CI/CD Pipeline Issues

**Strategy**: Workflow validation
- Test workflow changes in feature branches
- Validate action version compatibility
- Ensure all required secrets and tokens are available

### Docker Build Problems

**Strategy**: Multi-stage validation
- Test Docker builds locally before deployment
- Validate both regular and distroless variants
- Ensure all required files are included in build context

## Testing Strategy

### Unit Testing

**Scope**: Core functionality preservation
- Run existing test suite after each change category
- Verify no regression in liquor search functionality
- Validate notification system continues to work

### Integration Testing

**Scope**: End-to-end workflow validation
- Test complete build process with new configurations
- Validate Docker image functionality
- Verify release process with updated Goreleaser config

### Security Testing

**Scope**: Maintain security posture
- Run security scans after dependency updates
- Validate CodeQL analysis with v4 action
- Ensure container security scanning continues to work

### Compatibility Testing

**Scope**: Cross-platform functionality
- Test builds on all supported platforms (Linux, macOS, Windows)
- Validate ARM64 and AMD64 architectures
- Ensure package distributions work correctly

## Implementation Phases

### Phase 1: CI/CD Enhancements
- Update GitHub Actions to latest versions
- Add attestation support
- Test complete CI/CD pipeline

### Phase 2: Goreleaser Modernization
- Migrate to dockers_v2 configuration
- Add homebrew_casks support
- Configure tap token usage

### Phase 3: Docker Improvements
- Update Docker configurations
- Add distroless variant support
- Enhance metadata and labels

### Phase 4: Tooling Updates
- Update development tooling configurations
- Verify all make targets work correctly
- Update documentation as needed

## Risk Mitigation

### Breaking Changes
- Test each change in isolation
- Maintain rollback capability
- Document any manual intervention required

### Security Considerations
- Validate all new action versions for security
- Ensure secrets management remains secure
- Verify container security posture

### Compatibility Issues
- Test on multiple Go versions
- Validate cross-platform builds
- Ensure backward compatibility where possible

## Success Criteria

1. GitHub Actions workflows updated and passing
2. Goreleaser configuration modernized with dockers_v2 and homebrew_casks
3. Docker builds working with enhanced metadata
4. All existing go-find-liquor functionality preserved
5. Release process working with new configurations
6. Development workflow enhanced with updated tooling