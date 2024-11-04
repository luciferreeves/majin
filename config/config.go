package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server settings
	Port         string
	Environment  string
	SecretKey    string
	Debug        bool
	SecureCookie bool

	// Mail server settings
	MailServer struct {
		IMAPHost    string
		IMAPPort    int
		SMTPHost    string
		SMTPPort    int
		UseStartTLS bool
	}

	// Authorization settings
	AuthorizedDomains []string
}

// Load reads the environment variables and returns a Config struct
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{}

	// Server settings
	config.Port = getEnv("PORT", "3000")
	config.Environment = getEnv("ENV", "development")
	config.SecretKey = getEnv("SECRET_KEY", "")

	// Set Debug and SecureCookie based on environment
	config.Debug = config.Environment == "development"
	config.SecureCookie = config.Environment == "production"

	// Mail server settings
	config.MailServer.IMAPHost = getEnv("MAIL_IMAP_HOST", "")
	config.MailServer.IMAPPort = getEnvAsInt("MAIL_IMAP_PORT", 993)
	config.MailServer.SMTPHost = getEnv("MAIL_SMTP_HOST", "")
	config.MailServer.SMTPPort = getEnvAsInt("MAIL_SMTP_PORT", 587)
	config.MailServer.UseStartTLS = getEnvAsBool("MAIL_SMTP_USE_STARTTLS", true)

	// Authorized domains
	domains := getEnv("AUTHORIZED_DOMAINS", "")
	if domains != "" {
		config.AuthorizedDomains = splitAndTrim(domains)
	}

	// Validate required settings
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.SecretKey == "" {
		return errors.New("SECRET_KEY is required")
	}
	if c.MailServer.IMAPHost == "" {
		return errors.New("MAIL_IMAP_HOST is required")
	}
	if c.MailServer.SMTPHost == "" {
		return errors.New("MAIL_SMTP_HOST is required")
	}
	if len(c.AuthorizedDomains) == 0 {
		return errors.New("AUTHORIZED_DOMAINS is required")
	}
	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
