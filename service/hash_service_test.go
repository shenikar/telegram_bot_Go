package service

import (
	"testing"
)

func TestTypicalGetWordCases(t *testing.T) {
	service := NewHashService()

	word := "ab"
	hash := service.hashingWord(word)
	foundWord, found := service.GetWord(hash)
	if foundWord != word || !found {
		t.Errorf("Expected word %v to be found with hash %v", word, hash)
	}

	wordSpecChars := "a@"
	hashSpecChars := service.hashingWord(wordSpecChars)
	foundWordSpecChars, foundSpec := service.GetWord(hashSpecChars)
	if foundWordSpecChars != wordSpecChars || !found {
		t.Errorf("GetWord(%s) = %v, %v; expected %v, true", hashSpecChars, foundWordSpecChars, foundSpec, wordSpecChars)
	}
}

// краевые случаи метода GetWord
func TestGetWordCases(t *testing.T) {
	service := NewHashService()
	testCases := []struct {
		word     string
		expected string
		found    bool
	}{
		{"", "", false},
		{"~", "~", true},
		{"longwordlimit", "", false},
		{"ñ", "", false},
		{"nonehash", "", false},
	}
	for _, tc := range testCases {
		hash := service.hashingWord(tc.word)
		word, found := service.GetWord(hash)
		if word != tc.expected || found != tc.found {
			t.Errorf("GetWord(%s) = %v, %v; expected %v,%v", hash, word, found, tc.expected, tc.found)
		}
	}
}

// краевые случаи метода hashingWord
func TestHashingWordCases(t *testing.T) {
	service := NewHashService()
	testCases := []struct {
		word     string
		expected string
	}{
		{"", service.hashingWord("")},
		{"~", service.hashingWord("~")},
		{"longwordlimit", service.hashingWord("longwordlimit")},
		{"ñ", service.hashingWord("ñ")},
	}
	for _, tc := range testCases {
		hash := service.hashingWord(tc.word)
		if hash != tc.expected {
			t.Errorf("hashingWord(%s) = %v; expected %v", tc.word, hash, tc.expected)
		}
	}
}

// краевые случаи метода generateWord
func TestGenerateWordsCases(t *testing.T) {
	service := NewHashService()
	words := service.generateWords()
	if contains(words, "") {
		t.Errorf("generateWords() = %v; expected not to contain empty string", words)
	}
	if !contains(words, "~") {
		t.Errorf("generateWords() = %v; expected to contain '~'", words)
	}
	if contains(words, "ñ") {
		t.Errorf("generateWords() = %v; expected not to contain 'ñ'", words)
	}
}

// проверка существует ли слово в сгенерированном списке слов
func contains(s []string, item string) bool {
	for _, i := range s {
		if i == item {
			return true
		}
	}
	return false
}
