package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) LimitRequest(userID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM requests WHERE user_id=$1 AND created_at > now() - interval '1 hour'`
	err := r.db.Get(&count, query, userID)
	if err != nil {
		return false, err
	}
	return count >= 10, nil
}

func (r *UserRepo) SaveRequest(userID int) error {
	query := `INSERT INTO requests (user_id, created_at) VALUES ($1, $2)`
	_, err := r.db.Exec(query, userID, time.Now())
	return err
}
