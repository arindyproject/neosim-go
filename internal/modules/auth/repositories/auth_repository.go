package repositories

import (
	"context"
	"time"

	"neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/auth/models"

	"gorm.io/gorm"
)

// ─── Init ──────────────────────────────────────────────────────────────────────
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository membuat instance repository baru
func NewAuthRepository(db *gorm.DB) contracts.AuthRepository {
	return &authRepository{db: db}
}

// ─── End Init ──────────────────────────────────────────────────────────────────

// ─── Auth Token ────────────────────────────────────────────────────────────────

func (r *authRepository) SaveToken(ctx context.Context, token *models.AuthToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authRepository) GetTokenByJTI(ctx context.Context, jti string) (*models.AuthToken, error) {
	var token models.AuthToken
	result := r.db.WithContext(ctx).Where("jti = ? AND deleted_at IS NULL", jti).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &token, nil
}

func (r *authRepository) BlacklistToken(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).Model(&models.AuthToken{}).
		Where("jti = ?", jti).
		Update("is_blacklist", true).Error
}

func (r *authRepository) BlacklistAllUserTokens(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&models.AuthToken{}).
		Where("user_id = ? AND is_blacklist = false", userID).
		Update("is_blacklist", true).Error
}

func (r *authRepository) CountActiveTokens(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AuthToken{}).
		Where("user_id = ? AND is_blacklist = false AND expires_at > ? AND deleted_at IS NULL",
			userID, time.Now()).
		Count(&count).Error
	return count, err
}

// ─── Login History ─────────────────────────────────────────────────────────────

func (r *authRepository) SaveLoginHistory(ctx context.Context, history *models.LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *authRepository) GetUserLoginHistories(ctx context.Context, userID int64, limit int) ([]models.LoginHistory, error) {
	// Inisialisasi slice kosong (bukan nil) agar aman saat di-return
	histories := make([]models.LoginHistory, 0)

	// Jika limit yang diinput kurang dari atau sama dengan 0, beri default 10
	if limit <= 0 {
		limit = 10
	}

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}

// ─── Password History ──────────────────────────────────────────────────────────

func (r *authRepository) SavePasswordHistory(ctx context.Context, history *models.PasswordHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *authRepository) GetPasswordHistories(ctx context.Context, userID int64, limit int) ([]models.PasswordHistory, error) {
	var histories []models.PasswordHistory
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
