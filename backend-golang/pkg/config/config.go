package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

// AppConfig menyimpan konfigurasi aplikasi
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

// ServerConfig menyimpan konfigurasi server
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	JWTSecret    string `mapstructure:"jwt_secret"`
}

// DatabaseConfig menyimpan konfigurasi database
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LoggerConfig menyimpan konfigurasi logger
type LoggerConfig struct {
	Level        string `mapstructure:"level"`
	FilePath     string `mapstructure:"file_path"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress     bool   `mapstructure:"compress"`
	ConsoleOutput bool  `mapstructure:"console_output"`
}

// RedisConfig menyimpan konfigurasi Redis
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// RabbitMQConfig menyimpan konfigurasi RabbitMQ
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

var (
	// Global config instance
	cfg *Config
)

// LoadConfig memuat konfigurasi dari file
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	// Read environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return cfg, nil
}

// GetConfig mengembalikan instance konfigurasi
func GetConfig() *Config {
	return cfg
}

// setDefaults mengatur nilai default untuk konfigurasi
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "sample-stack-golang")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "development")

	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 15)
	v.SetDefault("server.write_timeout", 15)
	v.SetDefault("server.idle_timeout", 60)
	v.SetDefault("server.jwt_secret", "your-secret-key")

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "sample_stack")
	v.SetDefault("database.sslmode", "disable")

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.file_path", "logs/app.log")
	v.SetDefault("logger.max_size", 100)
	v.SetDefault("logger.max_backups", 3)
	v.SetDefault("logger.max_age", 28)
	v.SetDefault("logger.compress", true)
	v.SetDefault("logger.console_output", true)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// RabbitMQ defaults
	v.SetDefault("rabbitmq.host", "localhost")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.user", "guest")
	v.SetDefault("rabbitmq.password", "guest")
	v.SetDefault("rabbitmq.vhost", "/")
} 