package word

import (
	"encoding/json"
	"strings"

	"github.com/NhatHaoDev3324/zizone-be/internal/model"
	"github.com/NhatHaoDev3324/zizone-be/internal/tdo"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WordService interface {
	GetListWord(page, limit int, search string) (tdo.Meta, []model.Word, error)
	CreateWord(req *CreateWordRequest) error
	GetWordByID(id string) (*model.Word, error)
	UpdateWord(req *UpdateWordRequest) error
	DeleteWord(id string) error
}

type wordService struct {
	repo WordRepository
}

func NewWordService(repo WordRepository) WordService {
	return &wordService{repo}
}

func (s *wordService) GetListWord(page, limit int, search string) (tdo.Meta, []model.Word, error) {
	words, err := s.repo.FindAll()
	if err != nil {
		return tdo.Meta{}, nil, err
	}

	var result []model.Word
	if search != "" {
		search = strings.ToLower(search)
		for _, word := range words {
			wordStr := strings.ToLower(word.Word)
			wordType := strings.ToLower(word.WordType)
			meaning := strings.ToLower(word.Meaning)
			characters := strings.ToLower(string(word.Characters))
			examples := strings.ToLower(string(word.Examples))
			if strings.Contains(wordStr, search) || strings.Contains(wordType, search) || strings.Contains(meaning, search) || strings.Contains(characters, search) || strings.Contains(examples, search) {
				result = append(result, word)
			}
		}
	} else {
		result = words
	}

	total := len(result)
	totalPage := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		return tdo.NewMetaResponse(total, totalPage, page, limit), []model.Word{}, nil
	}
	if end > total {
		end = total
	}

	return tdo.NewMetaResponse(total, totalPage, page, limit), result[start:end], nil
}

func (s *wordService) CreateWord(req *CreateWordRequest) error {
	charactersData, err := json.Marshal(req.Characters)
	if err != nil {
		return err
	}

	examplesData, err := json.Marshal(req.Examples)
	if err != nil {
		return err
	}

	word := &model.Word{
		ID:         uuid.New(),
		Word:       req.Word,
		Pinyin:     req.Pinyin,
		WordType:   req.WordType,
		Meaning:    req.Meaning,
		MemoryTip:  req.MemoryTip,
		Characters: datatypes.JSON(charactersData),
		Examples:   datatypes.JSON(examplesData),
	}

	return s.repo.Create(word)
}

func (s *wordService) GetWordByID(id string) (*model.Word, error) {
	return s.repo.FindByID(id)
}

func (s *wordService) UpdateWord(req *UpdateWordRequest) error {
	uID, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	existingWord, err := s.repo.FindByID(uID.String())
	if err != nil {
		return err
	}

	charactersData, err := json.Marshal(req.Characters)
	if err != nil {
		return err
	}

	examplesData, err := json.Marshal(req.Examples)
	if err != nil {
		return err
	}

	existingWord.Word = req.Word
	existingWord.Pinyin = req.Pinyin
	existingWord.WordType = req.WordType
	existingWord.Meaning = req.Meaning
	existingWord.MemoryTip = req.MemoryTip
	existingWord.Characters = datatypes.JSON(charactersData)
	existingWord.Examples = datatypes.JSON(examplesData)

	return s.repo.Update(existingWord)
}

func (s *wordService) DeleteWord(id string) error {
	return s.repo.Delete(id)
}
