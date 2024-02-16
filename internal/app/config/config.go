package config

import (
	"os"
	"sync"
	"test-api/internal/db"
	logs "test-api/internal/log"
)

type Config struct {
	Database db.DbConfig
	Logger   logs.ConfigLogger
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	var err error
	once.Do(func() {
		instance, err = loadConfig()
		if err != nil {
			panic(err)
		}
	})

	return instance
}

func loadConfig() (*Config, error) {
	return &Config{
		Database: getDbConfig(),
		Logger:   getLoggerConfig(),
	}, nil
}

func getDbConfig() db.DbConfig {
	return db.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}
}

func getLoggerConfig() logs.ConfigLogger {
	return logs.ConfigLogger{
		EnableConsole: true,
		Level:         "info",
		EnableFile:    false,
		FileLocation:  os.Getenv("LOGFILE"),
	}
}
