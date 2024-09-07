package repository

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func GetConnect() (*sqlx.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	urlDb := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sqlx.Connect("pgx", urlDb)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	return db, nil
}
