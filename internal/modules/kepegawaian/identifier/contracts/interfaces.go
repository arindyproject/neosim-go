package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	// Create menyimpan identifier baru
	Create(ctx context.Context, identifier *models.KepegawaianIdentifier) error

	// Update memperbarui identifier berdasarkan ID
	Update(ctx context.Context, identifier *models.KepegawaianIdentifier) error

	// Delete soft delete identifier berdasarkan ID
	Delete(ctx context.Context, id int64, deletedBy int64) error

	// FindByID mencari identifier berdasarkan ID
	FindByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error)

	// FindAll mencari semua identifier dengan filter dan pagination
	FindAll(ctx context.Context, filter dto.FilterKepegawaianIdentifierRequest, page, limit int) ([]models.KepegawaianIdentifier, int64, error)

	// FindByPegawaiD mencari semua identifier milik satu pegawai
	FindByPegawaiD(ctx context.Context, kepegawaianID int64) ([]models.KepegawaianIdentifier, error)

	// FindByPegawaiDAndTipe mencari identifier berdasarkan pegawai dan tipe
	FindByPegawaiDAndTipe(ctx context.Context, kepegawaianID int64, tipe models.IdentifierType) ([]models.KepegawaianIdentifier, error)

	// FindPrimaryByTipe mencari identifier primary untuk tipe tertentu
	FindPrimaryByTipe(ctx context.Context, kepegawaianID int64, tipe models.IdentifierType) (*models.KepegawaianIdentifier, error)

	// FindExpiringSoon mencari identifier yang akan expired dalam N hari ke depan
	FindExpiringSoon(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error)

	// FindExpired mencari semua identifier yang sudah expired
	FindExpired(ctx context.Context) ([]models.KepegawaianIdentifier, error)

	// ExistsByID mengecek apakah identifier dengan ID tersebut ada
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// ExistsByNilaiAndTipe mengecek apakah nilai+tipe sudah dipakai pegawai lain (duplikasi)
	ExistsByNilaiAndTipe(ctx context.Context, tipe models.IdentifierType, nilai string, excludeID int64) (bool, error)

	// UnsetPrimaryByPegawaiDAndTipe mengubah semua is_primary = false untuk tipe tertentu
	// dipanggil sebelum set identifier baru sebagai primary
	UnsetPrimaryByPegawaiDAndTipe(ctx context.Context, kepegawaianID int64, tipe models.IdentifierType, updatedBy int64) error
}

// Service defines business logic operations
type Service interface {
	Create(ctx context.Context, req *dto.CreateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	GetByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)

	// List — filter pakai value, bukan pointer
	List(ctx context.Context, page, pageSize int, filter dto.FilterKepegawaianIdentifierRequest, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, int64, error)

	ListByPegawai(ctx context.Context, kepegawaianID int64, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	Delete(ctx context.Context, id int64, actor he.AuthContext) error
	GetExpiringSoon(ctx context.Context, days int, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, error)
	GetIdentifierTypes() []dto.IdentifierMetaResponse
}
