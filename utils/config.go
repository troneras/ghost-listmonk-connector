package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Config holds all configuration values
type Config struct {
	ListmonkURL string

	Port          string
	AUTH_USER     string
	AUTH_PASSWORD string
	JWT_SECRET    string

	// Magic link (email) configuration
	FrontendURL  string
	AWSRegion    string
	SESFromEmail string
	AWSAccessKey string
	AWSSecretKey string

	// Database configuration
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	// Redis configuration
	RedisAddr string
}

var (
	config     *Config
	configOnce sync.Once
)

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

	if envPort := os.Getenv("PORT"); envPort != "" {
		config.Port = envPort
	}
	if envAuthUser := os.Getenv("AUTH_USER"); envAuthUser != "" {
		config.AUTH_USER = envAuthUser
	}
	if envAuthPassword := os.Getenv("AUTH_PASSWORD"); envAuthPassword != "" {
		config.AUTH_PASSWORD = envAuthPassword
	}
	if envJWTSecret := os.Getenv("JWT_SECRET"); envJWTSecret != "" {
		config.JWT_SECRET = envJWTSecret
	}

	if envFrontendURL := os.Getenv("FRONTEND_URL"); envFrontendURL != "" {
		config.FrontendURL = envFrontendURL
	}
	if envAWSRegion := os.Getenv("AWS_REGION"); envAWSRegion != "" {
		config.AWSRegion = envAWSRegion
	}
	if envSESFromEmail := os.Getenv("SES_FROM_EMAIL"); envSESFromEmail != "" {
		config.SESFromEmail = envSESFromEmail
	}
	if envAWSAccessKey := os.Getenv("AWS_ACCESS_KEY_ID"); envAWSAccessKey != "" {
		config.AWSAccessKey = envAWSAccessKey
	}
	if envAWSSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY"); envAWSSecretKey != "" {
		config.AWSSecretKey = envAWSSecretKey
	}

	if envDBHost := os.Getenv("DB_HOST"); envDBHost != "" {
		config.DBHost = envDBHost
	}
	if envDBPort := os.Getenv("DB_PORT"); envDBPort != "" {
		config.DBPort = envDBPort
	}
	if envDBName := os.Getenv("DB_NAME"); envDBName != "" {
		config.DBName = envDBName
	}
	if envDBUser := os.Getenv("DB_USER"); envDBUser != "" {
		config.DBUser = envDBUser
	}
	if envDBPassword := os.Getenv("DB_PASSWORD"); envDBPassword != "" {
		config.DBPassword = envDBPassword
	}

	if envRedisAddr := os.Getenv("REDIS_ADDR"); envRedisAddr != "" {
		config.RedisAddr = envRedisAddr
	}

	// Validate required fields
	if config.ListmonkURL == "" {
		return nil, fmt.Errorf("LISTMONK_URL is not set")
	}
	if config.Port == "" {
		config.Port = "8808" // Default port if not set
	}
	if config.AUTH_USER == "" {
		return nil, fmt.Errorf("AUTH_USER is not set")
	}
	if config.AUTH_PASSWORD == "" {
		return nil, fmt.Errorf("AUTH_PASSWORD is not set")
	}
	if config.JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}
	if config.FrontendURL == "" {
		return nil, fmt.Errorf("FRONTEND_URL is not set")
	}
	if config.AWSRegion == "" {
		return nil, fmt.Errorf("AWS_REGION is not set")
	}
	if config.SESFromEmail == "" {
		return nil, fmt.Errorf("SES_FROM_EMAIL is not set")
	}
	if config.AWSAccessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID is not set")
	}
	if config.AWSSecretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY is not set")
	}
	if config.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is not set")
	}
	if config.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is not set")
	}
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is not set")
	}
	if config.RedisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is not set")
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
		case "PORT":
			config.Port = value
		case "AUTH_USER":
			config.AUTH_USER = value
		case "AUTH_PASSWORD":
			config.AUTH_PASSWORD = value
		case "JWT_SECRET":
			config.JWT_SECRET = value
		case "FRONTEND_URL":
			config.FrontendURL = value
		case "AWS_REGION":
			config.AWSRegion = value
		case "SES_FROM_EMAIL":
			config.SESFromEmail = value
		case "AWS_ACCESS_KEY_ID":
			config.AWSAccessKey = value
		case "AWS_SECRET_ACCESS_KEY":
			config.AWSSecretKey = value
		case "DB_HOST":
			config.DBHost = value
		case "DB_PORT":
			config.DBPort = value
		case "DB_NAME":
			config.DBName = value
		case "DB_USER":
			config.DBUser = value
		case "DB_PASSWORD":
			config.DBPassword = value
		case "REDIS_ADDR":
			config.RedisAddr = value

			// Add other configuration fields as needed
		}
	}

	return scanner.Err()
}

func GetConfig() *Config {
	configOnce.Do(func() {
		var err error
		config, err = LoadConfig()
		if err != nil {
			ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}
	})
	return config
}
