package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBHost       string
    DBPort       string
    DBUser       string
    DBPassword   string
    DBName       string
    DBSSLMode    string
    Port         string
    LogLevel     string
    MaxLimit     int
    DefaultLimit int
}

func Load() *Config {
    return &Config{
        DBHost:       getEnv("DB_HOST", "localhost"),
        DBPort:       getEnv("DB_PORT", "5433"),
        DBUser:       getEnv("DB_USER", "postgres"),
        DBPassword:   getEnv("DB_PASSWORD", "postgres"),
        DBName:       getEnv("DB_NAME", "phoneservice"),
        DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
        Port:         getEnv("PORT", "8080"),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
        MaxLimit:     getEnvAsInt("MAX_LIMIT", 100),
        DefaultLimit: getEnvAsInt("DEFAULT_LIMIT", 10),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
