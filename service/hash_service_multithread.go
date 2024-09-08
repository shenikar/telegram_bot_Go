package service

func (s *HashService) GetWordMulti(hash string) (string, bool) {
	wordChan := make(chan string)
	done := make(chan struct{})

	go s.generateWordsMulti(wordChan, done)
	for word := range wordChan {
		if s.hashingWord(word) == hash {
			close(done)
			return word, true
		}
	}
	return "", false
}

func (s *HashService) generateWordsMulti(wordChan chan<- string, done <-chan struct{}) {
	go func() {
		defer close(wordChan)
		// если длина пароля 1 символ

		for _, ch := range chars {
			select {
			case <-done:
				return
			case wordChan <- string(ch):
			}
		}

		// если длина пароля 2 символа
		for _, ch1 := range chars {
			for _, ch2 := range chars {
				select {
				case <-done:
					return
				case wordChan <- string(ch1) + string(ch2):
				}
			}
		}

		// если длина пароля 3 символа
		for _, ch1 := range chars {
			for _, ch2 := range chars {
				for _, ch3 := range chars {
					select {
					case <-done:
						return
					case wordChan <- string(ch1) + string(ch2) + string(ch3):
					}
				}
			}
		}

		// если длина пароля 4 символа
		for _, ch1 := range chars {
			for _, ch2 := range chars {
				for _, ch3 := range chars {
					for _, ch4 := range chars {
						select {
						case <-done:
							return
						case wordChan <- string(ch1) + string(ch2) + string(ch3) + string(ch4):
						}
					}
				}
			}
		}
	}()
}
