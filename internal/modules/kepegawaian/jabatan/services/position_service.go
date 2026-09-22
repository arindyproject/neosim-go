package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	departemenModels "neosim_go/internal/modules/master/departemen/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.departemenRepo,
// s.buildCreator, dan s.buildAuditMaps dipakai ulang langsung — tidak perlu
// field/param baru; departemenRepo sudah ada di struct service sejak awal.

// ── Create ────────────────────────────────────────────────────────────────────

func (s *service) CreatePosition(ctx context.Context, req *dto.CreatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error) {
	can, err := s.canCreatePosition(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Position baru.", nil)
	}

	if req.ParentID != nil {
		exists, err := s.repo.ExistsPositionByID(ctx, *req.ParentID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek parent position: " + err.Error())
		}
		if !exists {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Parent position tidak ditemukan.", nil)
		}
	}

	// department bersifat lintas-module (master/departemen), jadi divalidasi
	// lewat s.departemenRepo, bukan s.repo milik jabatan.
	var department *departemenModels.MasterDepartemen
	if req.DepartmentID != nil {
		department, err = s.departemenRepo.GetDepartemenByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek department: " + err.Error())
		}
		if department == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Department tidak ditemukan.", nil)
		}
	}

	dup, err := s.repo.ExistsByNameAndKategori(ctx, req.PositionKategoriID, req.Name, 0)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat nama: " + err.Error())
	}
	if dup {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity,
			"Nama Position sudah dipakai pada kategori yang sama.", nil)
	}

	m := &models.Position{
		Name:               req.Name,
		Description:        req.Description,
		PositionKategoriID: req.PositionKategoriID,
		ParentID:           req.ParentID,
		DepartmentID:       req.DepartmentID,
		LevelHierarki:      req.LevelHierarki,
		Kuota:              req.Kuota,
		IsAktif:            req.IsAktif,
		CreatedBy:          &actor.UserID,
		UpdatedBy:          &actor.UserID,
	}
	if err := s.repo.CreatePosition(ctx, m); err != nil {
		return nil, err
	}

	// reload dengan preload PositionKategori/Parent supaya response lengkap
	saved, err := s.repo.GetPositionByID(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, saved.CreatedBy)

	return dto.ToPositionResponse(dto.PositionResponseParams{
		Position:   saved,
		Creator:    creator,
		Updater:    creator, // saat create, creator dan updater sama
		Department: toDepartemenSimpel(department),
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *service) GetPositionByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.PositionResponse, error) {
	can, err := s.canReadPosition(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Position.", nil)
	}

	m, err := s.repo.GetPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Position tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)
	department := s.buildDepartemenSimpel(ctx, m.DepartmentID)

	return dto.ToPositionResponse(dto.PositionResponseParams{
		Position:   m,
		Creator:    creator,
		Updater:    updater,
		Department: department,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *service) ListPosition(ctx context.Context, page, pageSize int, filter *dto.FilterPositionRequest, actor he.AuthContext) ([]dto.PositionResponse, int64, error) {
	can, err := s.canReadPosition(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Position.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListPosition(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForPosition(ctx, items)
	departmentsMap := s.buildDepartemenMapForPosition(ctx, items)
	return dto.ToPositionListResponse(items, creatorsMap, updatersMap, departmentsMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *service) UpdatePosition(ctx context.Context, id int64, req *dto.UpdatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error) {
	can, err := s.canUpdatePosition(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Position.", nil)
	}

	m, err := s.repo.GetPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Position tidak ditemukan")
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity,
				"Position tidak boleh menjadi parent dari dirinya sendiri.", nil)
		}
		exists, err := s.repo.ExistsPositionByID(ctx, *req.ParentID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek parent position: " + err.Error())
		}
		if !exists {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Parent position tidak ditemukan.", nil)
		}
		m.ParentID = req.ParentID
	}

	var department *departemenModels.MasterDepartemen
	if req.DepartmentID != nil {
		department, err = s.departemenRepo.GetDepartemenByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, appErrors.Internal("gagal cek department: " + err.Error())
		}
		if department == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Department tidak ditemukan.", nil)
		}
		m.DepartmentID = req.DepartmentID
	}

	// tentukan nama & kategori akhir untuk pengecekan duplikat
	finalName := m.Name
	if req.Name != nil {
		finalName = *req.Name
	}
	finalKategoriID := m.PositionKategoriID
	if req.PositionKategoriID != nil {
		finalKategoriID = *req.PositionKategoriID
	}
	if req.Name != nil || req.PositionKategoriID != nil {
		dup, err := s.repo.ExistsByNameAndKategori(ctx, finalKategoriID, finalName, id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat nama: " + err.Error())
		}
		if dup {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity,
				"Nama Position sudah dipakai pada kategori yang sama.", nil)
		}
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.PositionKategoriID != nil {
		m.PositionKategoriID = *req.PositionKategoriID
	}
	if req.LevelHierarki != nil {
		m.LevelHierarki = *req.LevelHierarki
	}
	if req.Kuota != nil {
		m.Kuota = req.Kuota
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdatePosition(ctx, m); err != nil {
		return nil, err
	}

	// reload supaya PositionKategori/Parent di response ikut ter-update
	saved, err := s.repo.GetPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, saved.CreatedBy)
	updater := s.buildCreator(ctx, saved.UpdatedBy)

	// kalau department tidak diubah di request ini, department masih nil di
	// sini — ambil ulang dari saved.DepartmentID supaya response tetap lengkap.
	if department == nil {
		return dto.ToPositionResponse(dto.PositionResponseParams{
			Position:   saved,
			Creator:    creator,
			Updater:    updater,
			Department: s.buildDepartemenSimpel(ctx, saved.DepartmentID),
		}), nil
	}

	return dto.ToPositionResponse(dto.PositionResponseParams{
		Position:   saved,
		Creator:    creator,
		Updater:    updater,
		Department: toDepartemenSimpel(department),
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *service) DeletePosition(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeletePosition(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Position.", nil)
	}

	m, err := s.repo.GetPositionByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Position tidak ditemukan")
	}

	hasChildren, err := s.repo.HasChildren(ctx, id)
	if err != nil {
		return appErrors.Internal("gagal cek anak position: " + err.Error())
	}
	if hasChildren {
		return appErrors.Wrap(http.StatusUnprocessableEntity,
			"Position tidak bisa dihapus karena masih memiliki anak dalam struktur organisasi. Pindahkan atau hapus dulu anaknya.", nil)
	}

	return s.repo.DeletePosition(ctx, id, actor.UserID)
}

// ── helper khusus Position (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForPosition(ctx context.Context, items []models.Position) (map[int64]*he.UserData, map[int64]*he.UserData) {
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

// buildDepartemenSimpel mengambil satu department by ID untuk kebutuhan
// single-item response (GetByID, Update saat department tidak diubah).
// Mengikuti gaya buildCreator: error ditelan, nil kalau gagal/tidak ada —
// department gagal dimuat tidak boleh menggagalkan seluruh response Position.
func (s *service) buildDepartemenSimpel(ctx context.Context, departmentID *int64) *dto.DepartemenSimpelResponse {
	if departmentID == nil {
		return nil
	}
	d, err := s.departemenRepo.GetDepartemenByID(ctx, *departmentID)
	if err != nil || d == nil {
		return nil
	}
	return toDepartemenSimpel(d)
}

// buildDepartemenMapForPosition mengumpulkan department unik dari daftar
// Position lalu di-fetch satu-satu lewat s.departemenRepo.GetByID.
//
// CATATAN: kalau masterDepartemenContracts.Repository sudah/nanti punya
// method batch (mis. GetByIDs), ganti loop ini supaya jadi 1 query alih-alih
// N query — sama seperti buildAuditMapsForPosition di atas untuk users.
func (s *service) buildDepartemenMapForPosition(ctx context.Context, items []models.Position) map[int64]*dto.DepartemenSimpelResponse {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.DepartmentID != nil {
			idSet[*item.DepartmentID] = struct{}{}
		}
	}

	result := make(map[int64]*dto.DepartemenSimpelResponse, len(idSet))
	for id := range idSet {
		d, err := s.departemenRepo.GetDepartemenByID(ctx, id)
		if err != nil || d == nil {
			continue
		}
		result[id] = toDepartemenSimpel(d)
	}
	return result
}

// toDepartemenSimpel mapper kecil model Departemen -> DepartemenSimpelResponse.
//
// ASUMSI: field nama pada model Departemen adalah "Nama" (konvensi Indonesia
// yang dipakai modul master lain di project ini). Kalau modelmu pakai "Name",
// tinggal ganti d.Nama jadi d.Name di sini.
func toDepartemenSimpel(d *departemenModels.MasterDepartemen) *dto.DepartemenSimpelResponse {
	if d == nil {
		return nil
	}
	return &dto.DepartemenSimpelResponse{
		ID:   d.ID,
		Name: d.Name,
	}
}

// ── Tree ──────────────────────────────────────────────────────────────────────

// GetPositionTree menampilkan seluruh bagan organisasi Position sebagai tree
// bersarang, dirakit dari data yang sudah ada di DB (bukan hardcode) lewat
// satu query flat (FindAllPositions) + assembly O(n) di dto.ToPositionTree.
func (s *service) GetPositionTree(ctx context.Context, onlyAktif bool, actor he.AuthContext) ([]dto.PositionTreeNode, error) {
	can, err := s.canReadPosition(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat struktur organisasi Position.", nil)
	}

	items, err := s.repo.FindAllPositions(ctx, onlyAktif)
	if err != nil {
		return nil, err
	}

	departmentsMap := s.buildDepartemenMapForPosition(ctx, items)
	return dto.ToPositionTree(items, departmentsMap), nil
}
