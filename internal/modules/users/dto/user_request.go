package dto

import (
	"fmt"
	"neosim_go/config"
	"neosim_go/internal/modules/users/models"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ─── Request DTOs ──────────────────────────────────────────────────────────────

type CreateUserRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=150"`
	Email        string `json:"email"    validate:"required,email"`
	Name         string `json:"name"     validate:"required,min=1,max=255"`
	Password     string `json:"password" validate:"required"` // Validasi min length dipindah ke custom validator
	IsActive     *bool  `json:"is_active"`
	IsStaff      *bool  `json:"is_staff"`
	IsSuperadmin *bool  `json:"is_superadmin"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"         validate:"omitempty,min=1,max=255"`
	Email    *string `json:"email"        validate:"omitempty,email"`
	Username *string `json:"username"     validate:"omitempty,min=3,max=150"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"` // Validasi min length dipindah ke custom validator
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UpdateSettingsRequest request untuk update settings user
type UpdateSettingsRequest struct {
	Settings []models.UserSetting `json:"settings" validate:"required"`
}

type DeleteUserRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

// ─── Filters          ──────────────────────────────────────────────────────────
type UserFilter struct {
	Name         string `query:"name"`
	Username     string `query:"username"`
	Email        string `query:"email"`
	IsSuperadmin *bool  `query:"is_superadmin"`
	IsActive     *bool  `query:"is_active"`
	IsStaff      *bool  `query:"is_staff"`
}

type UserDeletedFilter struct {
	Name     string `query:"name"`
	Username string `query:"username"`
	Email    string `query:"email"`
}

// ─── Default Settings ──────────────────────────────────────────────────────────

// DefaultUserSettings mengembalikan default settings untuk user baru
func DefaultUserSettings() []models.UserSetting {
	return []models.UserSetting{
		{
			Key:         "is_dark_mode",
			Type:        "boolean",
			Value:       false,
			Description: "Aktifkan tema gelap pada antarmuka",
		},
	}
}

// ─── Custom Validator Registration ─────────────────────────────────────────────

// RegisterPasswordValidation mendaftarkan rule validasi dinamis berdasarkan config.
// Panggil fungsi ini sekali saat inisialisasi aplikasi (misal di main.go atau di setup validator).
func RegisterPasswordPolicyValidation(v *validator.Validate, cfg *config.Config) error {
	return v.RegisterValidation("password_policy", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		// 1. Validasi Panjang Minimal
		if len(password) < cfg.PasswordMinLength {
			return false
		}

		var hasUpper, hasNumber, hasSymbol bool

		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSymbol = true
			}
		}

		// 2. Validasi Huruf Besar (jika diaktifkan di config)
		if cfg.PasswordRequireUppercase && !hasUpper {
			return false
		}

		// 3. Validasi Angka (jika diaktifkan di config)
		if cfg.PasswordRequireNumber && !hasNumber {
			return false
		}

		// 4. Validasi Simbol/Karakter Unik (jika diaktifkan di config)
		if cfg.PasswordRequireSymbol && !hasSymbol {
			return false
		}

		return true
	})
}

// ValidatePasswordPolicy mengecek string password terhadap policy yang aktif di config
func ValidatePasswordPolicy(password string, cfg *config.Config) map[string]string {
	errs := make(map[string]string)

	if len(password) < cfg.PasswordMinLength {
		errs["new_password"] = fmt.Sprintf("Password minimal harus %d karakter", cfg.PasswordMinLength)
		return errs // Langsung return jika panjang tidak memenuhi
	}

	var hasUpper, hasNumber, hasSymbol bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if cfg.PasswordRequireUppercase && !hasUpper {
		errs["new_password"] = "Password harus mengandung minimal satu huruf besar"
	}
	if cfg.PasswordRequireNumber && !hasNumber {
		errs["new_password"] = "Password harus mengandung minimal satu angka"
	}
	if cfg.PasswordRequireSymbol && !hasSymbol {
		errs["new_password"] = "Password harus mengandung minimal satu simbol/karakter khusus"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
