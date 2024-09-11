package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   Databaseconfig
	Telegram   TelegramConfig
	MaxAttempt int
	Period     string
}

type Databaseconfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

type TelegramConfig struct {
	Token   string
	Timeout int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbConfig := Databaseconfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}

	maxAttemptStr := os.Getenv("MAX_ATTEMPTS")
	maxAttempt, err := strconv.Atoi(maxAttemptStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing MAX_ATTEMPTS: %v", err)
	}

	period := os.Getenv("PERIOD_ATTEMPTS")
	if period == "" {
		return nil, fmt.Errorf("PERIOD_ATTEMPTS is not set")
	}

	timeoutStr := os.Getenv("TIMEOUT_BOT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing TIMEOUT_BOT: %v", err)
	}

	telegramConfig := TelegramConfig{
		Token:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		Timeout: timeout,
	}

	return &Config{
		Database:   dbConfig,
		Telegram:   telegramConfig,
		MaxAttempt: maxAttempt,
		Period:     period,
	}, nil
}
