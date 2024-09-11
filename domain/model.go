package domain

import "time"

type User struct {
	ID             int
	AttemptLimit   int
	AttemptHistory []HashRequest
}

type HashRequest struct {
	Hash        string
	AttemptTime time.Time
}
