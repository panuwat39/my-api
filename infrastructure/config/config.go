package config

import (
	"os"
	"time"
)

type Config struct {
	AppEnv        string
	HTTPPort      string
	MongoURI      string
	MongoDatabase string

	JWTSecret    string
	JWTIssuer    string
	JWTAccessTTL time.Duration
}

func Load() Config {
	return Config{
		AppEnv:        getEnv("APP_ENV", "local"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", ""),
		MongoDatabase: getEnv("MONGO_DATABASE", ""),

		JWTSecret: getEnv(
			"JWT_SECRET",
			"",
		),
		JWTIssuer: getEnv(
			"JWT_ISSUER",
			"my-api",
		),
		JWTAccessTTL: getDurationEnv(
			"JWT_ACCESS_TTL",
			15*time.Minute,
		),
	}
}

func getEnv(
	key string,
	fallback string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getDurationEnv(
	key string,
	fallback time.Duration,
) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}
