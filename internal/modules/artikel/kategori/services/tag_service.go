package services

import (
	"neosim_go/config"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/artikel/kategori/contracts"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type tagService struct {
	repo     contracts.TagRepository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	cfg      *config.Config
}

// NewTagService membuat instance service baru
func NewTagService(
	repo contracts.TagRepository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	cfg *config.Config,
) contracts.TagService {
	return &tagService{repo: repo, rbacRepo: rbacRepo, authRepo: authRepo, cfg: cfg}
}

func (s *tagService) Create(req *dto.CreateTagRequest, createdBy *int64, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canCreateTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk membuat Tag baru.", nil)
	}
	m := &models.Tag{Name: req.Name, Description: req.Description, CreatedBy: createdBy, UpdatedBy: createdBy}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToTagResponse(m), nil
}

func (s *tagService) GetByID(id int64, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canReadTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tag.", nil)
	}
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}
	return dto.ToTagResponse(m), nil
}

func (s *tagService) List(page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error) {
	can, err := s.canReadTag(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tag.", nil)
	}
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax { pageSize = s.cfg.DefaultPageSize }
	items, total, err := s.repo.List(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToTagListResponse(items), total, nil
}

func (s *tagService) Update(id int64, req *dto.UpdateTagRequest, updatedBy *int64, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canUpdateTag(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Tag.", nil)
	}
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}
	if req.Name != nil { m.Name = *req.Name }
	if req.Description != nil { m.Description = req.Description }
	m.UpdatedBy = updatedBy
	m.UpdatedAt = time.Now()
	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return dto.ToTagResponse(m), nil
}

func (s *tagService) Delete(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteTag(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Tag.", nil)
	}
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Tag tidak ditemukan")
	}
	return s.repo.Delete(id)
}
