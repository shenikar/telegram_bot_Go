package service

import (
	"time"

	"telegram_bot_go/config"
	"telegram_bot_go/repository"
)

type UserService struct {
	repo *repository.UserRepo
	cfg  *config.Config
}

func NewUserService(repo *repository.UserRepo, cfg *config.Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

func (s *UserService) LimitAttempt(userID int) (bool, error) {
	timeLimit := time.Now().Add(-time.Duration(s.cfg.Period) * time.Hour)
	count, err := s.repo.CountAttempts(userID, timeLimit)
	if err != nil {
		return false, err
	}
	return count >= s.cfg.MaxAttempt, nil
}

func (s *UserService) SaveAttempt(userID int, hash, result string) error {
	return s.repo.SaveAttempt(userID, hash, result)
}
