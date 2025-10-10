package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// NotificationConfig stores configuration for notification methods
type NotificationConfig struct {
	Type       string            `mapstructure:"type"`
	Endpoint   string            `mapstructure:"endpoint"`
	Credential map[string]string `mapstructure:"credential"`
	Condense   bool              `mapstructure:"condense"`
}

// UserConfig represents configuration for a single user
type UserConfig struct {
	Name          string               `mapstructure:"name"`
	Items         []string             `mapstructure:"items"`
	Zipcode       string               `mapstructure:"zipcode"`
	Distance      int                  `mapstructure:"distance"`
	Notifications []NotificationConfig `mapstructure:"notifications"`
}

// Config stores all configuration for the application
type Config struct {
	// Global settings
	Interval  time.Duration `mapstructure:"interval"`
	UserAgent string        `mapstructure:"user_agent"`
	Verbose   bool          `mapstructure:"verbose"`

	// User-specific configurations
	Users []UserConfig `mapstructure:"users"`

	// Legacy fields for backward compatibility (will be populated if old format detected)
	Items         []string             `mapstructure:"items,omitempty"`
	Zipcode       string               `mapstructure:"zipcode,omitempty"`
	Distance      int                  `mapstructure:"distance,omitempty"`
	Notifications []NotificationConfig `mapstructure:"notifications,omitempty"`
}

func GetConfig() (Config, error) {
	var config Config

	// Set default values
	viper.SetDefault("distance", 10)
	viper.SetDefault("interval", 12*time.Hour)
	viper.SetDefault("verbose", false)

	// Only load default config.yaml if no custom config file was set via CLI
	if viper.ConfigFileUsed() == "" {
		if _, err := os.Stat("config.yaml"); err == nil {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")

			if err := viper.ReadInConfig(); err != nil {
				return config, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// Enable environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GFL")

	// Map environment variables for both new and legacy formats
	_ = viper.BindEnv("interval", "GFL_INTERVAL")
	_ = viper.BindEnv("user_agent", "GFL_USER_AGENT")
	_ = viper.BindEnv("verbose", "GFL_VERBOSE")

	// Legacy environment variables for backward compatibility
	_ = viper.BindEnv("items", "GFL_ITEMS")
	_ = viper.BindEnv("zipcode", "GFL_ZIPCODE")
	_ = viper.BindEnv("distance", "GFL_DISTANCE")
	_ = viper.BindEnv("notifications", "GFL_NOTIFICATIONS")

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Check for legacy configuration format and migrate if needed
	if isLegacyConfig(config) {
		migratedConfig, err := migrateLegacyConfig(config)
		if err != nil {
			return config, fmt.Errorf("failed to migrate legacy config: %w", err)
		}
		config = migratedConfig
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return config, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// isLegacyConfig detects if the configuration is in the old format
func isLegacyConfig(config Config) bool {
	// Legacy format has items, zipcode, or notifications at root level
	// and no users array
	return len(config.Users) == 0 && (len(config.Items) > 0 || config.Zipcode != "" || len(config.Notifications) > 0)
}

// migrateLegacyConfig converts legacy configuration to multi-user format
func migrateLegacyConfig(config Config) (Config, error) {
	if len(config.Items) == 0 {
		return config, fmt.Errorf("legacy configuration must have items specified")
	}

	if config.Zipcode == "" {
		return config, fmt.Errorf("legacy configuration must have zipcode specified")
	}

	// Create a single user from legacy configuration
	user := UserConfig{
		Name:          "default",
		Items:         config.Items,
		Zipcode:       config.Zipcode,
		Distance:      config.Distance,
		Notifications: config.Notifications,
	}

	// Set default distance if not specified
	if user.Distance == 0 {
		user.Distance = 10
	}

	// Create new config with migrated user
	newConfig := Config{
		Interval:  config.Interval,
		UserAgent: config.UserAgent,
		Verbose:   config.Verbose,
		Users:     []UserConfig{user},
	}

	fmt.Printf("Migrated legacy configuration to multi-user format with user '%s'\n", user.Name)

	return newConfig, nil
}

// validateConfig validates the configuration structure
func validateConfig(config Config) error {
	if len(config.Users) == 0 {
		return fmt.Errorf("at least one user must be configured")
	}

	for i, user := range config.Users {
		if user.Name == "" {
			return fmt.Errorf("user %d must have a name", i)
		}

		if len(user.Items) == 0 {
			return fmt.Errorf("user '%s' must have at least one item to search for", user.Name)
		}

		if user.Zipcode == "" {
			return fmt.Errorf("user '%s' must have a zipcode specified", user.Name)
		}

		if user.Distance <= 0 {
			return fmt.Errorf("user '%s' must have a positive distance", user.Name)
		}
	}

	return nil
}
