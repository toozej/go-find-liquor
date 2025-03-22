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
}

// Config stores all configuration for the application
type Config struct {
	// Search parameters
	Items     []string      `mapstructure:"items"`
	Zipcode   string        `mapstructure:"zipcode"`
	Distance  int           `mapstructure:"distance"`
	Interval  time.Duration `mapstructure:"interval"`
	UserAgent string        `mapstructure:"user_agent"`
	Verbose   bool          `mapstructure:"verbose"`

	// Notification settings
	Notifications []NotificationConfig `mapstructure:"notifications"`
}

func GetConfig() (Config, error) {
	var config Config

	// Set default values
	viper.SetDefault("distance", 10)
	viper.SetDefault("interval", 12*time.Hour)
	viper.SetDefault("verbose", false)

	// Check for config file
	if _, err := os.Stat("config.yaml"); err == nil {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			return config, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// TODO fix merging in .env config items to config struct,
	// ignoring items in .env which shouldn't be in config struct...
	// Check for .env file
	// if _, err := os.Stat(".env"); err == nil {
	// 	viper.SetConfigFile(".env")
	// 	if err := viper.MergeInConfig(); err != nil {
	// 		return config, fmt.Errorf("failed to read .env file: %w", err)
	// 	}
	// }

	// Enable environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GFL")

	// Map environment variables
	_ = viper.BindEnv("items", "GFL_ITEMS")
	_ = viper.BindEnv("zipcode", "GFL_ZIPCODE")
	_ = viper.BindEnv("distance", "GFL_DISTANCE")
	_ = viper.BindEnv("interval", "GFL_INTERVAL")
	_ = viper.BindEnv("user_agent", "GFL_USER_AGENT")
	_ = viper.BindEnv("verbose", "GFL_VERBOSE")
	_ = viper.BindEnv("notifications", "GFL_NOTIFICATIONS")

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
