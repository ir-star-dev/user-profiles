package configs

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	Db     DbConfig
	Secret string
}

type DbConfig struct {
	Dsn            string
	Driver         string
	MigrationsPath string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return &Config{
		Secret: os.Getenv("SECRET"),
		Db: DbConfig{
			Dsn:    buildDSN(),
			Driver: os.Getenv("DRIVER_NAME"),
			MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		},
	}, nil
}

func buildDSN() string {
    return fmt.Sprintf(
        "%s://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DRIVER_NAME"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
        os.Getenv("SSL_MODE"),
    )
}