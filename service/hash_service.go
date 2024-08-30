package service

import (
	"crypto/md5"
	"encoding/hex"
)

type HashService struct{}

func NewHashService() *HashService {
	return &HashService{}
}

func (s *HashService) GetWord(hash string) (string, bool) {
	return "", false
}

func (s *HashService) hashingWord(word string) string {
	hash := md5.Sum([]byte(word))
	return hex.EncodeToString(hash[:])
}
