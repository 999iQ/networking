package config

import (
	"fmt"
	"os"
)

type DbConnectConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Dbname      string
	Sslmode     string
	AutoMigrate bool
}

func LoadConfig() (DbConnectConfig, error) {
	config := DbConnectConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		Dbname:      os.Getenv("DB_NAME"),
		Sslmode:     os.Getenv("DB_SSLMODE"),
		AutoMigrate: os.Getenv("DB_AUTO_MIGRATE") == "true",
	}

	// check variables
	if config.Host == "" || config.Port == "" || config.User == "" || config.Password == "" ||
		config.Dbname == "" || config.Sslmode == "" {
		return DbConnectConfig{}, fmt.Errorf("необходимо установить все переменные окружения для базы данных: " +
			"DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE")
	}
	return config, nil
}
