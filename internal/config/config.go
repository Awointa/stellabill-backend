package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string
	Port        string
	DBConn      string
	JWTSecret   string
	FeatureFlags FeatureFlagConfig
}

type FeatureFlagConfig struct {
	DefaultEnabled bool
	LogDisabled    bool
	ConfigFile     string
}

func Load() Config {
	return Config{
		Env:       getEnv("ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		DBConn:    getEnv("DATABASE_URL", "postgres://localhost/stellarbill?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		FeatureFlags: FeatureFlagConfig{
			DefaultEnabled: getBoolEnv("FF_DEFAULT_ENABLED", false),
			LogDisabled:    getBoolEnv("FF_LOG_DISABLED", true),
			ConfigFile:     getEnv("FF_CONFIG_FILE", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if parsed, err := parseBool(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func parseBool(s string) (bool, error) {
	switch s {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}
