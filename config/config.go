package config

import (
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var ProjectConfig schemas.Config

func createConfig() schemas.Config {
	return schemas.Config{
		DB: schemas.DatabaseConfig{
			User:     getStrEnv("DB_USER_NAME", ""),
			Password: getStrEnv("DB_PASSWORD", ""),
			Host:     getStrEnv("DB_HOST", ""),
			Port:     getStrEnv("DB_PORT", ""),
			DBName:   getStrEnv("DB_NAME", ""),
		},
		Deploy: schemas.Deploy{
			Port: getIntEnv("DEPLOY_PORT", 5050),
		},
	}
}

func getStrEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		intVal, _ := strconv.Atoi(value)
		return intVal
	}

	return defaultVal
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func init() {
	loadEnv()
	ProjectConfig = createConfig()
}
