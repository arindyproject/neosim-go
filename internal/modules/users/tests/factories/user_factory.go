package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/users/models"

	"golang.org/x/crypto/bcrypt"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// UserFactory membuat data user untuk keperluan testing/seeding
type UserFactory struct {
	overrides map[string]interface{}
}

// NewUserFactory membuat instance UserFactory baru
func NewUserFactory() *UserFactory {
	return &UserFactory{
		overrides: make(map[string]interface{}),
	}
}

// With menambahkan override field sebelum build
func (f *UserFactory) With(field string, value interface{}) *UserFactory {
	f.overrides[field] = value
	return f
}

// Make membuat satu User model tanpa menyimpan ke DB
func (f *UserFactory) Make() *models.User {
	idx := rng.Intn(999999)

	username := fmt.Sprintf("user_%d", idx)
	email := fmt.Sprintf("user_%d@example.com", idx)
	name := fmt.Sprintf("User %d", idx)
	password := "password123"

	// Apply overrides
	if v, ok := f.overrides["username"]; ok {
		username = v.(string)
	}
	if v, ok := f.overrides["email"]; ok {
		email = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["password"]; ok {
		password = v.(string)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &models.User{
		Username:           username,
		Email:              email,
		Name:               name,
		Password:           string(hashed),
		IsActive:           true,
		IsVerified:         false,
		IsSuperuser:        false,
		IsStaff:            false,
		MustChangePassword: false,
	}

	// Apply boolean overrides
	if v, ok := f.overrides["is_active"]; ok {
		user.IsActive = v.(bool)
	}
	if v, ok := f.overrides["is_superuser"]; ok {
		user.IsSuperuser = v.(bool)
	}
	if v, ok := f.overrides["is_staff"]; ok {
		user.IsStaff = v.(bool)
	}
	if v, ok := f.overrides["is_verified"]; ok {
		user.IsVerified = v.(bool)
	}

	// Set default settings
	user.SetSettings(defaultSettings())

	return user
}

// MakeMany membuat banyak User model tanpa menyimpan ke DB
func (f *UserFactory) MakeMany(count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = NewUserFactory().Make()
	}
	return users
}

// ─── Preset Factories ──────────────────────────────────────────────────────────

// MakeSuperuser membuat user dengan role superuser
func MakeSuperuser() *models.User {
	return NewUserFactory().
		With("username", "superadmin").
		With("email", "superadmin@example.com").
		With("name", "Super Admin").
		With("is_superuser", true).
		With("is_staff", true).
		With("is_verified", true).
		Make()
}

// MakeStaff membuat user dengan role staff
func MakeStaff(idx int) *models.User {
	return NewUserFactory().
		With("username", fmt.Sprintf("staff_%d", idx)).
		With("email", fmt.Sprintf("staff_%d@example.com", idx)).
		With("name", fmt.Sprintf("Staff %d", idx)).
		With("is_staff", true).
		With("is_verified", true).
		Make()
}

// MakeInactiveUser membuat user yang tidak aktif
func MakeInactiveUser(idx int) *models.User {
	return NewUserFactory().
		With("username", fmt.Sprintf("inactive_%d", idx)).
		With("email", fmt.Sprintf("inactive_%d@example.com", idx)).
		With("is_active", false).
		Make()
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func defaultSettings() []models.UserSetting {
	return []models.UserSetting{
		{Key: "is_dark_mode", Type: "boolean", Value: false, Description: "Aktifkan tema gelap"},
		{Key: "language", Type: "string", Value: "id", Description: "Bahasa antarmuka"},
		{Key: "notification_email", Type: "boolean", Value: true, Description: "Notifikasi email"},
	}
}
