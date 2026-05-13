package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NhatHaoDev3324/zizone-be/internal/modules/word/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WordRepository interface {
	FindAll() ([]model.Word, error)
	Create(word *model.Word) error
	FindByID(id string) (*model.Word, error)
	Update(word *model.Word) error
	Delete(id string) error
}

type wordRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewWordRepository(db *gorm.DB, redis *redis.Client) WordRepository {
	return &wordRepository{db, redis}
}

func (r *wordRepository) FindAll() ([]model.Word, error) {
	ctx := context.Background()
	var words []model.Word

	cachedWords, err := r.redis.Get(ctx, "words:all").Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedWords), &words) == nil {
			return words, nil
		}
	}

	err = r.db.Model(&model.Word{}).Find(&words).Error
	if err == nil {
		wordsData, _ := json.Marshal(words)
		r.redis.Set(ctx, "words:all", wordsData, 30*time.Minute).Err()
	}

	return words, err
}

func (r *wordRepository) Create(word *model.Word) error {
	ctx := context.Background()
	err := r.db.Create(word).Error
	if err == nil {
		r.redis.Del(ctx, "words:all")
	}
	return err
}

func (r *wordRepository) FindByID(id string) (*model.Word, error) {
	var word model.Word
	err := r.db.First(&word, "id = ?", id).Error
	return &word, err
}

func (r *wordRepository) Update(word *model.Word) error {
	ctx := context.Background()
	err := r.db.Save(word).Error
	if err == nil {
		r.redis.Del(ctx, "words:all")
	}
	return err
}

func (r *wordRepository) Delete(id string) error {
	ctx := context.Background()
	err := r.db.Delete(&model.Word{}, "id = ?", id).Error
	if err == nil {
		r.redis.Del(ctx, "words:all")
	}
	return err
}
