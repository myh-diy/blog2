package config

import (
	"fmt"
	"os"
	"strings"
)

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Exporter ExporterConfig `mapstructure:"exporter"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	PublicURL string `mapstructure:"public_url"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	TokenTTLHours int    `mapstructure:"token_ttl_hours"`
	AdminUsername string `mapstructure:"admin_username"`
	AdminPassword string `mapstructure:"admin_password"`
}

type ExporterConfig struct {
	MetricsURL string `mapstructure:"metrics_url"`
}

type StorageConfig struct {
	UploadDir   string `mapstructure:"upload_dir"`
	BackupDir   string `mapstructure:"backup_dir"`
	MaxUploadMB int64  `mapstructure:"max_upload_mb"`
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("/app")

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.public_url", "http://localhost:8080")
	v.SetDefault("database.path", "./blog.db")
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.token_ttl_hours", 168)
	v.SetDefault("auth.admin_username", "admin")
	v.SetDefault("auth.admin_password", "admin")
	v.SetDefault("exporter.metrics_url", "http://127.0.0.1:9101/metrics")
	v.SetDefault("storage.upload_dir", "uploads")
	v.SetDefault("storage.backup_dir", "data/backups")
	v.SetDefault("storage.max_upload_mb", 20)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindings := map[string][]string{
		"server.port":           {"PORT", "SERVER_PORT"},
		"server.mode":           {"GIN_MODE", "SERVER_MODE"},
		"server.public_url":     {"PUBLIC_URL", "SERVER_PUBLIC_URL"},
		"database.path":         {"DB_PATH", "DATABASE_PATH"},
		"auth.jwt_secret":       {"JWT_SECRET", "AUTH_JWT_SECRET"},
		"auth.token_ttl_hours":  {"JWT_TTL_HOURS", "AUTH_TOKEN_TTL_HOURS"},
		"auth.admin_username":   {"ADMIN_USERNAME", "AUTH_ADMIN_USERNAME"},
		"auth.admin_password":   {"ADMIN_PASSWORD", "AUTH_ADMIN_PASSWORD"},
		"exporter.metrics_url":  {"EXPORTER_METRICS_URL"},
		"storage.upload_dir":    {"UPLOAD_DIR", "STORAGE_UPLOAD_DIR"},
		"storage.backup_dir":    {"BACKUP_DIR", "STORAGE_BACKUP_DIR"},
		"storage.max_upload_mb": {"MAX_UPLOAD_MB", "STORAGE_MAX_UPLOAD_MB"},
	}
	for key, names := range bindings {
		args := append([]string{key}, names...)
		if err := v.BindEnv(args...); err != nil {
			return Config{}, err
		}
	}

	if filename := os.Getenv("CONFIG_FILE"); filename != "" {
		v.SetConfigFile(filename)
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if cfg.Auth.TokenTTLHours <= 0 {
		return Config{}, fmt.Errorf("auth.token_ttl_hours must be greater than zero")
	}
	if cfg.Storage.MaxUploadMB <= 0 {
		return Config{}, fmt.Errorf("storage.max_upload_mb must be greater than zero")
	}
	cfg.Server.PublicURL = strings.TrimRight(cfg.Server.PublicURL, "/")
	return cfg, nil
}
