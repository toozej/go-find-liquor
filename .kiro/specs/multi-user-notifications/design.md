# Design Document

## Overview

This design implements multi-user support and notification condensing for the go-find-liquor application. The system will be restructured to support multiple users with individual configurations while maintaining global settings for system-wide behavior. Additionally, a notification condensing feature will allow multiple liquor findings from a single search run to be combined into a single notification message.

## Architecture

### Current Architecture
The current system has a flat configuration structure where all settings (items, zipcode, distance, notifications) are at the root level, with a single runner managing searches for one user.

### New Architecture
The new architecture will introduce a hierarchical configuration structure:

```
Global Settings (system-wide):
├── Search Interval
├── Verbose Logging
└── Users (array):
    ├── User 1:
    │   ├── Items to search
    │   ├── Zipcode & Distance
    │   └── Notifications (with condense setting)
    └── User 2:
        ├── Items to search
        ├── Zipcode & Distance
        └── Notifications (with condense setting)
```

## Components and Interfaces

### 1. Configuration Changes

#### Updated Config Structures
```go
// UserConfig represents configuration for a single user
type UserConfig struct {
    Name          string                `mapstructure:"name"`
    Items         []string              `mapstructure:"items"`
    Zipcode       string                `mapstructure:"zipcode"`
    Distance      int                   `mapstructure:"distance"`
    Notifications []NotificationConfig  `mapstructure:"notifications"`
}

// NotificationConfig enhanced with condensing capability
type NotificationConfig struct {
    Type       string            `mapstructure:"type"`
    Endpoint   string            `mapstructure:"endpoint"`
    Credential map[string]string `mapstructure:"credential"`
    Condense   bool              `mapstructure:"condense"`  // New field
}

// Config restructured for multi-user support
type Config struct {
    // Global settings
    Interval  time.Duration `mapstructure:"interval"`
    UserAgent string        `mapstructure:"user_agent"`
    Verbose   bool          `mapstructure:"verbose"`
    
    // User-specific configurations
    Users     []UserConfig  `mapstructure:"users"`
}
```

### 2. Runner Architecture Changes

#### Multi-User Runner
The runner will be enhanced to manage multiple user searches concurrently:

```go
type MultiUserRunner struct {
    config       Config
    userRunners  map[string]*UserRunner
    stopChan     chan struct{}
}

type UserRunner struct {
    userConfig   UserConfig
    searcher     *search.Searcher
    notifier     *notification.NotificationManager
    stopChan     chan struct{}
    runningCh    chan struct{}
}
```

### 3. Notification Enhancement

#### Condensed Notification Support
The notification system will be enhanced to support condensing multiple findings:

```go
type NotificationManager struct {
    notifiers []Notifier
    condense  bool  // New field to control condensing behavior
}

// New method for handling multiple items
func (m *NotificationManager) NotifyFoundItems(ctx context.Context, items []search.LiquorItem) error
```

## Data Models

### Configuration File Structure
The new YAML configuration will support the multi-user structure:

```yaml
# Global settings
interval: 6h
verbose: true
user_agent: "custom-agent"  # optional

# User configurations
users:
  - name: "user1"
    items:
      - "Blanton's"
      - "W.L. Weller Special Reserve"
    zipcode: "97201"
    distance: 15
    notifications:
      - type: gotify
        endpoint: "https://gotify.example.com"
        condense: false  # New field
        credential:
          token: "USER1_GOTIFY_TOKEN"
  
  - name: "user2"
    items:
      - "1942"
      - "54633"
    zipcode: "97210"
    distance: 10
    notifications:
      - type: slack
        condense: true  # New field
        credential:
          token: "USER2_SLACK_TOKEN"
          channel_id: "CHANNEL_ID"
```

### Backward Compatibility
The system will maintain backward compatibility by:
1. Detecting legacy configuration format
2. Automatically migrating to single-user format
3. Providing clear migration guidance in logs

## Error Handling

### Configuration Validation
- Validate that at least one user is configured
- Ensure each user has required fields (name, items, zipcode)
- Validate notification configurations per user
- Provide clear error messages for configuration issues

### Runtime Error Handling
- Individual user search failures should not affect other users
- Notification failures should be logged but not stop searches
- Graceful degradation when notification condensing fails

### Migration Error Handling
- Detect and handle legacy configuration formats
- Provide clear migration instructions
- Fail gracefully with helpful error messages

## Testing Strategy

### Unit Tests
- Configuration parsing and validation for multi-user setup
- Notification condensing logic
- User runner isolation and independence
- Backward compatibility with legacy configurations

### Integration Tests
- Multi-user search execution
- Notification delivery for condensed vs individual messages
- Configuration migration scenarios
- Error scenarios and recovery

### Test Data
- Sample multi-user configurations
- Legacy configuration examples
- Various notification condensing scenarios
- Error condition test cases

## Implementation Phases

### Phase 1: Configuration Enhancement
1. Update configuration structures for multi-user support
2. Add condense field to NotificationConfig
3. Implement configuration parsing and validation
4. Add backward compatibility support

### Phase 2: Notification Condensing
1. Enhance NotificationManager to support condensing
2. Implement condensed message formatting
3. Update notification methods to handle multiple items
4. Add configuration option processing

### Phase 3: Multi-User Runner
1. Create UserRunner for individual user management
2. Implement MultiUserRunner for coordinating multiple users
3. Add concurrent search execution
4. Implement proper isolation between users

### Phase 4: Testing and Documentation
1. Update existing unit tests for new structure
2. Add comprehensive multi-user test scenarios
3. Update configuration examples and README
4. Add migration documentation

## Security Considerations

### User Isolation
- Ensure user configurations are properly isolated
- Prevent cross-user data leakage in notifications
- Validate user-specific settings independently

### Configuration Security
- Validate notification credentials per user
- Ensure sensitive data is properly handled
- Maintain secure defaults for new configuration options

## Performance Considerations

### Concurrent Execution
- Users will run searches concurrently but with proper rate limiting
- Shared searcher instances to avoid overwhelming the target service
- Proper resource cleanup for user runners

### Memory Management
- Efficient handling of multiple user configurations
- Proper cleanup of notification managers
- Avoid memory leaks in long-running multi-user scenarios

## Migration Strategy

### Automatic Detection
The system will automatically detect legacy configurations and either:
1. Migrate them to single-user format automatically
2. Provide clear instructions for manual migration

### Migration Process
1. Detect legacy format (items/zipcode/distance at root level)
2. Create single user configuration with existing settings
3. Preserve all existing functionality
4. Log migration actions for user awareness