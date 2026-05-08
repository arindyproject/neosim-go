package services

import (
	"errors"
	"time"

	"neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"

	"golang.org/x/crypto/bcrypt"
)

// service implements the contracts.Service interface
type service struct {
	repo contracts.Repository
}

// NewService creates a new service instance
func NewService(repo contracts.Repository) contracts.Service {
	return &service{repo: repo}
}

// ─── CRUD ──────────────────────────────────────────────────────────────────────

func (s *service) CreateUser(req *dto.CreateUserRequest, createdBy *int64) (*dto.UserResponse, error) {
	// Cek username
	existing, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	// Cek email
	existing, err = s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah digunakan")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:           req.Username,
		Email:              req.Email,
		Name:               req.Name,
		Password:           hashedPassword,
		IsActive:           true,
		IsVerified:         false,
		MustChangePassword: false,
		CreatedBy:          createdBy,
		UpdatedBy:          createdBy,
	}

	// Set default settings dari dto
	if err := user.SetSettings(dto.DefaultUserSettings()); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) GetUserByID(id int64) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return dto.ToUserResponse(user), nil
}

func (s *service) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return dto.ToUserResponse(user), nil
}

func (s *service) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return dto.ToUserResponse(user), nil
}

func (s *service) ListUsers(page, pageSize int) ([]dto.UserResponse, int64, error) {
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

func (s *service) UpdateUser(id int64, req *dto.UpdateUserRequest, updatedBy *int64) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		existing, err := s.repo.GetByEmail(*req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("email sudah digunakan")
		}
		user.Email = *req.Email
	}
	if req.Photo != nil {
		user.Photo = req.Photo
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsStaff != nil {
		user.IsStaff = *req.IsStaff
	}
	if req.IsSuperuser != nil {
		user.IsSuperuser = *req.IsSuperuser
	}

	user.UpdatedBy = updatedBy
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) DeleteUser(id int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}
	return s.repo.Delete(id)
}

// ─── Password ──────────────────────────────────────────────────────────────────

func (s *service) ChangePassword(id int64, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}

	if !s.verifyPassword(req.OldPassword, user.Password) {
		return errors.New("password lama tidak sesuai")
	}

	hashed, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	user.Password = hashed
	user.PasswordChangedAt = &now
	user.MustChangePassword = false

	return s.repo.Update(user)
}

func (s *service) UpdateLastLogin(id int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}

	now := time.Now()
	user.LastLoginAt = &now

	return s.repo.Update(user)
}

// ─── Settings ──────────────────────────────────────────────────────────────────

func (s *service) GetSettings(id int64) ([]models.UserSetting, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return s.repo.GetSettings(id)
}

func (s *service) UpdateSettings(id int64, req *dto.UpdateSettingsRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}
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
