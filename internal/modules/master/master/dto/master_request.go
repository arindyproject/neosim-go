package dto

// =========================================================================
// 1. PEKERJAAN
// =========================================================================

// CreateMasterPekerjaanRequest request body untuk membuat Master Pekerjaan baru
type CreateMasterPekerjaanRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterPekerjaanRequest request body untuk update Master Pekerjaan
type UpdateMasterPekerjaanRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterPekerjaanRequest request body untuk filter Master Pekerjaan
type FilterMasterPekerjaanRequest struct {
	Name         string `query:"name"`
	KodeKemenkes string `query:"kode_kemenkes"`
}

// =========================================================================
// 2. PENDIDIKAN
// =========================================================================

// CreateMasterPendidikanRequest request body untuk membuat Master Pendidikan baru
type CreateMasterPendidikanRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterPendidikanRequest request body untuk update Master Pendidikan
type UpdateMasterPendidikanRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterPendidikanRequest request body untuk filter Master Pendidikan
type FilterMasterPendidikanRequest struct {
	Name         string `query:"name"`
	KodeKemenkes string `query:"kode_kemenkes"`
}

// =========================================================================
// 3. AGAMA
// =========================================================================

// CreateMasterAgamaRequest request body untuk membuat Master Agama baru
type CreateMasterAgamaRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterAgamaRequest request body untuk update Master Agama
type UpdateMasterAgamaRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterAgamaRequest request body untuk filter Master Agama
type FilterMasterAgamaRequest struct {
	Name         string `query:"name"`
	KodeKemenkes string `query:"kode_kemenkes"`
}

// =========================================================================
// 4. STATUS PERNIKAHAN
// =========================================================================

// CreateMasterStatusPernikahanRequest request body untuk membuat Master Status Pernikahan baru
type CreateMasterStatusPernikahanRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterStatusPernikahanRequest request body untuk update Master Status Pernikahan
type UpdateMasterStatusPernikahanRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=255"`
	KodeKemenkes *string `json:"kode_kemenkes" validate:"omitempty,max=50"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterStatusPernikahanRequest request body untuk filter Master Status Pernikahan
type FilterMasterStatusPernikahanRequest struct {
	Name         string `query:"name"`
	KodeKemenkes string `query:"kode_kemenkes"`
}
