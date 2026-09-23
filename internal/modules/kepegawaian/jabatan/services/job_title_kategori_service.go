package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateJobTitleKategori(ctx context.Context, req *dto.CreateJobTitleKategoriRequest, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error) {
	can, err := s.canCreateJobTitleKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat JobTitleKategori baru.", nil)
	}

	//ceck duplicate code
	data, err := s.repo.GetJobTitleKategoriByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan kode ini sudah ada", nil)
	}

	//ceck duplicate label
	data, err = s.repo.GetJobTitleKategoriByLabel(ctx, req.Label)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan label ini sudah ada", nil)
	}

	m := &models.JobTitleKategori{
		Code:      req.Code,
		Label:     req.Label,
		FHIRCode:  req.FHIRCode,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}
	if err := s.repo.CreateJobTitleKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixJobTitleKategoriSelectList)

	return dto.ToJobTitleKategoriResponse(dto.JobTitleKategoriResponseParams{
		JobTitleKategori: m,
		Creator:          creator,
		Updater:          creator, // saat create, creator dan updater sama
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetJobTitleKategoriByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error) {
	can, err := s.canReadJobTitleKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat JobTitleKategori.", nil)
	}

	m, err := s.repo.GetJobTitleKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitleKategori tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToJobTitleKategoriResponse(dto.JobTitleKategoriResponseParams{
		JobTitleKategori: m,
		Creator:          creator,
		Updater:          updater,
	}), nil
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (s *service) ListSelectJobTitleKategori(ctx context.Context, search string, actor he.AuthContext) ([]dto.JobTitleKategoriSimpelResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyJobTitleKategoriSelectList(search)

	can, err := s.canReadJobTitleKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar JobTitleKategori Pegawai.", nil)
	}

	// 1. Cek Cache
	var cachedRes []dto.JobTitleKategoriSimpelResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return cachedRes, nil
	}

	items, err := s.repo.ListSelectTipe(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, appErrors.Wrap(http.StatusNotFound, "data JobTitleKategori tidak ditemukan", nil)
	}

	res := dto.ToJobTitleKategoriSimpelResponse(items)

	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListJobTitleKategori(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleKategoriRequest, actor he.AuthContext) ([]dto.JobTitleKategoriResponse, int64, error) {
	can, err := s.canReadJobTitleKategori(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar JobTitleKategori.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJobTitleKategori(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForJobTitleKategori(ctx, items)
	return dto.ToJobTitleKategoriListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateJobTitleKategori(ctx context.Context, id int64, req *dto.UpdateJobTitleKategoriRequest, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error) {
	can, err := s.canUpdateJobTitleKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah JobTitleKategori.", nil)
	}

	m, err := s.repo.GetJobTitleKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitleKategori tidak ditemukan")
	}

	//ceck duplicate code
	if req.Code != nil && *req.Code != m.Code {
		data, err := s.repo.GetJobTitleKategoriByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "JobTitleKategori dengan kode ini sudah digunakan", nil)
		}
	}

	//ceck duplicate label
	if req.Label != nil && *req.Label != m.Label {
		data, err := s.repo.GetJobTitleKategoriByLabel(ctx, *req.Label)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "JobTitleKategori dengan label ini sudah digunakan", nil)
		}
	}

	// update fields yang diubah
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.FHIRCode != nil {
		m.FHIRCode = req.FHIRCode
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJobTitleKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	res := dto.ToJobTitleKategoriResponse(dto.JobTitleKategoriResponseParams{
		JobTitleKategori: m,
		Creator:          creator,
		Updater:          updater,
	})
	ctxs := context.Background()
	s.cache.InvalidateList(ctxs, cachePrefixJobTitleKategoriSelectList)
	return res, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteJobTitleKategori(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteJobTitleKategori(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus JobTitleKategori.", nil)
	}

	m, err := s.repo.GetJobTitleKategoriByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("JobTitleKategori tidak ditemukan")
	}
	err = s.repo.DeleteJobTitleKategori(ctx, id, actor.UserID)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateList(ctxs, cachePrefixJobTitleKategoriSelectList)
	}
	return err
}

// ── helper khusus JobTitleKategori (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForJobTitleKategori(ctx context.Context, items []models.JobTitleKategori) (map[int64]*he.UserData, map[int64]*he.UserData) {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			idSet[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			idSet[*item.UpdatedBy] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, ids) // ← 1 query total, bukan 40
	if err != nil {
		return map[int64]*he.UserData{}, map[int64]*he.UserData{}
	}

	userMap := make(map[int64]*he.UserData, len(users))
	for _, u := range users {
		userMap[u.ID] = &he.UserData{ID: u.ID, Username: u.Username, Name: u.Name}
	}
	// creator dan updater sekarang share map yang sama — reuse otomatis, kode lebih pendek juga
	return userMap, userMap
}
