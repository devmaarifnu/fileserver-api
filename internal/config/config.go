package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Tokens   []TokenConfig  `mapstructure:"tokens"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
	Domain  string `mapstructure:"domain"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	BasePath          string   `mapstructure:"base_path"`
	MaxFileSize       int64    `mapstructure:"max_file_size"`
	AllowedExtensions []string `mapstructure:"allowed_extensions"`
}

// TokenConfig holds authentication token configuration
type TokenConfig struct {
	ID          string   `mapstructure:"id"`
	Key         string   `mapstructure:"key"`
	Name        string   `mapstructure:"name"`
	Permissions []string `mapstructure:"permissions"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	ValidateFileContent bool `mapstructure:"validate_file_content"`
	SanitizeFilename    bool `mapstructure:"sanitize_filename"`
}

// HasPermission checks if a token has a specific permission
func (t *TokenConfig) HasPermission(permission string) bool {
	for _, p := range t.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./")
	}

	// Enable environment variables override
	v.SetEnvPrefix("CDN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.App.Port)
	}

	if c.Storage.BasePath == "" {
		return fmt.Errorf("storage base path is required")
	}

	if c.Storage.MaxFileSize <= 0 {
		return fmt.Errorf("invalid max file size: %d", c.Storage.MaxFileSize)
	}

	if len(c.Storage.AllowedExtensions) == 0 {
		return fmt.Errorf("no allowed file extensions configured")
	}

	if len(c.Tokens) == 0 {
		return fmt.Errorf("no authentication tokens configured")
	}

	// Validate tokens
	for i, token := range c.Tokens {
		if token.ID == "" {
			return fmt.Errorf("token %d: ID is required", i)
		}
		if token.Key == "" {
			return fmt.Errorf("token %d: key is required", i)
		}
		if token.Name == "" {
			return fmt.Errorf("token %d: name is required", i)
		}
		if len(token.Permissions) == 0 {
			return fmt.Errorf("token %d: at least one permission is required", i)
		}
	}

	return nil
}

// FindTokenByKey finds a token by its key
func (c *Config) FindTokenByKey(key string) *TokenConfig {
	for i := range c.Tokens {
		if c.Tokens[i].Key == key {
			return &c.Tokens[i]
		}
	}
	return nil
}

// GetBaseURL returns the base URL based on environment
func (c *Config) GetBaseURL() string {
	if c.App.Env == "production" {
		return fmt.Sprintf("https://%s", c.App.Domain)
	}
	return fmt.Sprintf("http://localhost:%d", c.App.Port)
}
