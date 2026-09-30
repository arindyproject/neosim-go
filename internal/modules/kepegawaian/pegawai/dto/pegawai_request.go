package dto

// CreateKepegawaianPegawaiRequest request body untuk membuat KepegawaianPegawai baru
type CreateKepegawaianPegawaiRequest struct {
	UserID *int64 `json:"user_id" validate:"omitempty,gt=0"`

	// Identitas
	NIK          string  `json:"nik" validate:"required,len=16,numeric"`
	IHSNumber    *string `json:"ihs_number" validate:"omitempty,max=50"`
	NomorPegawai string  `json:"nomor_pegawai" validate:"required,min=1,max=30"`
	NamaLengkap  string  `json:"nama_lengkap" validate:"required,min=1,max=150"`

	// Biodata
	JenisKelamin     string  `json:"jenis_kelamin" validate:"required,oneof=Laki-laki Perempuan"`
	TanggalLahir     string  `json:"tanggal_lahir" validate:"required,datetime=2006-01-02"`
	TempatLahir      string  `json:"tempat_lahir" validate:"required,min=1,max=100"`
	GolonganDarah    *string `json:"golongan_darah" validate:"omitempty,oneof=A B AB O A+ A- B+ B- AB+ AB- O+ O-"`
	Agama            string  `json:"agama" validate:"required,oneof=Islam Kristen Katolik Hindu Buddha Konghucu Lainnya"`
	StatusPerkawinan string  `json:"status_perkawinan" validate:"required,oneof='Belum Kawin' Kawin 'Cerai Hidup' 'Cerai Mati'"`
	Kewarganegaraan  string  `json:"kewarganegaraan" validate:"required,min=1,max=5"`

	// Kepegawaian
	TanggalMasuk  string  `json:"tanggal_masuk" validate:"required,datetime=2006-01-02"`
	TanggalKeluar *string `json:"tanggal_keluar" validate:"omitempty,datetime=2006-01-02"`
	JenisID       int64   `json:"jenis_id" validate:"required,gt=0"`
	StatusID      int64   `json:"status_id" validate:"required,gt=0"`
	FotoURL       *string `json:"foto_url" validate:"omitempty,max=500"`
	IsAktif       *bool   `json:"is_aktif"`
}

// UpdateKepegawaianPegawaiRequest request body untuk update KepegawaianPegawai
type UpdateKepegawaianPegawaiRequest struct {
	UserID *int64 `json:"user_id" validate:"omitempty,gt=0"`

	NIK          *string `json:"nik" validate:"omitempty,len=16,numeric"`
	IHSNumber    *string `json:"ihs_number" validate:"omitempty,max=50"`
	NomorPegawai *string `json:"nomor_pegawai" validate:"omitempty,min=1,max=30"`
	NamaLengkap  *string `json:"nama_lengkap" validate:"omitempty,min=1,max=150"`

	JenisKelamin     *string `json:"jenis_kelamin" validate:"omitempty,oneof=Laki-laki Perempuan"`
	TanggalLahir     *string `json:"tanggal_lahir" validate:"omitempty,datetime=2006-01-02"`
	TempatLahir      *string `json:"tempat_lahir" validate:"omitempty,min=1,max=100"`
	GolonganDarah    *string `json:"golongan_darah" validate:"omitempty,oneof=A B AB O A+ A- B+ B- AB+ AB- O+ O-"`
	Agama            *string `json:"agama" validate:"omitempty,oneof=Islam Kristen Katolik Hindu Buddha Konghucu Lainnya"`
	StatusPerkawinan *string `json:"status_perkawinan" validate:"omitempty,oneof='Belum Kawin' Kawin 'Cerai Hidup' 'Cerai Mati'"`
	Kewarganegaraan  *string `json:"kewarganegaraan" validate:"omitempty,min=1,max=5"`

	TanggalMasuk  *string `json:"tanggal_masuk" validate:"omitempty,datetime=2006-01-02"`
	TanggalKeluar *string `json:"tanggal_keluar" validate:"omitempty,datetime=2006-01-02"`
	JenisID       *int64  `json:"jenis_id" validate:"omitempty,gt=0"`
	StatusID      *int64  `json:"status_id" validate:"omitempty,gt=0"`
	FotoURL       *string `json:"foto_url" validate:"omitempty,max=500"`
	IsAktif       *bool   `json:"is_aktif"`
}

// FilterKepegawaianPegawaiRequest query filter untuk list KepegawaianPegawai
type FilterKepegawaianPegawaiRequest struct {
	Name         string `query:"name"` // cari nama_lengkap (ILIKE)
	NIK          string `query:"nik"`
	NomorPegawai string `query:"nomor_pegawai"`
	JenisKelamin string `query:"jenis_kelamin"`
	JenisID      *int64 `query:"jenis_id"`
	StatusID     *int64 `query:"status_id"`
	IsAktif      *bool  `query:"is_aktif"`
}
