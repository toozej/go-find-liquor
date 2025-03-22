package config

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestGetConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	tests := []struct {
		name          string
		mockEnv       map[string]string
		mockEnvFile   string
		expectError   bool
		expectZipcode string
	}{
		{
			name: "Valid environment variable",
			mockEnv: map[string]string{
				"zipcode": "testuser",
			},
			expectError:   false,
			expectZipcode: "",
		},
		{
			name:          "Valid .env file",
			mockEnvFile:   "zipcode=testenvfileuser\n",
			expectError:   false,
			expectZipcode: "testenvfileuser",
		},
		{
			name:          "No environment variables or .env file",
			expectError:   false,
			expectZipcode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset Viper settings before each test
			viper.Reset()

			// Mock environment variables
			for key, value := range tt.mockEnv {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Mock .env file if applicable
			if tt.mockEnvFile != "" {
				if err := afero.WriteFile(fs, ".env", []byte(tt.mockEnvFile), 0644); err != nil {
					t.Fatalf("Failed to write mock .env file: %v", err)
				}
				viper.SetFs(fs) // Ensure Viper uses the mocked filesystem
				viper.SetConfigFile(".env")
				if err := viper.ReadInConfig(); err != nil {
					t.Fatalf("failed to read mock .env file: %v", err)
				}
			}

			// Call function
			conf, err := GetConfig()
			if err != nil {
				t.Errorf("GetConfig() returned error %v", err)
			}

			// Verify output
			if conf.Zipcode != tt.expectZipcode {
				t.Errorf("expected zipcode %q, got %q", tt.expectZipcode, conf.Zipcode)
			}
		})
	}
}
