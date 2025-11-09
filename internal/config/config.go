package config
import (
	"os"
	"github.com/joho/godotenv"
	"github.com/alibaba0010/postgres-api/internal/logger"
	"go.uber.org/zap"

)


type Config struct {
	Port string
	DB_HOST string
	DB_PORT string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME string
}

func LoadConfig() Config {
		err := godotenv.Load()
	if err != nil {
		// log.Println("No .env file found, using system environment variables...")
		logger.Log.Warn("No .env file found", zap.Error(err))
	}
	return Config{
		Port: getEnv("PORT", "2000"),
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     "5432",
		DB_USERNAME: getEnv("DB_USERNAME", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", "password"),
		DB_NAME:     getEnv("DB_NAME", "postgres"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
