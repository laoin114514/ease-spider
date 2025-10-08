package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
}

var AppConfig *Config

func Init() error {
	godotenv.Load()
	AppConfig = &Config{
		Port:        os.Getenv("PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
	}
	return nil
}
