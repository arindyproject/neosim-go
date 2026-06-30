package dto

// CreateArtikelRequest request body untuk membuat Artikel baru
type CreateArtikelRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateArtikelRequest request body untuk update Artikel
type UpdateArtikelRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterArtikelRequest request body untuk filter Artikel
type FilterArtikelRequest struct {
	Name string `query:"name"`
}
