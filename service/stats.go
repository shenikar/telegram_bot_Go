package service

import (
	"log"
	"telegram_bot_go/domain"
	"telegram_bot_go/repository"
	"time"
)

type StatsService struct {
	userRepo *repository.UserRepo
}

func NewStatsService(userRepo *repository.UserRepo) *StatsService {
	return &StatsService{userRepo: userRepo}
}

func (s *StatsService) GetStats(userID int) ([]domain.HashRequest, error) {
	timeLimit := time.Now().Add(-24 * time.Hour)
	requests, err := s.userRepo.GetAttemptHistory(userID, timeLimit)
	if err != nil {
		log.Printf("Error retrieving requests: %v", err)
		return nil, err
	}
	return requests, nil
}

func (s *StatsService) GetLast24HoursRequests() ([]domain.HashRequest, error) {
	timeLimit := time.Now().Add(-24 * time.Hour)
	return s.userRepo.GetRequestsLast24Hours(timeLimit)
}
