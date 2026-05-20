package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	authContracts "neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	userContracts "neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"

	appErrors "neosim_go/internal/shared/errors"

	"golang.org/x/crypto/bcrypt"
)

// ─── Auth Context ──────────────────────────────────────────────────────────────

// AuthContext berisi informasi user yang sedang login
// Di-inject dari handler → service untuk authorization
//type AuthContext struct {
//	UserID      int64
//	IsSuperadmin bool
//}

// ─── Service ───────────────────────────────────────────────────────────────────

type service struct {
	repo     userContracts.Repository
	rbacRepo contracts.RBACRepository // ← inject RBAC repo
	authRepo authContracts.AuthRepository
}

func NewUserService(repo userContracts.Repository, rbacRepo contracts.RBACRepository, authRepo authContracts.AuthRepository) userContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}

// ─── Authorization Helper ──────────────────────────────────────────────────────

// canUpdateUser mengecek apakah actor boleh update targetUserID
//
// Boleh update jika:
// 1. Superadmin
// 2. Dirinya sendiri (actor == target)
// 3. Memiliki permission "users:update"
// 4. Memiliki role "hrd"

func (s *service) canCreateuser(actor userContracts.AuthContext) (bool, error) {
	// 1. Superadmin — boleh semua
	if actor.IsSuperadmin {
		return true, nil
	}

	// 2. Cek permission users:update
	hasPermission, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersCreate)
	if err != nil {
		return false, err
	}
	if hasPermission {
		return true, nil
	}

	// 3. Cek any Roles
	hasRole, err := rbacMiddlewares.HasAnyRole(s.rbacRepo, actor.UserID, "admin", "superadmin", "hrd")
	if err != nil {
		return false, err
	}
	if hasRole {
		return true, nil
	}

	return false, nil

}

func (s *service) canUpdateUser(actor userContracts.AuthContext, targetUserID int64) (bool, error) {
	// 1. Superadmin — boleh semua
	if actor.IsSuperadmin {
		return true, nil
	}

	// 2. Dirinya sendiri
	if actor.UserID == targetUserID {
		return true, nil
	}

	// 3. Cek permission users:update
	hasPermission, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersUpdate)
	if err != nil {
		return false, err
	}
	if hasPermission {
		return true, nil
	}

	// 4. Cek role hrd
	hasRole, err := rbacMiddlewares.HasRole(s.rbacRepo, actor.UserID, "hrd")
	if err != nil {
		return false, err
	}
	if hasRole {
		return true, nil
	}

	return false, nil
}

// canDeleteUser mengecek apakah actor boleh delete user
//
// Boleh delete jika:
// 1. Superadmin
// 2. Memiliki permission "users:delete"
func (s *service) canDeleteUser(actor userContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	return rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersDelete)
}

// canReadUser mengecek apakah actor boleh membaca data user
//
// Boleh read jika:
// 1. Superadmin
// 2. Dirinya sendiri
// 3. Memiliki permission "users:read"
func (s *service) canReadUser(actor userContracts.AuthContext, targetUserID int64) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if actor.UserID == targetUserID {
		return true, nil
	}
	return rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersRead)
}

// ─── CRUD ──────────────────────────────────────────────────────────────────────

func (s *service) CreateUser(req *dto.CreateUserRequest, actor userContracts.AuthContext) (*dto.UserSimpleResponse, error) {
	//----------------------------------------------------------------------------
	can, err := s.canCreateuser(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}

	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki Hak Akses untuk mebuat User Baru.", nil)
	}
	//----------------------------------------------------------------------------

	existing, _ := s.repo.GetByUsername(req.Username)
	if existing != nil {
		return nil, appErrors.BadRequest("username sudah digunakan")
	}
	existing, _ = s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, appErrors.BadRequest("email sudah digunakan")
	}

	hashed, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, appErrors.Internal("gagal memproses password")
	}

	defaultSettingsList := models.DefaultSettings() // Pastikan huruf 'D' besar jika diakses dari luar package models
	settingsBytes, err := json.Marshal(defaultSettingsList)
	if err != nil {
		return nil, appErrors.Internal("gagal memproses setting bawaan")
	}

	isActive := true // Default IsActive = true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isStaff := false // Default IsStaff = false
	if req.IsStaff != nil {
		isStaff = *req.IsStaff
	}

	isSuperuser := false // Default IsSuperuser = false
	if req.IsSuperadmin != nil {
		isSuperuser = *req.IsSuperadmin
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Name:         req.Name,
		Password:     hashed,
		IsActive:     isActive,
		IsStaff:      isStaff,
		IsSuperadmin: isSuperuser,
		Settings:     models.JSONB(settingsBytes),
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, appErrors.Internal("gagal membuat user")
	}

	return dto.ToUserSimpleResponse(user), nil
}

func (s *service) GetUserByID(id int64, actor userContracts.AuthContext) (*dto.UserResponse, error) {
	// Authorization
	can, err := s.canReadUser(actor, id)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Anda tidak memiliki izin untuk melihat data ini.", nil)
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	var creatorDTO *models.UserCreator
	creatorUser, err := s.repo.GetByID(*user.CreatedBy)
	if err == nil {
		creatorDTO = &models.UserCreator{
			ID:       creatorUser.ID,
			Username: creatorUser.Username,
			Name:     creatorUser.Name,
		}
	}

	histories, _ := s.authRepo.GetUserLoginHistories(id, 10)
	return dto.ToUserResponse(user, histories, creatorDTO), nil
}

func (s *service) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	var creatorDTO *models.UserCreator
	creatorUser, err := s.repo.GetByID(*user.CreatedBy)
	if err == nil {
		creatorDTO = &models.UserCreator{
			ID:       creatorUser.ID,
			Username: creatorUser.Username,
			Name:     creatorUser.Name,
		}
	}

	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)
	return dto.ToUserResponse(user, histories, creatorDTO), nil
}

func (s *service) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil || user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	var creatorDTO *models.UserCreator
	creatorUser, err := s.repo.GetByID(*user.CreatedBy)
	if err == nil {
		creatorDTO = &models.UserCreator{
			ID:       creatorUser.ID,
			Username: creatorUser.Username,
			Name:     creatorUser.Name,
		}
	}
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)
	return dto.ToUserResponse(user, histories, creatorDTO), nil
}

func (s *service) ListUsers(page, pageSize int) ([]dto.UserSimpleResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	users, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToUserListResponse(users), total, nil
}

// UpdateUser — authorization: superadmin ATAU diri sendiri ATAU punya permission users:update ATAU role hrd
func (s *service) UpdateUser(id int64, req *dto.UpdateUserRequest, actor userContracts.AuthContext) (*dto.UserResponse, error) {
	// Authorization
	can, err := s.canUpdateUser(actor, id)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Hanya superadmin, diri sendiri, atau yang memiliki permission 'users:update' / role 'hrd' yang bisa mengubah data ini.", nil)
	}

	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	// Update fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		existing, _ := s.repo.GetByEmail(*req.Email)
		if existing != nil && existing.ID != id {
			return nil, appErrors.BadRequest("email sudah digunakan")
		}
		user.Email = *req.Email
	}
	if req.Photo != nil {
		user.Photo = req.Photo
	}

	// Field sensitif hanya bisa diubah oleh superadmin / yang punya permission users:manage
	if req.IsActive != nil || req.IsStaff != nil || req.IsSuperadmin != nil {
		canManage := actor.IsSuperadmin
		if !canManage {
			canManage, _ = rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersManage)
		}
		if !canManage {
			return nil, appErrors.Wrap(http.StatusForbidden,
				"Akses ditolak. Hanya superadmin atau yang memiliki permission 'users:manage' yang bisa mengubah status user.", nil)
		}

		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}
		if req.IsStaff != nil {
			user.IsStaff = *req.IsStaff
		}
		if req.IsSuperadmin != nil {
			user.IsSuperadmin = *req.IsSuperadmin
		}
	}

	updatedBy := actor.UserID
	user.UpdatedBy = &updatedBy
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, appErrors.Internal("gagal mengupdate user")
	}

	var creatorDTO *models.UserCreator
	creatorUser, err := s.repo.GetByID(*user.CreatedBy)
	if err == nil {
		creatorDTO = &models.UserCreator{
			ID:       creatorUser.ID,
			Username: creatorUser.Username,
			Name:     creatorUser.Name,
		}
	}
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)
	return dto.ToUserResponse(user, histories, creatorDTO), nil
}

// DeleteUser — hanya superadmin atau yang punya permission users:delete
func (s *service) DeleteUser(id int64, actor userContracts.AuthContext) error {
	can, err := s.canDeleteUser(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Hanya superadmin atau yang memiliki permission 'users:delete' yang bisa menghapus user.", nil)
	}

	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return appErrors.NotFound("user tidak ditemukan")
	}

	// Tidak bisa hapus diri sendiri
	if user.ID == actor.UserID {
		return appErrors.BadRequest("tidak bisa menghapus akun sendiri")
	}

	return s.repo.Delete(id)
}

// ─── Password ──────────────────────────────────────────────────────────────────

func (s *service) ChangePassword(
	id int64,
	req *dto.ChangePasswordRequest,
	actor userContracts.AuthContext, // ✅ ini penting
) error {
	// Hanya diri sendiri atau superadmin
	if !actor.IsSuperadmin && actor.UserID != id {
		return appErrors.Wrap(http.StatusForbidden, "Akses ditolak. Hanya bisa mengubah password sendiri.", nil)
	}

	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return appErrors.NotFound("user tidak ditemukan")
	}

	if !s.verifyPassword(req.OldPassword, user.Password) {
		return appErrors.BadRequest("password lama tidak sesuai")
	}

	hashed, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return appErrors.Internal("gagal memproses password")
	}

	now := time.Now()
	user.Password = hashed
	user.PasswordChangedAt = &now

	return s.repo.Update(user)
}

func (s *service) UpdateLastLogin(id int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}
	now := time.Now()
	user.LastLoginAt = &now
	return s.repo.Update(user)
}

// ─── Settings ──────────────────────────────────────────────────────────────────

func (s *service) GetSettings(id int64) ([]models.UserSetting, error) {
	return s.repo.GetSettings(id)
}

func (s *service) UpdateSettings(id int64, req *dto.UpdateSettingsRequest) error {
	return s.repo.UpdateSettings(id, req.Settings)
}

// ─── Private Helpers ───────────────────────────────────────────────────────────

func (s *service) verifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
