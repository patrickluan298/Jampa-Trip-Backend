package testutils

import (
	"os"
	"time"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	Environment      string
}

// LoadTestConfig loads test configuration
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL:      getEnv("TEST_DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable"),
		RedisURL:         getEnv("TEST_REDIS_URL", "localhost:6379"),
		JWTSecret:        getEnv("TEST_JWT_SECRET", "test-secret-key-for-testing-only"),
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
		Environment:      "test",
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// MockConfig returns a mock configuration for testing
func MockConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL:      "mock://test",
		RedisURL:         "mock://test",
		JWTSecret:        "mock-secret-key",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
		Environment:      "test",
	}
}

