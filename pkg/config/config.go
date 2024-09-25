package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ZincHost      string
	ZincUser      string
	ZincPassword  string
	ZincIndex     string
	BatchSize     int
	EmailDataPath string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = Config{
		ZincHost:      getEnv("ZINC_HOST", "http://localhost:4080"),
		ZincUser:      getEnv("ZINC_USER", "admin"),
		ZincPassword:  getEnv("ZINC_PASSWORD", "Complexpass#123"),
		ZincIndex:     getEnv("ZINC_INDEX", "enronJELM"),
		BatchSize:     getEnvAsInt("BATCH_SIZE", 500),
		EmailDataPath: getEnv("EMAIL_DATA_PATH", "/home/jaimonetti/dev/Golang/enron_data/maildir/"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
