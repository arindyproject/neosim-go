package sorting

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Config mendefinisikan aturan sorting untuk satu entitas/modul.
type Config struct {
	// Allowed memetakan nilai sort_by dari client -> ekspresi kolom DB.
	// Hanya key di sini yang boleh dipakai (mencegah SQL injection lewat Order()).
	Allowed map[string]string

	DefaultColumn string // kolom default jika sort_by kosong/tidak valid (nilai DB, bukan key)
	DefaultDesc   bool   // arah default saat sort_by kosong/tidak valid
	TieBreaker    string // kolom pembeda agar urutan stabil antar halaman, mis. "id"
}

// Build menghasilkan klausa ORDER BY yang aman.
//   - sort_by kosong/tidak valid -> DefaultColumn + DefaultDesc
//   - sort_by valid              -> sort_order "desc" = DESC, selain itu ASC
func (c Config) Build(sortBy, sortOrder string) string {
	col, ok := c.Allowed[strings.ToLower(strings.TrimSpace(sortBy))]

	var desc bool
	if ok {
		desc = strings.EqualFold(strings.TrimSpace(sortOrder), "desc")
	} else {
		col = c.DefaultColumn
		desc = c.DefaultDesc
	}

	dir := "ASC"
	if desc {
		dir = "DESC"
	}

	order := fmt.Sprintf("%s %s", col, dir)
	if c.TieBreaker != "" && c.TieBreaker != col {
		order += fmt.Sprintf(", %s %s", c.TieBreaker, dir)
	}
	return order
}

// Scope dipakai langsung di GORM: db.Scopes(cfg.Scope(sortBy, sortOrder))
func (c Config) Scope(sortBy, sortOrder string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(c.Build(sortBy, sortOrder))
	}
}

// IsValid mengecek apakah sort_by ada di whitelist (berguna untuk respons 400 di handler).
func (c Config) IsValid(sortBy string) bool {
	if strings.TrimSpace(sortBy) == "" {
		return true // kosong = pakai default
	}
	_, ok := c.Allowed[strings.ToLower(strings.TrimSpace(sortBy))]
	return ok
}
