package service

func (s *HashService) GetWordMulti(hash string) (string, bool) {
	wordChan := make(chan string)
	done := make(chan struct{})

	go s.generateWordsMulti("", 4, wordChan, done)

	for word := range wordChan {
		if s.hashingWord(word) == hash {
			close(done)
			return word, true
		}
	}
	return "", false
}

func (s *HashService) generateWordsMulti(currentWord string, maxLen int, wordChan chan<- string, done <-chan struct{}) {
	if len(currentWord) == maxLen {
		select {
		case <-done:
			return
		case wordChan <- currentWord:
		}
		return
	}

	for _, char := range chars {
		select {
		case <-done:
			return
		default:
			nextWord := currentWord + string(char)
			go s.generateWordsMulti(nextWord, maxLen, wordChan, done)
		}
	}
}
