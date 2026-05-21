package repositories

import (
	"errors"

	"neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/models"

	"gorm.io/gorm"
)

// ─── Init ──────────────────────────────────────────────────────────────────────
// repository implements the contracts.Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}

// ─── End Init ──────────────────────────────────────────────────────────────────

// Create creates a new user
func (r *repository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(id int64) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil

}

// GetByUsername retrieves a user by username
func (r *repository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// List retrieves paginated list of users
func (r *repository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates an existing user
func (r *repository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a user
func (r *repository) Delete(id int64) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetSettings retrieves user settings
func (r *repository) GetSettings(id int64) ([]models.UserSetting, error) {
	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return user.GetSettings()
}

// UpdateSettings updates user settings
func (r *repository) UpdateSettings(id int64, settings []models.UserSetting) error {
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}

	if err := user.SetSettings(settings); err != nil {
		return err
	}

	return r.Update(user)
}
