package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	TokenExpireHours int
}

func Load() Config {
	expireHours, err := strconv.Atoi(getEnv("TOKEN_EXPIRE_HOURS", "72"))
	if err != nil {
		expireHours = 72
	}

	return Config{
		Port:             getEnvAny([]string{"PORT", "APP_PORT"}, "8080"),
		DBHost:           getEnvAny([]string{"DB_HOST", "MYSQL_HOST"}, "127.0.0.1"),
		DBPort:           getEnvAny([]string{"DB_PORT", "MYSQL_PORT"}, "3306"),
		DBUser:           getEnvAny([]string{"DB_USER", "MYSQL_USER"}, "gomall"),
		DBPassword:       getEnvAny([]string{"DB_PASSWORD", "MYSQL_PASSWORD"}, "gomall123"),
		DBName:           getEnvAny([]string{"DB_NAME", "MYSQL_DATABASE"}, "gomall_lite"),
		JWTSecret:        getEnv("JWT_SECRET", "gomall-lite-secret"),
		TokenExpireHours: expireHours,
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return fallback
}
