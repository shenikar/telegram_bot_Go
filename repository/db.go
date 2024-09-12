package repository

import (
	"fmt"
	"telegram_bot_go/config"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func GetConnect(dbConfig config.Databaseconfig) (*sqlx.DB, error) {
	urlDb := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db, err := sqlx.Connect("pgx", urlDb)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	return db, nil
}
