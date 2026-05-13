package dto

type CreateWordRequest struct {
	Word       string      `json:"word" binding:"required"`
	Pinyin     string      `json:"pinyin" binding:"required"`
	WordType   string      `json:"word_type" binding:"required"`
	Meaning    string      `json:"meaning" binding:"required"`
	MemoryTip  string      `json:"memory_tip" binding:"required"`
	Characters []Character `json:"characters"`
	Examples   []Example   `json:"examples" binding:"required"`
}

type Character struct {
	Hanzi         string `json:"hanzi" binding:"required"`
	Pinyin        string `json:"pinyin" binding:"required"`
	CharacterType string `json:"character_type" binding:"required"`
	Structure     string `json:"structure" binding:"required"`
	Imagination   string `json:"imagination" binding:"required"`
	Meaning       string `json:"meaning" binding:"required"`
	Position      int    `json:"position" binding:"required"`
}

type Example struct {
	Sentence    string `json:"sentence" binding:"required"`
	Pinyin      string `json:"pinyin" binding:"required"`
	Translation string `json:"translation" binding:"required"`
}

type UpdateWordRequest struct {
	ID         string      `json:"id" binding:"required"`
	Word       string      `json:"word" binding:"required"`
	Pinyin     string      `json:"pinyin" binding:"required"`
	WordType   string      `json:"word_type" binding:"required"`
	Meaning    string      `json:"meaning" binding:"required"`
	MemoryTip  string      `json:"memory_tip" binding:"required"`
	Characters []Character `json:"characters"`
	Examples   []Example   `json:"examples" binding:"required"`
}
