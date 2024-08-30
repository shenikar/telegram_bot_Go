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
	for _, word := range s.generateWords() {
		if s.hashingWord(word) == hash {
			return word, true
		}
	}
	return "", false
}

func (s *HashService) hashingWord(word string) string {
	hash := md5.Sum([]byte(word))
	return hex.EncodeToString(hash[:])
}

func (s *HashService) generateWords() []string {
	var words []string
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=~"
	// если длина пароля 1 символ
	for _, ch := range chars {
		words = append(words, string(ch))
	}
	// если длина пароля 2 символ
	for _, ch1 := range chars {
		for _, ch2 := range chars {
			words = append(words, string(ch1)+string(ch2))
		}
	}
	return words
}
