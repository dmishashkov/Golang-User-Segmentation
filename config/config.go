package config

import (
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"os"
	"strconv"
)

var ProjectConfig schemas.Config

func createConfig() schemas.Config {
	return schemas.Config{
		DB: schemas.DatabaseConfig{
			User:     getStrEnv("DB_USER_NAME", "postgres"),
			Password: getStrEnv("DB_PASSWORD", ""),
			Host:     getStrEnv("DB_HOST", "database"),
			Port:     getStrEnv("DB_DOCKER_PORT", "5432"),
			DBName:   getStrEnv("DB_NAME", "avito2023"),
		},
		Deploy: schemas.Deploy{
			Port: getIntEnv("SERVER_DOCKER_PORT", 5050),
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

func init() {
	ProjectConfig = createConfig()
}
