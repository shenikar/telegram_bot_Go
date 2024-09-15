package service

import (
	"sync"
)

func (s *HashService) GetWordMulti(hash string) (string, bool) {
	resultChan := make(chan string, 1)
	var wg sync.WaitGroup

	for _, char := range chars {
		wg.Add(1)
		go func(char rune) {
			defer wg.Done()
			s.generateWordsMulti(string(char), 4, hash, resultChan)
		}(char)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	if result, ok := <-resultChan; ok {
		return result, true
	}
	return "", false
}

func (s *HashService) generateWordsMulti(currentWord string, maxLen int, targetHash string, resultChan chan string) {
	if s.hashingWord(currentWord) == targetHash {
		select {
		case resultChan <- currentWord:
		default:
		}
		return
	}
	if len(currentWord) == maxLen {
		return
	}

	for _, char := range chars {
		nextWord := currentWord + string(char)
		s.generateWordsMulti(nextWord, maxLen, targetHash, resultChan)
	}
}
