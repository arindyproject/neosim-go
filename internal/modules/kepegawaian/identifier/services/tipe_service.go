package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

func (s *service) CreateTipe(req *dto.CreateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canCreateTipe(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Tipe baru.", nil)
	}

	// Check Duplicate code
	data, err := s.repo.GetTipeByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan kode ini sudah ada", nil)
	}

	// Check Duplicate label
	data, err = s.repo.GetTipeByLabel(req.Label)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan label ini sudah ada", nil)
	}

	// Buat instance model Tipe baru dengan bidang-bidang yang disesuaikan
	m := &models.Tipe{
		Code:        req.Code,
		Label:       req.Label,
		Penerbit:    req.Penerbit,
		FHIRSystem:  req.FHIRSystem,
		HasExpiry:   req.HasExpiry,
		IsNakes:     req.IsNakes,
		IsRequired:  req.IsRequired,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}

	if err := s.repo.CreateTipe(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

func (s *service) GetTipeByID(id int64, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

func (s *service) GetTipeByCode(code string, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByCode(code)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

func (s *service) GetTipeByLabel(label string, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByLabel(label)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

func (s *service) ListTipe(page, pageSize int, filter *dto.FilterTipeRequest, actor he.AuthContext) ([]dto.TipeResponse, int64, error) {
	can, err := s.canReadTipe(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tipe.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.ListTipe(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForTipe(items)
	return toTipeResponses(items, creatorsMap, updatersMap), total, nil
}

func (s *service) UpdateTipe(id int64, req *dto.UpdateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canUpdateTipe(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	// Check Duplicate code jika ada perubahan
	if req.Code != nil && *req.Code != m.Code {
		data, err := s.repo.GetTipeByCode(*req.Code)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan kode ini sudah digunakan", nil)
		}
	}

	// Check Duplicate label jika ada perubahan
	if req.Label != nil && *req.Label != m.Label {
		data, err := s.repo.GetTipeByLabel(*req.Label)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan label ini sudah digunakan", nil)
		}
	}

	// Update parsial sesuai pointer/field pada DTO Update
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.Penerbit != nil {
		m.Penerbit = req.Penerbit
	}
	if req.FHIRSystem != nil {
		m.FHIRSystem = req.FHIRSystem
	}
	if req.HasExpiry != nil {
		m.HasExpiry = *req.HasExpiry
	}
	if req.IsNakes != nil {
		m.IsNakes = *req.IsNakes
	}
	if req.IsRequired != nil {
		m.IsRequired = *req.IsRequired
	}
	if req.Description != nil {
		m.Description = req.Description
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateTipe(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

func (s *service) DeleteTipe(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteTipe(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Tipe tidak ditemukan")
	}
	return s.repo.DeleteTipe(id)
}

// ── helper khusus Tipe (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForTipe(items []models.Tipe) (map[int64]*he.UserData, map[int64]*he.UserData) {
	fetchUser := func(id int64) (*he.UserData, error) {
		user, err := s.userRepo.GetByID(id)
		if err != nil || user == nil {
			return nil, err
		}
		return &he.UserData{ID: user.ID, Username: user.Username, Name: user.Name}, nil
	}

	creatorIDs := make(map[int64]struct{})
	updaterIDs := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			creatorIDs[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			updaterIDs[*item.UpdatedBy] = struct{}{}
		}
	}

	creatorsMap := make(map[int64]*he.UserData)
	for id := range creatorIDs {
		if data, err := fetchUser(id); err == nil && data != nil {
			creatorsMap[id] = data
		}
	}

	updatersMap := make(map[int64]*he.UserData)
	for id := range updaterIDs {
		if data, ok := creatorsMap[id]; ok {
			updatersMap[id] = data
		} else if data, err := fetchUser(id); err == nil && data != nil {
			updatersMap[id] = data
		}
	}

	return creatorsMap, updatersMap
}

func toTipeResponses(
	items []models.Tipe,
	creatorsMap, updatersMap map[int64]*he.UserData,
) []dto.TipeResponse {
	responses := make([]dto.TipeResponse, 0, len(items))
	for _, item := range items {
		var creator, updater *he.UserData
		if item.CreatedBy != nil {
			creator = creatorsMap[*item.CreatedBy]
		}
		if item.UpdatedBy != nil {
			updater = updatersMap[*item.UpdatedBy]
		}
		responses = append(responses, *dto.ToTipeResponse(dto.TipeResponseParams{
			Tipe:    &item,
			Creator: creator,
			Updater: updater,
		}))
	}
	return responses
}
