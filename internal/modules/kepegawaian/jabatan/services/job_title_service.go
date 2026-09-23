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
//
// ASUMSI: s.repo (contracts.Repository) sudah/akan meng-embed
// JobTitleKategoriRepository & JobTitleRumpunProfesiRepository juga (sama
// seperti PositionRepository), jadi punya method ExistsJobTitleKategoriByID
// dan ExistsJobTitleRumpunProfesiByID — belum pernah kamu kirim kode dua
// repo lookup itu, jadi kalau nama methodnya beda tinggal disesuaikan.

// ── ListSelect ────────────────────────────────────────────────────────────────
func (s *service) ListSelectJobTitle(ctx context.Context, search string, actor he.AuthContext) ([]dto.JobTitleSimpelResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyJobTitlesSelectList(search)

	can, err := s.canReadJobTitle(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tipe Kontak Pegawai.", nil)
	}

	// 1. Cek Cache
	var cachedRes []dto.JobTitleSimpelResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return cachedRes, nil
	}

	items, err := s.repo.ListSelectJobTitle(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, appErrors.Wrap(http.StatusNotFound, "data JobTitle tidak ditemukan", nil)
	}

	res := dto.ToJobTitleSimpelResponse(items)

	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *service) CreateJobTitle(ctx context.Context, req *dto.CreateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canCreateJobTitle(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat JobTitle baru.", nil)
	}

	kategoriExists, err := s.repo.CheckJobTitleKategori(ctx, req.KategoriID)
	if err != nil {
		return nil, appErrors.Internal("gagal cek kategori: " + err.Error())
	}
	if !kategoriExists {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Kategori JobTitle tidak ditemukan.", nil)
	}

	if req.RumpunProfesiID != nil {
		exists, err := s.repo.CheckJobTitleRumpunProfesi(ctx, *req.RumpunProfesiID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek rumpun profesi: " + err.Error())
		}
		if !exists {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Rumpun profesi tidak ditemukan.", nil)
		}
	}

	dup, err := s.repo.ExistsByCode(ctx, req.Code, 0)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code: " + err.Error())
	}
	if dup {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Code JobTitle sudah dipakai.", nil)
	}

	m := &models.JobTitle{
		Code:            req.Code,
		Label:           req.Label,
		Description:     req.Description,
		KategoriID:      req.KategoriID,
		RumpunProfesiID: req.RumpunProfesiID,
		MemerlukanSTR:   req.MemerlukanSTR,
		MemerlukanSIP:   req.MemerlukanSIP,
		JenjangMin:      req.JenjangMin,
		Point:           req.Point,
		FHIRCode:        req.FHIRCode,
		FHIRSystem:      req.FHIRSystem,
		IsAktif:         req.IsAktif,
		CreatedBy:       &actor.UserID,
		UpdatedBy:       &actor.UserID,
	}
	if err := s.repo.CreateJobTitle(ctx, m); err != nil {
		return nil, err
	}

	// reload dengan preload Kategori/RumpunProfesi supaya response lengkap
	saved, err := s.repo.GetJobTitleByID(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, saved.CreatedBy)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixJobTitlesSelectList)

	return dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: saved,
		Creator:  creator,
		Updater:  creator, // saat create, creator dan updater sama
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *service) GetJobTitleByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canReadJobTitle(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitle tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: m,
		Creator:  creator,
		Updater:  updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *service) ListJobTitle(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRequest, actor he.AuthContext) ([]dto.JobTitleResponse, int64, error) {
	can, err := s.canReadJobTitle(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar JobTitle.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJobTitle(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForJobTitle(ctx, items)
	return dto.ToJobTitleListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *service) UpdateJobTitle(ctx context.Context, id int64, req *dto.UpdateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canUpdateJobTitle(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitle tidak ditemukan")
	}

	if req.KategoriID != nil {
		exists, err := s.repo.CheckJobTitleKategori(ctx, *req.KategoriID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek kategori: " + err.Error())
		}
		if !exists {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Kategori JobTitle tidak ditemukan.", nil)
		}
		m.KategoriID = *req.KategoriID
	}

	if req.RumpunProfesiID != nil {
		exists, err := s.repo.CheckJobTitleRumpunProfesi(ctx, *req.RumpunProfesiID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek rumpun profesi: " + err.Error())
		}
		if !exists {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Rumpun profesi tidak ditemukan.", nil)
		}
		m.RumpunProfesiID = req.RumpunProfesiID
	}

	if req.Code != nil {
		dup, err := s.repo.ExistsByCode(ctx, *req.Code, id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code: " + err.Error())
		}
		if dup {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Code JobTitle sudah dipakai.", nil)
		}
		m.Code = *req.Code
	}

	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.MemerlukanSTR != nil {
		m.MemerlukanSTR = *req.MemerlukanSTR
	}
	if req.MemerlukanSIP != nil {
		m.MemerlukanSIP = *req.MemerlukanSIP
	}
	if req.Point != nil {
		m.Point = req.Point
	}
	if req.JenjangMin != nil {
		m.JenjangMin = req.JenjangMin
	}
	if req.FHIRCode != nil {
		m.FHIRCode = req.FHIRCode
	}
	if req.FHIRSystem != nil {
		m.FHIRSystem = req.FHIRSystem
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJobTitle(ctx, m); err != nil {
		return nil, err
	}

	// reload supaya Kategori/RumpunProfesi di response ikut ter-update
	saved, err := s.repo.GetJobTitleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, saved.CreatedBy)
	updater := s.buildCreator(ctx, saved.UpdatedBy)

	res := dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: saved,
		Creator:  creator,
		Updater:  updater,
	})

	ctxs := context.Background()
	s.cache.InvalidateList(ctxs, cachePrefixJobTitlesSelectList)
	return res, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *service) DeleteJobTitle(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteJobTitle(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("JobTitle tidak ditemukan")
	}
	err = s.repo.DeleteJobTitle(ctx, id, actor.UserID)
	if err == nil {
		ctxs := context.Background()
		s.cache.InvalidateList(ctxs, cachePrefixJobTitlesSelectList)
	}
	return err
}

// ── helper khusus JobTitle (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForJobTitle(ctx context.Context, items []models.JobTitle) (map[int64]*he.UserData, map[int64]*he.UserData) {
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
