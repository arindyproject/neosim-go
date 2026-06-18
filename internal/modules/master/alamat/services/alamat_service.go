package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	alamatContracts "neosim_go/internal/modules/master/alamat/contracts"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"

	//cache

	"neosim_go/internal/shared/cache"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo     alamatContracts.Repository
	rbacRepo rbacContracts.RBACRepository //RBAC
	authRepo authContracts.AuthRepository //AUTH
	cache    *cache.Manager               // <--- Gunakan Cache Manager
}

// NewMasterAlamatService membuat instance service baru
func NewMasterAlamatService(
	repo alamatContracts.Repository,
	rbacRepo rbacContracts.RBACRepository, //RBAC
	authRepo authContracts.AuthRepository, //AUTH
	cacheManager *cache.Manager, // <--- Terima Cache Manager
) alamatContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo, //RBAC
		authRepo: authRepo, //AUTH
		cache:    cacheManager,
	}
}

// ─── Cache Keys Generator ──────────────────────────────────────────────────────
// Negara-------------------------------------------------------------------------
func cacheKeyNegaraDetail(id int64) string { return fmt.Sprintf("master_alamat:negara:detail:%d", id) }
func cacheKeyNegaraList(page, pageSize int, filter *dto.FilterNegaraRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:negara:list:p%d:ps%d:f%s", page, pageSize, f)
} // Negara-----------------------------------------------------------------------

// Provinsi-----------------------------------------------------------------------
func cacheKeyProvinsiDetail(id int64) string {
	return fmt.Sprintf("master_alamat:provinsi:detail:%d", id)
}
func cacheKeyProvinsiGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:provinsi:getdetail:%d", id)
}
func cacheKeyProvinsiList(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) string {
	nID := "all"
	if negaraID != nil {
		nID = fmt.Sprintf("%d", *negaraID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:provinsi:list:p%d:ps%d:nid%s:f%s", page, pageSize, nID, f)
} // Provinsi---------------------------------------------------------------------

// Kota/Kabupaten-----------------------------------------------------------------
func cacheKeyKotaDetail(id int64) string { return fmt.Sprintf("master_alamat:kota:detail:%d", id) }
func cacheKeyKotaGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kota:getdetail:%d", id)
}
func cacheKeyKotaList(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) string {
	pID := "all"
	if provinsiID != nil {
		pID = fmt.Sprintf("%d", *provinsiID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:kota:list:p%d:ps%d:pid%s:f%s", page, pageSize, pID, f)
} // Kota/Kabupaten---------------------------------------------------------------

// Kecamatan----------------------------------------------------------------------
func cacheKeyKecamatanDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kecamatan:detail:%d", id)
}
func cacheKeyKecamatanGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kecamatan:getdetail:%d", id)
}
func cacheKeyKecamatanList(page, pageSize int, kotaID *int64, filter *dto.FilterKecamatanRequest) string {
	kID := "all"
	if kotaID != nil {
		kID = fmt.Sprintf("%d", *kotaID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:kecamatan:list:p%d:ps%d:kid%s:f%s", page, pageSize, kID, f)
} // Kecamatan--------------------------------------------------------------------

// Kelurahan/Desa-----------------------------------------------------------------
func cacheKeyDesaDetail(id int64) string { return fmt.Sprintf("master_alamat:desa:detail:%d", id) }
func cacheKeyDesaGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:desa:getdetail:%d", id)
}
func cacheKeyDesaList(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) string {
	kecID := "all"
	if kecamatanID != nil {
		kecID = fmt.Sprintf("%d", *kecamatanID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:desa:list:p%d:ps%d:kecid%s:f%s", page, pageSize, kecID, f)
} // Kelurahan/Desa---------------------------------------------------------------

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
	ctx := context.Background()
	cacheKey := cacheKeyNegaraDetail(id)

	// 1. Cek Cache
	var cachedRes dto.NegaraResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Jika tidak ada di cache, ambil dari DB
	m, err := s.repo.GetByIDNegara(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	res := dto.ToNegaraResponse(m)

	// 3. Simpan ke Cache (misal TTL 24 jam)
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListNegara(page, pageSize int, filter *dto.FilterNegaraRequest) ([]dto.NegaraResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	ctx := context.Background()
	cacheKey := cacheKeyNegaraList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.NegaraResponse `json:"items"`
		Total int64                `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Negara tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	// 2. Ambil dari DB
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
		return nil, 0, appErrors.NotFound("Negara tidak ditemukan")
	}

	res := dto.ToNegaraListResponse(items)

	// 3. Simpan ke Cache (misal TTL 1 jam untuk list)
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.NegaraResponse `json:"items"`
		Total int64                `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
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
	res := dto.ToNegaraResponse(m)

	// HAPUS CACHE LIST karena ada data baru
	s.cache.InvalidateList(context.Background(), "master_alamat:negara:list:")

	return res, nil
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
	res := dto.ToNegaraResponse(m)

	// HAPUS CACHE DETAIL & LIST karena data berubah
	s.cache.InvalidateDetail(context.Background(), cacheKeyNegaraDetail(id))
	s.cache.InvalidateList(context.Background(), "master_alamat:negara:list:")

	return res, nil
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
	err = s.repo.DeleteNegara(id)
	if err == nil {
		// HAPUS CACHE DETAIL & LIST karena data dihapus
		s.cache.InvalidateDetail(context.Background(), cacheKeyNegaraDetail(id))
		s.cache.InvalidateList(context.Background(), "master_alamat:negara:list:")
	}
	return err
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Negara =========================================================================

// Provinsi ========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDProvinsi(id int64) (*dto.ProvinsiDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyProvinsiDetail(id)

	var cachedRes dto.ProvinsiDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := &dto.ProvinsiDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		NegaraID:       m.NegaraID,
		NegaraName:     m.Negara.Name,
		TotalKota:      totalKota,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}

	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyProvinsiList(page, pageSize, negaraID, filter)

	var cachedRes struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Provinsi tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := dto.ToProvinsiListResponse(items)
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToProvinsiResponse(m)

	// Invalidate cache provinsi + list child (kota, kecamatan, desa) karena response mereka mengandung data provinsi
	ctx := context.Background()
	s.cache.InvalidateList(ctx, "master_alamat:provinsi:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToProvinsiResponse(m)

	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyProvinsiDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyProvinsiGetDetail(id))
	s.cache.InvalidateList(ctx, "master_alamat:provinsi:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	err = s.repo.DeleteProvinsi(id)
	if err == nil {
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyProvinsiDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyProvinsiGetDetail(id))
		s.cache.InvalidateList(ctx, "master_alamat:provinsi:list:")
		s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
		s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
		s.cache.InvalidateList(ctx, "master_alamat:desa:list:")
	}
	return err
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Provinsi =======================================================================

// Kota/Kabupaten =====================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKotaKabupaten(id int64) (*dto.KotaKabupatenDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyKotaDetail(id)

	var cachedRes dto.KotaKabupatenDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := &dto.KotaKabupatenDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		ProvinsiID:     m.ProvinsiID,
		ProvinsiName:   m.Provinsi.Name,
		NegaraID:       m.Provinsi.NegaraID,
		NegaraName:     m.Provinsi.Negara.Name,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}

	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]dto.KotaKabupatenResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyKotaList(page, pageSize, provinsiID, filter)

	var cachedRes struct {
		Items []dto.KotaKabupatenResponse `json:"items"`
		Total int64                       `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := dto.ToKotaKabupatenListResponse(items)
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KotaKabupatenResponse `json:"items"`
		Total int64                       `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKotaKabupatenResponse(m)

	ctx := context.Background()
	s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKotaKabupatenResponse(m)

	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyKotaDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyKotaGetDetail(id))
	s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	err = s.repo.DeleteKotaKabupaten(id)
	if err == nil {
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyKotaDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyKotaGetDetail(id))
		s.cache.InvalidateList(ctx, "master_alamat:kota:list:")
		s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
		s.cache.InvalidateList(ctx, "master_alamat:desa:list:")
	}
	return err
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kota/Kabupaten =====================================================================

// Kecamatan ===========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKecamatan(id int64) (*dto.KecamatanDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyKecamatanDetail(id)

	var cachedRes dto.KecamatanDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := &dto.KecamatanDetailResponse{
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
	}

	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyKecamatanList(page, pageSize, kotaKabupatenID, filter)

	var cachedRes struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kecamatan tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := dto.ToKecamatanListResponse(items)
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKecamatanResponse(m)

	ctx := context.Background()
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKecamatanResponse(m)

	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyKecamatanDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyKecamatanGetDetail(id))
	s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	err = s.repo.DeleteKecamatan(id)
	if err == nil {
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyKecamatanDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyKecamatanGetDetail(id))
		s.cache.InvalidateList(ctx, "master_alamat:kecamatan:list:")
		s.cache.InvalidateList(ctx, "master_alamat:desa:list:")
	}
	return err
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kecamatan ===========================================================================

// Kelurahan/Desa ========================================================================
// ─────────────── GetByID ────────────────────────────────────────────────────────
func (s *service) GetByIDKelurahanDesa(id int64) (*dto.KelurahanDesaDetailResponse, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyDesaDetail(id)

	var cachedRes dto.KelurahanDesaDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	//db---------------------------------------------------------------------------
	m, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	//out--------------------------------------------------------------------------
	res := &dto.KelurahanDesaDetailResponse{
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
	}

	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── List ───────────────────────────────────────────────────────────
func (s *service) ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error) {
	// ─── Logic ──────────────────────────────────────────────────────────────────
	//cache------------------------------------------------------------------------
	ctx := context.Background()
	cacheKey := cacheKeyDesaList(page, pageSize, kecamatanID, filter)

	var cachedRes struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	//db---------------------------------------------------------------------------
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

	//out--------------------------------------------------------------------------
	res := dto.ToKelurahanDesaListResponse(items)
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKelurahanDesaResponse(m)

	ctx := context.Background()
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	res := dto.ToKelurahanDesaResponse(m)

	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyDesaDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyDesaGetDetail(id))
	s.cache.InvalidateList(ctx, "master_alamat:desa:list:")

	return res, nil
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

	//out--------------------------------------------------------------------------
	err = s.repo.DeleteKelurahanDesa(id)
	if err == nil {
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyDesaDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyDesaGetDetail(id))
		s.cache.InvalidateList(ctx, "master_alamat:desa:list:")
	}
	return err
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kelurahan/Desa ========================================================================
