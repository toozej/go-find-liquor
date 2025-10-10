# Requirements Document

## Introduction

This feature enhances the go-find-liquor application to support multiple users with individual configurations while adding notification condensing capabilities. The system will allow each user to maintain their own search preferences (liquor items, notification methods, zipcode, and distance) while sharing global settings like search intervals and logging configuration. Additionally, users can choose to condense multiple liquor notifications from a single search run into a single notification message.

## Requirements

### Requirement 1

**User Story:** As a system administrator, I want to configure notification condensing globally, so that I can control whether multiple liquor findings are combined into single notifications.

#### Acceptance Criteria

1. WHEN the system loads configuration THEN it SHALL read a "condense" boolean field from the NotificationConfig struct
2. WHEN "condense" is set to false (default) THEN the system SHALL send individual notifications for each liquor item found
3. WHEN "condense" is set to true THEN the system SHALL combine all liquor findings from a single search run into one notification message
4. WHEN the config.example.yaml is updated THEN it SHALL include the new "condense" configuration option with appropriate documentation

### Requirement 2

**User Story:** As a user of a shared go-find-liquor installation, I want to have my own search configuration, so that I can search for different liquor items in my area without affecting other users.

#### Acceptance Criteria

1. WHEN the system loads configuration THEN it SHALL support multiple user configurations within a single installation
2. WHEN a user is configured THEN they SHALL have their own list of liquor items to search for
3. WHEN a user is configured THEN they SHALL have their own notification methods and settings
4. WHEN a user is configured THEN they SHALL have their own zipcode and search distance
5. WHEN multiple users exist THEN each user's search SHALL run independently without interfering with others

### Requirement 3

**User Story:** As a system administrator, I want to maintain global settings for all users, so that I can control system-wide behavior like search intervals and logging.

#### Acceptance Criteria

1. WHEN the system is configured THEN search intervals SHALL remain globally configured for all users
2. WHEN the system is configured THEN verbose logging settings SHALL remain globally configured for all users
3. WHEN global settings are changed THEN they SHALL apply to all user searches uniformly
4. WHEN user-specific settings are changed THEN they SHALL NOT affect global settings or other users

### Requirement 4

**User Story:** As a developer, I want existing unit tests to work with multi-user functionality, so that the system maintains reliability and test coverage.

#### Acceptance Criteria

1. WHEN unit tests are run THEN existing tests SHALL be updated to support multi-user configuration structure
2. WHEN new functionality is added THEN new unit tests SHALL be created to cover multi-user scenarios
3. WHEN tests are executed THEN they SHALL validate both single-user and multi-user configurations
4. WHEN notification condensing is tested THEN tests SHALL verify both condensed and individual notification behaviors

### Requirement 5

**User Story:** As a user, I want clear documentation on multi-user setup, so that I can configure the system correctly for my needs.

#### Acceptance Criteria

1. WHEN the config.example.yaml is updated THEN it SHALL demonstrate multi-user configuration structure
2. WHEN the README.md is updated THEN it SHALL include instructions for setting up multiple users
3. WHEN documentation is provided THEN it SHALL explain the difference between user-specific and global settings
4. WHEN examples are given THEN they SHALL show realistic multi-user scenarios with different configurations

### Requirement 6

**User Story:** As a user receiving notifications, I want condensed notifications to be clear and informative, so that I can understand all the liquor items found in a single message.

#### Acceptance Criteria

1. WHEN multiple liquor items are found in one search run AND condensing is enabled THEN the system SHALL create a single notification containing all items
2. WHEN a condensed notification is sent THEN it SHALL clearly list each liquor item found with relevant details
3. WHEN no liquor items are found THEN no notification SHALL be sent regardless of condensing setting
4. WHEN only one liquor item is found THEN the notification format SHALL be consistent whether condensing is enabled or disabled