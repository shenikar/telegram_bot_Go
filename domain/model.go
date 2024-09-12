package domain

import "time"

type User struct {
	ID             int
	AttemptHistory []HashRequest
}

type HashRequest struct {
	Hash        string    `db:"hash"`
	AttemptTime time.Time `db:"created_at"`
}
