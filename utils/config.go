package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration values
type Config struct {
	ListmonkURL string
	APIKey      string
	Port        string
	// Add other configuration fields as needed
}

// LoadConfig reads configuration from .env.local file and environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Read from .env.local file
	if err := loadEnvFile(".env.local", config); err != nil {
		ErrorLogger.Printf("Error loading .env.local file: %v", err)
		// Continue execution, as we'll fall back to environment variables
	}

	// Override with environment variables if they exist
	if envListmonkURL := os.Getenv("LISTMONK_URL"); envListmonkURL != "" {
		config.ListmonkURL = envListmonkURL
	}
	if envAPIKey := os.Getenv("API_KEY"); envAPIKey != "" {
		config.APIKey = envAPIKey
	}

	if envPort := os.Getenv("API_KEY"); envPort != "" {
		config.Port = envPort
	}

	// Validate required fields
	if config.ListmonkURL == "" {
		return nil, fmt.Errorf("LISTMONK_URL is not set")
	}
	if config.APIKey == "" {
		return nil, fmt.Errorf("API_KEY is not set")
	}
	if config.Port == "" {
		config.Port = "8808"
	}

	return config, nil
}

func loadEnvFile(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "LISTMONK_URL":
			config.ListmonkURL = value
		case "API_KEY":
			config.APIKey = value
		case "PORT":
			config.Port = value
			// Add other configuration fields as needed
		}
	}

	return scanner.Err()
}
