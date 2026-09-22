package dto

// CreatePositionRequest request body untuk membuat Position baru
type CreatePositionRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`

	PositionKategoriID int64 `json:"position_kategori_id" validate:"required,gt=0"`

	// ParentID kosong = puncak hierarki (mis. Direktur).
	ParentID *int64 `json:"parent_id" validate:"omitempty,gt=0"`

	// DepartmentID kosong = jabatan tidak terikat 1 unit (mis. Direktur, Wakil Direktur).
	DepartmentID *int64 `json:"department_id" validate:"omitempty,gt=0"`

	LevelHierarki int16  `json:"level_hierarki" validate:"required,gt=0"`
	Kuota         *int16 `json:"kuota" validate:"omitempty,gt=0"`
	IsAktif       bool   `json:"is_aktif"`
}

// UpdatePositionRequest request body untuk update Position
type UpdatePositionRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`

	PositionKategoriID *int64 `json:"position_kategori_id" validate:"omitempty,gt=0"`
	ParentID           *int64 `json:"parent_id" validate:"omitempty,gt=0"`
	DepartmentID       *int64 `json:"department_id" validate:"omitempty,gt=0"`

	LevelHierarki *int16 `json:"level_hierarki" validate:"omitempty,gt=0"`
	Kuota         *int16 `json:"kuota" validate:"omitempty,gt=0"`
	IsAktif       *bool  `json:"is_aktif"`
}

// FilterPositionRequest request query untuk filter Position
type FilterPositionRequest struct {
	Name               string `query:"name"`
	PositionKategoriID *int64 `query:"position_kategori_id"`
	ParentID           *int64 `query:"parent_id"`
	DepartmentID       *int64 `query:"department_id"`
	IsAktif            *bool  `query:"is_aktif"`

	// IsRoot true → hanya posisi dengan parent_id NULL (puncak bagan organisasi).
	IsRoot *bool `query:"is_root"`
}
