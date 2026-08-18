package seeders

import (
	"fmt"
	"log"

	"neosim_go/internal/modules/users/models"
	"neosim_go/internal/modules/users/tests/factories"

	"gorm.io/gorm"
)

type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

func (s *UserSeeder) Run() error {
	log.Println("🌱 Seeding users...")

	if err := s.seedSuperuser(); err != nil {
		return err
	}

	if err := s.seedStaff(); err != nil {
		return err
	}

	if err := s.seedRegularUsers(); err != nil {
		return err
	}

	if err := s.syncSequence(); err != nil {
		log.Printf("   ⚠️ Warning: Gagal sinkronisasi sequence: %v", err)
	}

	log.Println("✅ Users seeding selesai!")
	return nil
}

func (s *UserSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data user...")

	if err := s.db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}

	if err := s.db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}

	log.Println("✅ Data user dihapus.")
	return s.Run()
}

// ─── Seeder Methods ────────────────────────────────────────────────────────────

func (s *UserSeeder) seedSuperuser() error {
	// Khusus Superadmin: Set nama khusus "Super Admin" dan ID = 1
	user := factories.MakeSuperadminUser()
	user.ID = 1
	user.Name = "Super Admin" // 👈 Custom nama spesifik untuk Superadmin
	user.Username = "superadmin"

	var count int64
	s.db.Model(&models.User{}).Where("id = ? OR username = ?", user.ID, user.Username).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  Superuser '%s' (%s, ID: %d) sudah ada, skip.", user.Name, user.Username, user.ID)
		return nil
	}

	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	log.Printf("   ✅ Superuser '%s' (%s, ID: %d) dibuat.", user.Name, user.Username, user.ID)
	return nil
}

func (s *UserSeeder) seedStaff() error {
	staffCount := 50

	// ID 2, 3, 4
	for i := 1; i <= staffCount; i++ {
		targetID := int64(1 + i)
		user := factories.MakeStaffsUser(i)
		user.ID = targetID

		var count int64
		s.db.Model(&models.User{}).Where("id = ? OR username = ?", user.ID, user.Username).Count(&count)
		if count > 0 {
			log.Printf("   ⏭️  Staff '%s' (ID: %d) sudah ada, skip.", user.Username, user.ID)
			continue
		}

		if err := s.db.Create(user).Error; err != nil {
			return err
		}

		log.Printf("   ✅ Staff '%s' (ID: %d) dibuat.", user.Username, user.ID)
	}

	return nil
}

func (s *UserSeeder) seedRegularUsers() error {
	regularCount := 50
	startID := int64(5) // ID 5 s/d 14

	for i := 0; i < regularCount; i++ {
		targetID := startID + int64(i)
		idx := i + 1

		user := factories.NewUserFactory().
			With("id", targetID).
			With("username", fmt.Sprintf("user_%d", idx)).
			With("email", fmt.Sprintf("user_%d@example.com", idx)).
			With("name", fmt.Sprintf("User %d", idx)).
			Make()

		var count int64
		s.db.Model(&models.User{}).Where("id = ? OR username = ?", user.ID, user.Username).Count(&count)
		if count > 0 {
			log.Printf("   ⏭️  User '%s' (ID: %d) sudah ada, skip.", user.Username, user.ID)
			continue
		}

		if err := s.db.Create(user).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat user '%s': %v", user.Username, err)
			continue
		}
		log.Printf("   ✅ User '%s' (ID: %d) dibuat.", user.Username, user.ID)
	}

	return nil
}

func (s *UserSeeder) syncSequence() error {
	var maxID int64
	s.db.Model(&models.User{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)

	if maxID > 0 {
		query := fmt.Sprintf("SELECT setval('users_id_seq', %d, true)", maxID)
		return s.db.Exec(query).Error
	}
	return nil
}
