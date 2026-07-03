package dto

// CreateArtikelKategoriRequest request body untuk membuat ArtikelKategori baru
type CreateArtikelKategoriRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateArtikelKategoriRequest request body untuk update ArtikelKategori
type UpdateArtikelKategoriRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterArtikelKategoriRequest request body untuk filter ArtikelKategori
type FilterArtikelKategoriRequest struct {
	Name string `query:"name"`
}
