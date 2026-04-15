package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	DBDSN    string

	RedisAddr string

	JWTSecret string

	APIAddr        string
	TrackerAddr    string
	TrackerBaseURL string

	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	c := &Config{
		DBHost:         env("DB_HOST", "127.0.0.1"),
		DBPort:         env("DB_PORT", "3306"),
		DBUser:         env("DB_USER", "phishguard"),
		DBPass:         env("DB_PASS", "phishguard_dev"),
		DBName:         env("DB_NAME", "phishguard"),
		RedisAddr:      env("REDIS_ADDR", "127.0.0.1:6379"),
		JWTSecret:      env("JWT_SECRET", "dev-secret-change-me"),
		APIAddr:        env("API_ADDR", ":8080"),
		TrackerAddr:    env("TRACKER_ADDR", ":8090"),
		TrackerBaseURL: env("TRACKER_BASE_URL", "http://localhost:8090"),
		AdminEmail:     env("ADMIN_EMAIL", "admin@phishguard.local"),
		AdminPassword:  env("ADMIN_PASSWORD", "changeme123"),
	}
	c.DBDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	return c
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
