package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbName == "" || dbPassword == "" || dbUser == "" {
		return nil, fmt.Errorf("some database configuration enviromental variables are missing")
	}

	config := &Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbName,
		DBPassword: dbPassword,
		DBName:     dbUser,
	}

	return config, nil
}
