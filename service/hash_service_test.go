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

	// Пустая строка
	hash := service.hashingWord("")
	word, found := service.GetWord(hash)
	if word != "" || found {
		t.Errorf("GetWord(%s) = %v, %v; expected empty string and false", hash, word, found)
	}

	// Специальные символы
	hashSpecChars := service.hashingWord("~")
	wordSpecChars, foundSpecChars := service.GetWord(hashSpecChars)
	if wordSpecChars != "~" || !foundSpecChars {
		t.Errorf("GetWord(%s) = %v, %v; expected '~' and true", hashSpecChars, wordSpecChars, foundSpecChars)
	}

	// Несуществующее слово
	hashNonExistent := service.hashingWord("longwordlimit")
	wordNonExistent, foundNonExistent := service.GetWord(hashNonExistent)
	if wordNonExistent != "" || foundNonExistent {
		t.Errorf("GetWord(%s) = %v, %v; expected empty string and false", hashNonExistent, wordNonExistent, foundNonExistent)
	}

	// Сложные символы
	hashSpecial := service.hashingWord("ñ")
	wordSpecial, foundSpecial := service.GetWord(hashSpecial)
	if wordSpecial != "" || foundSpecial {
		t.Errorf("GetWord(%s) = %v, %v; expected empty string and false", hashSpecial, wordSpecial, foundSpecial)
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
	hash := service.hashingWord("a")
	generatedWord := service.generateWords("", 4, hash)
	if generatedWord != "a" {
		t.Errorf("generateWords() = %v; expected 'a'", generatedWord)
	}

	hash = service.hashingWord("ab")
	generatedWord = service.generateWords("", 4, hash)
	if generatedWord != "ab" {
		t.Errorf("generateWords() = %v; expected 'ab'", generatedWord)
	}

	hash = service.hashingWord("nonexistent")
	generatedWord = service.generateWords("", 4, hash)
	if generatedWord != "" {
		t.Errorf("generateWords() = %v; expected empty string", generatedWord)
	}
}
