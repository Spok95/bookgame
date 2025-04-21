package game

import (
	"encoding/json"
	"os"
)

// Paragraph описывает один параграф книги
type Paragraph struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

// Story хранит карту всех параграфов
type Story struct {
	Paragraphs map[string]Paragraph
}

// LoadStory загружает story.json в память
func LoadStory(filename string) (*Story, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paragraphs map[string]Paragraph
	err = json.NewDecoder(file).Decode(&paragraphs)
	if err != nil {
		return nil, err
	}
	return &Story{Paragraphs: paragraphs}, nil
}

// Get возвращает текст нужного параграфа
func (s *Story) Get(num string) (Paragraph, bool) {
	para, ok := s.Paragraphs[num]
	return para, ok
}
