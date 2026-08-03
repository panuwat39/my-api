package config

import "os"

type Config struct {
	AppEnv        string
	HTTPPort      string
	MongoURI      string
	MongoDatabase string
}

func Load() Config {
	return Config{
		AppEnv:        getEnv("APP_ENV", "local"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", ""),
		MongoDatabase: getEnv("MONGO_DATABASE", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
