package service

import (
	"crypto/md5"
	"encoding/hex"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=~"

type HashService struct{}

type HashWorder interface {
	GetWord(hash string) (string, bool)
	GetWordMulti(hash string) (string, bool)
}

func NewHashService() *HashService {
	return &HashService{}
}

func (s *HashService) GetWord(hash string) (string, bool) {
	word := s.generateWords("", 4, hash)
	if word != "" {
		return word, true
	}
	return "", false
}

func (s *HashService) hashingWord(word string) string {
	hash := md5.Sum([]byte(word))
	return hex.EncodeToString(hash[:])
}

func (s *HashService) generateWords(currentWord string, maxLen int, targetHash string) string {
	if s.hashingWord(currentWord) == targetHash {
		return currentWord
	}
	if len(currentWord) == maxLen {
		return ""
	}

	for _, char := range chars {
		nextWord := currentWord + string(char)
		foundWord := s.generateWords(nextWord, maxLen, targetHash)
		if foundWord != "" {
			return foundWord
		}
	}
	return ""

}
