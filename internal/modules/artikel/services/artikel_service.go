package services

import (
	"errors"
	"time"
	"net/http"

	artikelContracts "neosim_go/internal/modules/artikel/contracts"
	"neosim_go/internal/modules/artikel/dto"
	"neosim_go/internal/modules/artikel/models"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo artikelContracts.Repository
	rbacRepo rbacContracts.RBACRepository	//RBAC
	authRepo authContracts.AuthRepository	//AUTH
}

// NewArtikelService membuat instance service baru
func NewArtikelService(
	repo artikelContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,	//RBAC
	authRepo authContracts.AuthRepository,	//AUTH
) artikelContracts.Service {
	return &service{
		repo: repo,
		rbacRepo: rbacRepo,	//RBAC
		authRepo: authRepo,	//AUTH
	}
}


// ─── Permission ─────────────────────────────────────────────────────────────────
func (s *service) canReadArtikel(actor artikelContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyRead); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canCreateArtikel(actor artikelContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canUpdateArtikel(actor artikelContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canDeleteArtikel(actor artikelContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}

// ─── Service ────────────────────────────────────────────────────────────────────

// ------------ Create ------------------------------------------------------------
func (s *service) Create(req *dto.CreateArtikelRequest, createdBy *int64, actor artikelContracts.AuthContext) (*dto.ArtikelResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Artikel baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m := &models.Artikel{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToArtikelResponse(m), nil
}

// ------------ GetByID -----------------------------------------------------------
func (s *service) GetByID(id int64, actor artikelContracts.AuthContext) (*dto.ArtikelResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat Artikel.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}
	return dto.ToArtikelResponse(m), nil
}

// ------------ List --------------------------------------------------------------
func (s *service) List(page, pageSize int,filter *dto.FilterArtikelRequest, actor artikelContracts.AuthContext) ([]dto.ArtikelResponse, int64, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadArtikel(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Artikel.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.List(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToArtikelListResponse(items), total, nil
}

// ------------ Update ------------------------------------------------------------
func (s *service) Update(id int64, req *dto.UpdateArtikelRequest, updatedBy *int64, actor artikelContracts.AuthContext) (*dto.ArtikelResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Artikel.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("artikel tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = updatedBy
	m.UpdatedAt = time.Now()

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return dto.ToArtikelResponse(m), nil
}

// ------------ Delete ------------------------------------------------------------
func (s *service) Delete(id int64, actor artikelContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteArtikel(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Artikel.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("artikel tidak ditemukan")
	}
	return s.repo.Delete(id)
}
