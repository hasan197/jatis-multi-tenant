package config

// Config menyimpan konfigurasi aplikasi
type Config struct {
	DatabaseURL string
}

// Load memuat konfigurasi dari environment variables
func Load() (*Config, error) {
	return &Config{
		DatabaseURL: "postgres://postgres:postgres@postgres:5432/sample_db?sslmode=disable",
	}, nil
} 