# Implementation Plan

- [x] 1. Update configuration structures for multi-user support
  - Modify config.go to add UserConfig struct with name, items, zipcode, distance, and notifications fields
  - Add Condense boolean field to NotificationConfig struct
  - Restructure main Config struct to have global settings (interval, user_agent, verbose) and Users array
  - Implement backward compatibility detection for legacy configuration format
  - _Requirements: 2.2, 2.3, 2.4, 3.1, 3.2, 3.3_

- [x] 1.1 Write unit tests for configuration parsing
  - Create tests for multi-user configuration parsing
  - Test backward compatibility with legacy configurations
  - Test validation of user-specific settings
  - _Requirements: 4.1, 4.3_

- [x] 2. Enhance notification system with condensing capability
  - Add Condense field processing in NotificationManager constructor
  - Implement NotifyFoundItems method to handle multiple liquor items in single notification
  - Create condensed message formatting logic that lists all found items clearly
  - Update existing NotifyFound method to work with new condensing logic
  - _Requirements: 1.1, 1.2, 1.3, 6.1, 6.2, 6.3, 6.4_

- [x] 2.1 Write unit tests for notification condensing
  - Test condensed vs individual notification behavior
  - Test message formatting for multiple items
  - Test edge cases with single item and no items found
  - _Requirements: 4.1, 4.3, 6.1, 6.2, 6.4_

- [x] 3. Create UserRunner for individual user search management
  - Implement UserRunner struct with user-specific configuration, searcher, and notifier
  - Add Start method for UserRunner to handle periodic searches for one user
  - Implement runSearch method that collects all found items before sending notifications
  - Add proper error handling and logging with user identification
  - _Requirements: 2.2, 2.3, 2.5_

- [x] 4. Implement MultiUserRunner for coordinating multiple users
  - Create MultiUserRunner struct to manage multiple UserRunner instances
  - Implement concurrent execution of user searches with proper isolation
  - Add Start method that initializes and starts all user runners concurrently
  - Implement proper shutdown and cleanup for all user runners
  - Remove or replace existing Runner struct and related functions in runner package if not longer used in either UserRunner or MultiUserRunner functionality
  - _Requirements: 2.2, 2.5, 3.1, 3.2, 3.3_

- [x] 4.1 Write integration tests for multi-user execution
  - Test concurrent user search execution
  - Test user isolation and independence
  - Test proper cleanup and shutdown procedures
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 4.2 Consolidate or otherwise refactor runner package
  - Determine whether the single user and multi-user runner structs and related functions can be refactored, ideally with one Runner struct and related functions and then multi-user is a list of those or similar.
  - Refactor runner package as necessary, and adjust its usage in other packages as necessary
  - Refactor runner package tests and others as necessary
  - _Requirements: 2.2, 2.5, 3.1, 3.2, 3.3_

- [x] 5. Update main application entry point
  - Modify main.go and runner initialization to use new MultiUserRunner
  - Update command-line interface to work with multi-user configuration
  - Ensure proper error handling and logging for multi-user scenarios
  - Maintain backward compatibility for existing single-user setups
  - _Requirements: 2.2, 2.5, 3.1, 3.2, 3.3_

- [x] 6. Update configuration examples and documentation
  - Update config loading either in the config package or cmd package to ensure that if a config file is specified as a CLI argument, then the default config.yaml file should not be loaded and the user-specified config file is loaded. 
  - Update config.example.yaml to demonstrate multi-user configuration structure
  - Add examples showing both condensed and individual notification settings
  - Include backward compatibility examples and migration guidance
  - Update README.md with multi-user setup instructions and configuration explanations
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 6.1 Update existing unit tests for new configuration structure
  - Modify config_test.go to test new multi-user configuration parsing
  - Update any other existing tests that depend on old configuration structure
  - Ensure all tests pass with new multi-user configuration format
  - _Requirements: 4.1, 4.3_