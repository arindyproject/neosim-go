package dto

// CreateArtikelKetegoriRequest request body untuk membuat ArtikelKetegori baru
type CreateArtikelKetegoriRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateArtikelKetegoriRequest request body untuk update ArtikelKetegori
type UpdateArtikelKetegoriRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterArtikelKetegoriRequest request body untuk filter ArtikelKetegori
type FilterArtikelKetegoriRequest struct {
	Name        string `query:"name"`
}

