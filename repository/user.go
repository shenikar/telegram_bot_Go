package repository

import (
	"telegram_bot_go/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CountAttempts(userID int, timeLimit time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM requests WHERE user_id=$1 AND created_at > $2`
	err := r.db.Get(&count, query, userID, timeLimit)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepo) SaveAttempt(userID int, hash string) error {
	query := `INSERT INTO requests (user_id, hash, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, userID, hash, time.Now())
	return err
}

func (r *UserRepo) GetAttemptHistory(userID int) ([]domain.HashRequest, error) {
	var attempts []domain.HashRequest
	query := `SELECT hash, created_at FROM requests WHERE user_id = $1`
	err := r.db.Select(&attempts, query, userID)
	if err != nil {
		return nil, err
	}
	return attempts, nil
}
