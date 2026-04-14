package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	WeatherAPIKey string
	VKBotToken    string
	SessionSecret string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/sup_anapa?sslmode=disable"),
		WeatherAPIKey: getEnv("WEATHER_API_KEY", ""),
		VKBotToken:    getEnv("VK_BOT_TOKEN", ""),
		SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
