package dto

import "neosim_go/internal/shared/types"

// CreateKepegawaianPegawaiRequest request body untuk membuat KepegawaianPegawai baru
type CreateKepegawaianPegawaiRequest struct {
	UserID *int64 `json:"user_id" validate:"omitempty,gt=0"`

	// Identitas
	NIK          string  `json:"nik" validate:"required,len=16,numeric"`
	IHSNumber    *string `json:"ihs_number" validate:"omitempty,max=50"`
	NomorPegawai string  `json:"nomor_pegawai" validate:"required,min=1,max=30"`
	NamaLengkap  string  `json:"nama_lengkap" validate:"required,min=1,max=150"`

	// Biodata
	JenisKelaminID     int64           `json:"jenis_kelamin_id" validate:"required,gt=0"`
	TanggalLahir       *types.DateOnly `json:"tanggal_lahir" validate:"required" swaggertype:"string" format:"date" example:"1990-01-01"`
	TempatLahir        string          `json:"tempat_lahir" validate:"required,min=1,max=100"`
	GolonganDarahID    *int64          `json:"golongan_darah_id" validate:"omitempty,gt=0"`
	AgamaID            int64           `json:"agama_id" validate:"required,gt=0"`
	StatusPernikahanID int64           `json:"status_pernikahan_id" validate:"required,gt=0"`
	Kewarganegaraan    *string         `json:"kewarganegaraan" validate:"omitempty,oneof=WNI WNA"`

	// Kepegawaian
	TanggalMasuk  *types.DateOnly `json:"tanggal_masuk" validate:"required" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalKeluar *types.DateOnly `json:"tanggal_keluar" swaggertype:"string" format:"date" example:"2026-12-31"`
	JenisID       int64           `json:"jenis_id" validate:"required,gt=0"`
	StatusID      int64           `json:"status_id" validate:"required,gt=0"`
	FotoURL       *string         `json:"foto_url" validate:"omitempty,max=500"`
	IsAktif       *bool           `json:"is_aktif"`
}

// UpdateKepegawaianPegawaiRequest request body untuk update KepegawaianPegawai
type UpdateKepegawaianPegawaiRequest struct {
	UserID *int64 `json:"user_id" validate:"omitempty,gt=0"`

	NIK          *string `json:"nik" validate:"omitempty,len=16,numeric"`
	IHSNumber    *string `json:"ihs_number" validate:"omitempty,max=50"`
	NomorPegawai *string `json:"nomor_pegawai" validate:"omitempty,min=1,max=30"`
	NamaLengkap  *string `json:"nama_lengkap" validate:"omitempty,min=1,max=150"`

	JenisKelaminID     *int64          `json:"jenis_kelamin_id" validate:"omitempty,gt=0"`
	TanggalLahir       *types.DateOnly `json:"tanggal_lahir" swaggertype:"string" format:"date" example:"1990-01-01"`
	TempatLahir        *string         `json:"tempat_lahir" validate:"omitempty,min=1,max=100"`
	GolonganDarahID    *int64          `json:"golongan_darah_id" validate:"omitempty,gt=0"`
	AgamaID            *int64          `json:"agama_id" validate:"omitempty,gt=0"`
	StatusPernikahanID *int64          `json:"status_pernikahan_id" validate:"omitempty,gt=0"`
	Kewarganegaraan    *string         `json:"kewarganegaraan" validate:"omitempty,oneof=WNI WNA"`

	TanggalMasuk  *types.DateOnly `json:"tanggal_masuk" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalKeluar *types.DateOnly `json:"tanggal_keluar" swaggertype:"string" format:"date" example:"2026-12-31"`
	JenisID       *int64          `json:"jenis_id" validate:"omitempty,gt=0"`
	StatusID      *int64          `json:"status_id" validate:"omitempty,gt=0"`
	FotoURL       *string         `json:"foto_url" validate:"omitempty,max=500"`
	IsAktif       *bool           `json:"is_aktif"`
}

// FilterKepegawaianPegawaiRequest query filter untuk list KepegawaianPegawai
type FilterKepegawaianPegawaiRequest struct {
	Name               string `query:"name"` // cari nama_lengkap (ILIKE)
	NIK                string `query:"nik"`
	NomorPegawai       string `query:"nomor_pegawai"`
	JenisKelaminID     *int64 `query:"jenis_kelamin_id"`
	AgamaID            *int64 `query:"agama_id"`
	StatusPernikahanID *int64 `query:"status_pernikahan_id"`
	JenisID            *int64 `query:"jenis_id"`
	StatusID           *int64 `query:"status_id"`
	IsAktif            *bool  `query:"is_aktif"`

	// Sorting
	SortBy    string `query:"sort_by"`    // user_id | nik | ihs_number | nomor_pegawai | nama_lengkap | tanggal_lahir | tempat_lahir | golongan_darah_id | agama_id | status_pernikahan_id | kewarganegaraan | tanggal_masuk | tanggal_keluar | jenis_id | status_id | foto_url | is_aktif | created_at | updated_at
	SortOrder string `query:"sort_order"` // asc | desc
}
