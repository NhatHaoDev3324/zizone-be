package tdo

type Meta struct {
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

func NewMetaResponse(total, totalPage, page, limit int) Meta {
	return Meta{
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
	}
}
