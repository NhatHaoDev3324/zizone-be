package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Word struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Word      string    `json:"word"`
	Pinyin    string    `json:"pinyin"`
	WordType  string    `json:"word_type"`
	Meaning   string    `json:"meaning"`
	MemoryTip string    `json:"memory_tip"`

	Characters datatypes.JSON `gorm:"type:jsonb" json:"characters"`
	Examples   datatypes.JSON `gorm:"type:jsonb" json:"examples"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Character struct {
	Hanzi         string `json:"hanzi"`
	Pinyin        string `json:"pinyin"`
	CharacterType string `json:"character_type"`
	Structure     string `json:"structure"`
	Imagination   string `json:"imagination"`
	Meaning       string `json:"meaning"`
	Position      int    `json:"position"`
}

type Example struct {
	Sentence    string `json:"sentence"`
	Pinyin      string `json:"pinyin"`
	Translation string `json:"translation"`
}
