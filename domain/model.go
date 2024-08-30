package domain

// Интерфейс для работы с хешами
type HashWorder interface {
	GetWord(hash string) (string, bool)
}
