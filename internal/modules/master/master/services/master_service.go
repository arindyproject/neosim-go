package services

import (
	"errors"
	"time"
	"net/http"

	masterContracts "neosim_go/internal/modules/master/master/contracts"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo masterContracts.Repository
	rbacRepo rbacContracts.RBACRepository	//RBAC
	authRepo authContracts.AuthRepository	//AUTH
}

// NewMasterService membuat instance service baru
func NewMasterService(
	repo masterContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,	//RBAC
	authRepo authContracts.AuthRepository,	//AUTH
) masterContracts.Service {
	return &service{
		repo: repo,
		rbacRepo: rbacRepo,	//RBAC
		authRepo: authRepo,	//AUTH
	}
}


// ─── Permission ─────────────────────────────────────────────────────────────────
func (s *service) canReadMaster(actor masterContracts.AuthContext) (bool, error) {
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

func (s *service) canCreateMaster(actor masterContracts.AuthContext) (bool, error) {
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

func (s *service) canUpdateMaster(actor masterContracts.AuthContext) (bool, error) {
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

func (s *service) canDeleteMaster(actor masterContracts.AuthContext) (bool, error) {
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
func (s *service) Create(req *dto.CreateMasterRequest, createdBy *int64, actor masterContracts.AuthContext) (*dto.MasterResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Master baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m := &models.Master{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToMasterResponse(m), nil
}

// ------------ GetByID -----------------------------------------------------------
func (s *service) GetByID(id int64, actor masterContracts.AuthContext) (*dto.MasterResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat Master.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Master tidak ditemukan")
	}
	return dto.ToMasterResponse(m), nil
}

// ------------ List --------------------------------------------------------------
func (s *service) List(page, pageSize int,filter *dto.FilterMasterRequest, actor masterContracts.AuthContext) ([]dto.MasterResponse, int64, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadMaster(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Master.", nil)
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
	return dto.ToMasterListResponse(items), total, nil
}

// ------------ Update ------------------------------------------------------------
func (s *service) Update(id int64, req *dto.UpdateMasterRequest, updatedBy *int64, actor masterContracts.AuthContext) (*dto.MasterResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Master.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Master tidak ditemukan")
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
	return dto.ToMasterResponse(m), nil
}

// ------------ Delete ------------------------------------------------------------
func (s *service) Delete(id int64, actor masterContracts.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Master.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Master tidak ditemukan")
	}
	return s.repo.Delete(id)
}
