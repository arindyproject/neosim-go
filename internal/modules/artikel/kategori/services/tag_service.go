package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

func (s *service) CreateTag(req *dto.CreateTagRequest, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canCreateTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Tag baru.", nil)
	}

	m := &models.Tag{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateTag(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

func (s *service) GetTagByID(id int64, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canReadTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tag.", nil)
	}

	m, err := s.repo.GetTagByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

func (s *service) ListTag(page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error) {
	can, err := s.canReadTag(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tag.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListTag(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForTag(items)
	return toTagResponses(items, creatorsMap, updatersMap), total, nil
}

func (s *service) UpdateTag(id int64, req *dto.UpdateTagRequest, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canUpdateTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Tag.", nil)
	}

	m, err := s.repo.GetTagByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateTag(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

func (s *service) DeleteTag(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteTag(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Tag.", nil)
	}

	m, err := s.repo.GetTagByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Tag tidak ditemukan")
	}
	return s.repo.DeleteTag(id)
}

// ── helper khusus Tag (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForTag(items []models.Tag) (map[int64]*he.UserData, map[int64]*he.UserData) {
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

func toTagResponses(
	items []models.Tag,
	creatorsMap, updatersMap map[int64]*he.UserData,
) []dto.TagResponse {
	responses := make([]dto.TagResponse, 0, len(items))
	for _, item := range items {
		var creator, updater *he.UserData
		if item.CreatedBy != nil {
			creator = creatorsMap[*item.CreatedBy]
		}
		if item.UpdatedBy != nil {
			updater = updatersMap[*item.UpdatedBy]
		}
		responses = append(responses, *dto.ToTagResponse(dto.TagResponseParams{
			Tag: &item,
			Creator:       creator,
			Updater:       updater,
		}))
	}
	return responses
}
