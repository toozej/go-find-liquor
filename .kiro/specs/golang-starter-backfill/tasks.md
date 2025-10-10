# Implementation Plan

- [x] 1. Update GitHub Actions CI/CD pipeline
- [x] 1.1 Fetch latest GitHub Actions workflow from golang-starter template
  - Download `.github/workflows/cicd.yaml` from golang-starter main branch
  - Compare with current go-find-liquor `.github/workflows/cicd.yaml`
  - Identify differences in action versions (CodeQL v3→v4, setup-go versions)
  - Identify new attestation support features
  - _Requirements: 2.2, 2.3_

- [x] 1.2 Merge GitHub Actions workflow changes
  - Update `github/codeql-action` references from v3 to v4
  - Add attestation support to goreleaser job if not present
  - Preserve go-find-liquor specific configurations (project names, secrets, etc.)
  - Resolve any merge conflicts between template and project-specific settings
  - _Requirements: 2.4, 2.5_

- [x] 1.3 Validate updated GitHub Actions workflow
  - Run workflow syntax validation
  - Test that all existing security scanning continues to work (Trivy, Snyk, Gitleaks)
  - Verify all go-find-liquor specific job steps are preserved
  - _Requirements: 2.5_

- [x] 2. Modernize Goreleaser configuration
- [x] 2.1 Fetch latest Goreleaser configuration from golang-starter template
  - Download `.goreleaser.yml` from golang-starter main branch
  - Compare with current go-find-liquor `.goreleaser.yml`
  - Identify differences in dockers vs dockers_v2 configuration
  - Identify new homebrew_casks section and OCI annotations
  - _Requirements: 3.2, 3.4_

- [x] 2.2 Migrate Goreleaser dockers to dockers_v2 format
  - Replace legacy `dockers` section with modern `dockers_v2` configuration
  - Preserve go-find-liquor specific image names and registry configurations
  - Add comprehensive OCI annotations and labels for better metadata
  - Configure both regular and distroless Docker image variants
  - _Requirements: 3.2, 3.5_

- [x] 2.3 Add homebrew_casks support for macOS distribution
  - Add homebrew_casks section to .goreleaser.yml
  - Configure tap token usage with TAP_GITHUB_TOKEN environment variable
  - Adapt cask metadata for go-find-liquor (name, description, etc.)
  - Add proper installation hooks and completions
  - _Requirements: 3.3, 3.4_

- [x] 2.4 Resolve Goreleaser configuration conflicts
  - Merge template improvements while preserving go-find-liquor specific settings
  - Update project name, binary name, and repository references
  - Ensure Docker image names and registry paths are correct
  - Validate final .goreleaser.yml syntax
  - _Requirements: 3.5_

- [x] 3. Update Docker configurations and metadata
- [x] 3.1 Fetch latest Docker configurations from golang-starter template
  - Download `Dockerfile.goreleaser.distroless` from golang-starter main branch
  - Compare with current go-find-liquor Dockerfile.goreleaser (note this does not have and should not have .distroless in the filename)
  - Identify improvements in build practices and metadata
  - _Requirements: 4.2_

- [x] 3.2 Update Docker build configurations
  - Update Dockerfile.goreleaser with template best practices
  - Update Docker labels and annotations for enhanced metadata
  - Remove any obsolete Docker files (like Dockerfile.docs) if present
  - _Requirements: 4.2, 4.3_

- [x] 3.3 Preserve go-find-liquor specific Docker settings
  - Ensure binary name and entrypoint scripts are correct for go-find-liquor
  - Verify any go-find-liquor specific dependencies or configurations
  - Update metadata to reflect go-find-liquor project information
  - _Requirements: 4.3_

- [x] 4. Enhance development tooling and configurations
- [x] 4.1 Fetch latest development tooling configurations from golang-starter template
  - Download `.pre-commit-config.yaml` from golang-starter main branch
  - Download `Makefile` from golang-starter main branch
  - Download `checkmake.ini` from golang-starter main branch
  - Compare each with current go-find-liquor versions
  - _Requirements: 5.2, 5.3_

- [x] 4.2 Update pre-commit configuration
  - Merge template pre-commit hooks with current configuration
  - Preserve any go-find-liquor specific pre-commit settings
  - Update hook versions to match template
  - _Requirements: 5.2_

- [x] 4.3 Update Makefile with template improvements
  - Merge template Makefile targets with current go-find-liquor Makefile
  - Preserve go-find-liquor specific targets and variables
  - Update any outdated commands or practices
  - Ensure all existing functionality is maintained
  - _Requirements: 5.3_

- [x] 4.4 Update additional development configurations
  - Update checkmake.ini configuration if needed
  - Review and update shell scripts in scripts/ directory
  - Update .dockerignore and .gitignore if needed
  - Gather and update .vscode/launch.json if needed
  - Merge any other development tooling improvements from template
  - _Requirements: 5.4_

- [ ] 5. Verify project-specific functionality preservation
- [ ] 5.1 Run comprehensive functionality tests
  - Execute complete test suite to ensure no regression in liquor search functionality
  - Validate notification system continues to work correctly
  - Test all go-find-liquor specific features and configurations
  - _Requirements: 6.3_

- [ ] 5.2 Validate project identity preservation
  - Verify binary name remains "go-find-liquor" throughout all configurations
  - Confirm project name and repository references are correct
  - Check that Docker image names and registry paths are accurate
  - Ensure all go-find-liquor specific metadata is preserved
  - _Requirements: 6.4_

- [ ] 6. Validate complete build and release process
- [ ] 6.1 Test local build process with new configurations
  - Run local build using updated Makefile targets
  - Test Docker builds with new configurations
  - Verify Goreleaser configuration with dry-run
  - _Requirements: 3.5, 4.4_

- [ ] 6.2 Validate cross-platform compatibility
  - Test builds on supported platforms (Linux, macOS, Windows)
  - Verify ARM64 and AMD64 architectures work correctly
  - Ensure package distributions are generated properly
  - Test homebrew cask functionality if applicable
  - _Requirements: 6.5_

- [ ] 6.3 End-to-end release process validation
  - Validate complete release workflow with updated configurations
  - Test Docker image functionality and distribution
  - Verify all signing and attestation processes work
  - Confirm release artifacts are generated correctly
  - _Requirements: 6.5_