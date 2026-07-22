package config

import "os"

type Config struct {
	Port        string
	DBPath      string
	JWTSecret   string
	ExporterURL string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./blog.db"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		ExporterURL: getEnv("EXPORTER_METRICS_URL", "http://127.0.0.1:9101/metrics"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
