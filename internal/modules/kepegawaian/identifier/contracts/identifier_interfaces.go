package contracts

import (
	"context"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianIdentifierRepository defines database operations for KepegawaianIdentifier.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/identifier_repository.go).
type KepegawaianIdentifierRepository interface {
	// Create menyimpan identifier baru
	CreateIdentifier(ctx context.Context, m *models.KepegawaianIdentifier) error

	// GetIdentifierByID mencari identifier berdasarkan ID
	GetIdentifierByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error)

	// ListIdentifier mencari semua identifier dengan filter dan pagination
	ListIdentifier(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest) ([]models.KepegawaianIdentifier, int64, error)

	// FindByPegawaiID mencari semua identifier milik satu pegawai
	FindByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianIdentifier, int64, error)

	// FindByPegawaiIDAndTipe mencari identifier milik pegawai untuk tipe tertentu
	FindByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianIdentifier, error)

	// FindPrimaryByTipe mencari identifier primary milik pegawai untuk tipe tertentu
	FindPrimaryByTipe(ctx context.Context, pegawaiID, tipeID int64) (*models.KepegawaianIdentifier, error)

	// Update memperbarui identifier berdasarkan ID
	UpdateIdentifier(ctx context.Context, m *models.KepegawaianIdentifier) error

	// Delete soft delete identifier berdasarkan ID
	DeleteIdentifier(ctx context.Context, id int64, deletedBy int64) error

	// FindExpiringSoonIdentifier mencari identifier yang akan expired dalam N hari ke depan
	FindExpiringSoonIdentifier(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error)

	// FindExpiredIdentifier mencari identifier yang sudah expired
	FindExpiredIdentifier(ctx context.Context) ([]models.KepegawaianIdentifier, error)

	// ExistsIdentifierByID mengecek keberadaan identifier berdasarkan ID
	ExistsIdentifierByID(ctx context.Context, id int64) (bool, error)

	// ExistsByNilaiAndTipe mengecek duplikasi nilai identifier untuk tipe tertentu,
	// excludeID > 0 mengecualikan record itu sendiri (dipakai saat update)
	ExistsByNilaiAndTipe(ctx context.Context, tipeID int64, nilai string, excludeID int64) (bool, error)

	// UnsetPrimaryByPegawaiIDAndTipe melepas status primary identifier lama
	// milik pegawai untuk tipe tertentu, sebelum identifier baru dijadikan primary
	UnsetPrimaryByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64, updatedBy int64) error
}

// KepegawaianIdentifierService defines business logic operations for KepegawaianIdentifier.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/identifier_service.go).
type KepegawaianIdentifierService interface {
	CreateIdentifier(ctx context.Context, req *dto.CreateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	GetIdentifierByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	ListIdentifier(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, int64, error)
	ListByPegawai(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, int64, error)
	UpdateIdentifier(ctx context.Context, id int64, req *dto.UpdateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	DeleteIdentifier(ctx context.Context, id int64, actor he.AuthContext) error

	GetExpiringSoonIdentifier(ctx context.Context, days int, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, error)
	GetExpiredIdentifier(ctx context.Context, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, error)
}
