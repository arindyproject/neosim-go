package services

import (
	"net/http"
	"time"

	"neosim_go/internal/modules/master/alamat/contracts"
	alamatContracts "neosim_go/internal/modules/master/alamat/contracts"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo     alamatContracts.Repository
	rbacRepo rbacContracts.RBACRepository //RBAC
	authRepo authContracts.AuthRepository //AUTH
}

// NewMasterAlamatService membuat instance service baru
func NewMasterAlamatService(
	repo alamatContracts.Repository,
	rbacRepo rbacContracts.RBACRepository, //RBAC
	authRepo authContracts.AuthRepository, //AUTH
) alamatContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo, //RBAC
		authRepo: authRepo, //AUTH
	}
}

// ─── Permission ─────────────────────────────────────────────────────────────────
func (s *service) canCreateMasterAlamat(actor alamatContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canUpdateMasterAlamat(actor alamatContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canDeleteMasterAlamat(actor alamatContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}

// ─── Service ────────────────────────────────────────────────────────────────────

// Negara =========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDNegara(id int64) (*dto.NegaraResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDNegara(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}
	return dto.ToNegaraResponse(m), nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListNegara(page, pageSize int, filter *dto.FilterNegaraRequest) ([]dto.NegaraResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListNegara(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	if len(items) == 0 {
		// Mengembalikan error spesifik bahwa data tidak ditemukan
		return nil, 0, appErrors.NotFound("Negara tidak ditemukan")
	}
	return dto.ToNegaraListResponse(items), total, nil
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
func (s *service) CreateNegara(req *dto.CreateNegaraRequest, actor alamatContracts.AuthContext) (*dto.NegaraResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat MasterAlamat baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m := &models.MasterAlamatNegara{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateNegara(m); err != nil {
		return nil, err
	}
	return dto.ToNegaraResponse(m), nil
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
func (s *service) UpdateNegara(id int64, req *dto.UpdateNegaraRequest, actor alamatContracts.AuthContext) (*dto.NegaraResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah MasterAlamat.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDNegara(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateNegara(m); err != nil {
		return nil, err
	}
	return dto.ToNegaraResponse(m), nil
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
func (s *service) DeleteNegara(id int64, actor alamatContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus MasterAlamat.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDNegara(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Negara tidak ditemukan")
	}
	return s.repo.DeleteNegara(id)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Negara =========================================================================

// Provinsi ========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDProvinsi(id int64) (*dto.ProvinsiDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	totalKota, err := s.repo.CountKotaByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}
	totalKecamatan, err := s.repo.CountKecamatanByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}
	totalDesa, err := s.repo.CountDesaByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ProvinsiDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		NegaraID:       m.NegaraID,
		NegaraName:     m.Negara.Name,
		TotalKota:      totalKota,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListProvinsi(page, pageSize, negaraID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Provinsi tidak ditemukan")
	}
	return dto.ToProvinsiListResponse(items), total, nil
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
func (s *service) CreateProvinsi(req *dto.CreateProvinsiRequest, actor alamatContracts.AuthContext) (*dto.ProvinsiResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Provinsi baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	negara, err := s.repo.GetByIDNegara(req.NegaraID)
	if err != nil {
		return nil, err
	}
	if negara == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	m := &models.MasterAlamatProvinsi{
		NegaraID:  req.NegaraID,
		Code:      req.Code,
		Name:      req.Name,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}
	if err := s.repo.CreateProvinsi(m); err != nil {
		return nil, err
	}
	return dto.ToProvinsiResponse(m), nil
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
func (s *service) UpdateProvinsi(id int64, req *dto.UpdateProvinsiRequest, actor alamatContracts.AuthContext) (*dto.ProvinsiResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Provinsi.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	if req.NegaraID != nil {
		m.NegaraID = *req.NegaraID
	}
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateProvinsi(m); err != nil {
		return nil, err
	}
	return dto.ToProvinsiResponse(m), nil
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
func (s *service) DeleteProvinsi(id int64, actor alamatContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Provinsi.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Provinsi tidak ditemukan")
	}
	return s.repo.DeleteProvinsi(id)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Provinsi =======================================================================

// Kota/Kabupaten =====================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKotaKabupaten(id int64) (*dto.KotaKabupatenDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	totalKecamatan, err := s.repo.CountKecamatanByKotaID(m.ID)
	if err != nil {
		return nil, err
	}
	totalDesa, err := s.repo.CountDesaByKotaID(m.ID)
	if err != nil {
		return nil, err
	}

	return &dto.KotaKabupatenDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		ProvinsiID:     m.ProvinsiID,
		ProvinsiName:   m.Provinsi.Name,
		NegaraID:       m.Provinsi.NegaraID,
		NegaraName:     m.Provinsi.Negara.Name,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]dto.KotaKabupatenResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListKotaKabupaten(page, pageSize, provinsiID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}
	return dto.ToKotaKabupatenListResponse(items), total, nil
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
func (s *service) CreateKotaKabupaten(req *dto.CreateKotaKabupatenRequest, actor alamatContracts.AuthContext) (*dto.KotaKabupatenResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kota/Kabupaten baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	provinsi, err := s.repo.GetByIDProvinsi(req.ProvinsiID)
	if err != nil {
		return nil, err
	}
	if provinsi == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	m := &models.MasterAlamatKotaKabupaten{
		ProvinsiID: req.ProvinsiID,
		Code:       req.Code,
		Name:       req.Name,
		CreatedBy:  &actor.UserID,
		UpdatedBy:  &actor.UserID,
	}
	if err := s.repo.CreateKotaKabupaten(m); err != nil {
		return nil, err
	}
	return dto.ToKotaKabupatenResponse(m), nil
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
func (s *service) UpdateKotaKabupaten(id int64, req *dto.UpdateKotaKabupatenRequest, actor alamatContracts.AuthContext) (*dto.KotaKabupatenResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kota/Kabupaten.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	if req.ProvinsiID != nil {
		m.ProvinsiID = *req.ProvinsiID
	}
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKotaKabupaten(m); err != nil {
		return nil, err
	}
	return dto.ToKotaKabupatenResponse(m), nil
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
func (s *service) DeleteKotaKabupaten(id int64, actor alamatContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kota/Kabupaten.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}
	return s.repo.DeleteKotaKabupaten(id)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kota/Kabupaten =====================================================================

// Kecamatan ===========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKecamatan(id int64) (*dto.KecamatanDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	totalDesa, err := s.repo.CountDesaByKecamatanID(m.ID)
	if err != nil {
		return nil, err
	}

	return &dto.KecamatanDetailResponse{
		ID:                m.ID,
		Code:              m.Code,
		Name:              m.Name,
		KotaKabupatenID:   m.KotaKabupatenID,
		KotaKabupatenName: m.KotaKabupaten.Name,
		ProvinsiID:        m.KotaKabupaten.ProvinsiID,
		ProvinsiName:      m.KotaKabupaten.Provinsi.Name,
		NegaraID:          m.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        m.KotaKabupaten.Provinsi.Negara.Name,
		TotalDesa:         totalDesa,
	}, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListKecamatan(page, pageSize, kotaKabupatenID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kecamatan tidak ditemukan")
	}
	return dto.ToKecamatanListResponse(items), total, nil
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
func (s *service) CreateKecamatan(req *dto.CreateKecamatanRequest, actor alamatContracts.AuthContext) (*dto.KecamatanResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kecamatan baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	kota, err := s.repo.GetByIDKotaKabupaten(req.KotaKabupatenID)
	if err != nil {
		return nil, err
	}
	if kota == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	m := &models.MasterAlamatKecamatan{
		KotaKabupatenID: req.KotaKabupatenID,
		Code:            req.Code,
		Name:            req.Name,
		CreatedBy:       &actor.UserID,
		UpdatedBy:       &actor.UserID,
	}
	if err := s.repo.CreateKecamatan(m); err != nil {
		return nil, err
	}
	return dto.ToKecamatanResponse(m), nil
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
func (s *service) UpdateKecamatan(id int64, req *dto.UpdateKecamatanRequest, actor alamatContracts.AuthContext) (*dto.KecamatanResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kecamatan.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	if req.KotaKabupatenID != nil {
		m.KotaKabupatenID = *req.KotaKabupatenID
	}
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKecamatan(m); err != nil {
		return nil, err
	}
	return dto.ToKecamatanResponse(m), nil
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
func (s *service) DeleteKecamatan(id int64, actor alamatContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kecamatan.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kecamatan tidak ditemukan")
	}
	return s.repo.DeleteKecamatan(id)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kecamatan ===========================================================================

// Kelurahan/Desa ========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKelurahanDesa(id int64) (*dto.KelurahanDesaDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	return &dto.KelurahanDesaDetailResponse{
		ID:                m.ID,
		Code:              m.Code,
		Name:              m.Name,
		PostalCode:        m.PostalCode,
		KecamatanID:       m.KecamatanID,
		KecamatanName:     m.Kecamatan.Name,
		KotaKabupatenID:   m.Kecamatan.KotaKabupatenID,
		KotaKabupatenName: m.Kecamatan.KotaKabupaten.Name,
		ProvinsiID:        m.Kecamatan.KotaKabupaten.ProvinsiID,
		ProvinsiName:      m.Kecamatan.KotaKabupaten.Provinsi.Name,
		NegaraID:          m.Kecamatan.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        m.Kecamatan.KotaKabupaten.Provinsi.Negara.Name,
	}, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListKelurahanDesa(page, pageSize, kecamatanID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}
	return dto.ToKelurahanDesaListResponse(items), total, nil
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
func (s *service) CreateKelurahanDesa(req *dto.CreateKelurahanDesaRequest, actor alamatContracts.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kelurahan/Desa baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	kecamatan, err := s.repo.GetByIDKecamatan(req.KecamatanID)
	if err != nil {
		return nil, err
	}
	if kecamatan == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	m := &models.MasterAlamatKelurahanDesa{
		KecamatanID: req.KecamatanID,
		Code:        req.Code,
		Name:        req.Name,
		PostalCode:  req.PostalCode,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKelurahanDesa(m); err != nil {
		return nil, err
	}
	return dto.ToKelurahanDesaResponse(m), nil
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
func (s *service) UpdateKelurahanDesa(id int64, req *dto.UpdateKelurahanDesaRequest, actor alamatContracts.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kelurahan/Desa.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	if req.KecamatanID != nil {
		m.KecamatanID = *req.KecamatanID
	}
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.PostalCode != nil {
		m.PostalCode = req.PostalCode
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKelurahanDesa(m); err != nil {
		return nil, err
	}
	return dto.ToKelurahanDesaResponse(m), nil
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
func (s *service) DeleteKelurahanDesa(id int64, actor alamatContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kelurahan/Desa.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}
	return s.repo.DeleteKelurahanDesa(id)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kelurahan/Desa ========================================================================

func (s *service) GetDetailProvinsi(id int64, actor contracts.AuthContext) (*dto.ProvinsiDetailResponse, error) {
	provinsi, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return nil, err
	}
	if provinsi == nil {
		return nil, appErrors.BadRequest("")
	}

	totalKota, err := s.repo.CountKotaByProvinsiID(provinsi.ID)
	if err != nil {
		return nil, err
	}

	//totalKecamatan, err := s.repo.CountKecamatanByProvinsiID(provinsi.ID)
	//if err != nil {
	//	return nil, err
	//}

	totalDesa, err := s.repo.CountDesaByProvinsiID(provinsi.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ProvinsiDetailResponse{
		ID:         provinsi.ID,
		Code:       provinsi.Code,
		Name:       provinsi.Name,
		NegaraID:   provinsi.NegaraID,
		NegaraName: provinsi.Negara.Name,
		TotalKota:  totalKota,
		//TotalKecamatan: totalKecamatan,
		TotalDesa: totalDesa,
	}, nil
}

func (s *service) GetDetailKotaKabupaten(id int64, actor contracts.AuthContext) (*dto.KotaKabupatenDetailResponse, error) {
	kota, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return nil, err
	}
	if kota == nil {
		return nil, appErrors.BadRequest("")
	}

	totalKecamatan, err := s.repo.CountKecamatanByKotaID(kota.ID)
	if err != nil {
		return nil, err
	}

	totalDesa, err := s.repo.CountDesaByKotaID(kota.ID)
	if err != nil {
		return nil, err
	}

	return &dto.KotaKabupatenDetailResponse{
		ID:             kota.ID,
		Code:           kota.Code,
		Name:           kota.Name,
		ProvinsiID:     kota.ProvinsiID,
		ProvinsiName:   kota.Provinsi.Name,
		NegaraID:       kota.Provinsi.NegaraID,
		NegaraName:     kota.Provinsi.Negara.Name,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}, nil
}

func (s *service) GetDetailKecamatan(id int64, actor contracts.AuthContext) (*dto.KecamatanDetailResponse, error) {
	kecamatan, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return nil, err
	}
	if kecamatan == nil {
		return nil, appErrors.BadRequest("")
	}

	totalDesa, err := s.repo.CountDesaByKecamatanID(kecamatan.ID)
	if err != nil {
		return nil, err
	}

	return &dto.KecamatanDetailResponse{
		ID:                kecamatan.ID,
		Code:              kecamatan.Code,
		Name:              kecamatan.Name,
		KotaKabupatenID:   kecamatan.KotaKabupatenID,
		KotaKabupatenName: kecamatan.KotaKabupaten.Name,
		ProvinsiID:        kecamatan.KotaKabupaten.ProvinsiID,
		ProvinsiName:      kecamatan.KotaKabupaten.Provinsi.Name,
		NegaraID:          kecamatan.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        kecamatan.KotaKabupaten.Provinsi.Negara.Name,
		TotalDesa:         totalDesa,
	}, nil
}

func (s *service) GetDetailKelurahanDesa(id int64, actor contracts.AuthContext) (*dto.KelurahanDesaDetailResponse, error) {
	desa, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return nil, err
	}
	if desa == nil {
		return nil, appErrors.BadRequest("")
	}

	return &dto.KelurahanDesaDetailResponse{
		ID:                desa.ID,
		Code:              desa.Code,
		Name:              desa.Name,
		PostalCode:        desa.PostalCode,
		KecamatanID:       desa.KecamatanID,
		KecamatanName:     desa.Kecamatan.Name,
		KotaKabupatenID:   desa.Kecamatan.KotaKabupatenID,
		KotaKabupatenName: desa.Kecamatan.KotaKabupaten.Name,
		ProvinsiID:        desa.Kecamatan.KotaKabupaten.ProvinsiID,
		ProvinsiName:      desa.Kecamatan.KotaKabupaten.Provinsi.Name,
		NegaraID:          desa.Kecamatan.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        desa.Kecamatan.KotaKabupaten.Provinsi.Negara.Name,
	}, nil
}

// ------------ Create ------------------------------------------------------------

// ------------ List --------------------------------------------------------------

// ------------ Update ------------------------------------------------------------

// ------------ Delete ------------------------------------------------------------
