package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacDto "neosim_go/internal/modules/rbac/dto"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	userContracts "neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"

	appErrors "neosim_go/internal/shared/errors"

	"golang.org/x/crypto/bcrypt"
)

// ─── Service ───────────────────────────────────────────────────────────────────

type service struct {
	repo     userContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

func NewUserService(
	repo userContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) userContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}

// ─── Authorization Helpers ─────────────────────────────────────────────────────

func (s *service) canCreateUser(actor userContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersCreate); err != nil || has {
		return has, err
	}
	return rbacMiddlewares.HasAnyRole(s.rbacRepo, actor.UserID, "admin", "superadmin", "hrd")
}

func (s *service) canUpdateUser(actor userContracts.AuthContext, targetUserID int64) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if actor.UserID == targetUserID {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersUpdate); err != nil || has {
		return has, err
	}
	return rbacMiddlewares.HasRole(s.rbacRepo, actor.UserID, "hrd")
}

func (s *service) canDeleteUser(actor userContracts.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	return rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersDelete)
}

func (s *service) canReadUser(actor userContracts.AuthContext, targetUserID int64) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if actor.UserID == targetUserID {
		return true, nil
	}
	return rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermUsersRead)
}

// ─── RBAC Data Builder ─────────────────────────────────────────────────────────

// buildUserRBAC mengambil roles dan permissions untuk user — dipanggil saat butuh UserResponse
func (s *service) buildUserRBAC(userID int64) ([]rbacDto.RoleResponse, []string) {
	// Ambil roles
	roles, err := s.rbacRepo.GetUserRoles(userID)
	var roleResponses []rbacDto.RoleResponse
	if err == nil {
		roleResponses = rbacDto.ToRoleListResponse(roles)
	}

	// Ambil semua permissions (dari role + direct)
	perms, err := s.rbacRepo.GetUserAllPermissions(userID)
	if err != nil {
		perms = []string{}
	}

	return roleResponses, perms
}

// buildCreator mengambil data creator user
func (s *service) buildCreator(createdBy *int64) *models.UserCreator {
	if createdBy == nil {
		return nil
	}
	creator, err := s.repo.GetByID(*createdBy)
	if err != nil || creator == nil {
		return nil
	}
	return &models.UserCreator{
		ID:       creator.ID,
		Username: creator.Username,
		Name:     creator.Name,
	}
}

// ─── CRUD ──────────────────────────────────────────────────────────────────────

func (s *service) CreateUser(req *dto.CreateUserRequest, actor userContracts.AuthContext) (*dto.UserSimpleResponse, error) {
	can, err := s.canCreateUser(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat user baru.", nil)
	}

	if existing, _ := s.repo.GetByUsername(req.Username); existing != nil {
		return nil, appErrors.BadRequest("username sudah digunakan")
	}
	if existing, _ := s.repo.GetByEmail(req.Email); existing != nil {
		return nil, appErrors.BadRequest("email sudah digunakan")
	}

	hashed, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, appErrors.Internal("gagal memproses password")
	}

	defaultSettingsList := models.DefaultSettings()
	settingsBytes, err := json.Marshal(defaultSettingsList)
	if err != nil {
		return nil, appErrors.Internal("gagal memproses setting bawaan")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	isStaff := false
	if req.IsStaff != nil {
		isStaff = *req.IsStaff
	}
	isSuperadmin := false
	if req.IsSuperadmin != nil {
		isSuperadmin = *req.IsSuperadmin
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Name:         req.Name,
		Password:     hashed,
		IsActive:     isActive,
		IsStaff:      isStaff,
		IsSuperadmin: isSuperadmin,
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

	// Ambil RBAC data
	roles, permissions := s.buildUserRBAC(user.ID)

	// Ambil creator
	creator := s.buildCreator(user.CreatedBy)

	// Ambil login histories
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)

	return dto.ToUserResponse(dto.UserResponseParams{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		Histories:   histories,
		Creator:     creator,
	}), nil
}

func (s *service) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	roles, permissions := s.buildUserRBAC(user.ID)
	creator := s.buildCreator(user.CreatedBy)
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)

	return dto.ToUserResponse(dto.UserResponseParams{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		Histories:   histories,
		Creator:     creator,
	}), nil
}

func (s *service) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil || user == nil {
		return nil, appErrors.NotFound("user tidak ditemukan")
	}

	roles, permissions := s.buildUserRBAC(user.ID)
	creator := s.buildCreator(user.CreatedBy)
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)

	return dto.ToUserResponse(dto.UserResponseParams{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		Histories:   histories,
		Creator:     creator,
	}), nil
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

func (s *service) UpdateUser(id int64, req *dto.UpdateUserRequest, actor userContracts.AuthContext) (*dto.UserResponse, error) {
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

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		if existing, _ := s.repo.GetByEmail(*req.Email); existing != nil && existing.ID != id {
			return nil, appErrors.BadRequest("email sudah digunakan")
		}
		user.Email = *req.Email
	}
	if req.Photo != nil {
		user.Photo = req.Photo
	}

	// Field sensitif — hanya superadmin / users:manage
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

	user.UpdatedBy = &actor.UserID
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, appErrors.Internal("gagal mengupdate user")
	}

	roles, permissions := s.buildUserRBAC(user.ID)
	creator := s.buildCreator(user.CreatedBy)
	histories, _ := s.authRepo.GetUserLoginHistories(user.ID, 10)

	return dto.ToUserResponse(dto.UserResponseParams{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		Histories:   histories,
		Creator:     creator,
	}), nil
}

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
	if user.ID == actor.UserID {
		return appErrors.BadRequest("tidak bisa menghapus akun sendiri")
	}

	return s.repo.Delete(id)
}

// ─── Password ──────────────────────────────────────────────────────────────────

func (s *service) ChangePassword(id int64, req *dto.ChangePasswordRequest, actor userContracts.AuthContext) error {
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
