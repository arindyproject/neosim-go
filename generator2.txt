package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"
)

// ─── Config ────────────────────────────────────────────────────────────────────

type ModuleConfig struct {
	ModuleName    string // contoh: roles
	ModuleTitle   string // contoh: Role (PascalCase singular)
	ModulePlural  string // contoh: Roles (PascalCase plural)
	PackageName   string // contoh: roles
	ProjectModule string // contoh: neosim_go
	Timestamp     string // contoh: 20240507120000
}

// ─── Main ──────────────────────────────────────────────────────────────────────

func main() {
	name := flag.String("name", "", "Nama module (snake_case, contoh: login_history)")
	project := flag.String("project", "neosim_go", "Nama Go module project (dari go.mod)")
	flag.Parse()

	if *name == "" {
		log.Fatal("❌ Nama module wajib diisi. Contoh: go run ./cmd/gen/main.go -name=roles")
	}

	cfg := ModuleConfig{
		ModuleName:    *name,
		ModuleTitle:   toPascalCase(*name),
		ModulePlural:  toPascalCase(*name) + "s",
		PackageName:   toPackageName(*name),
		ProjectModule: *project,
		Timestamp:     time.Now().Format("20060102150405"),
	}

	basePath := filepath.Join("internal", "modules", cfg.ModuleName)

	fmt.Printf("\n🚀 Membuat module: %s\n", cfg.ModuleName)
	fmt.Printf("   Path: %s\n\n", basePath)

	files := buildFileList(cfg, basePath)
	for _, f := range files {
		if err := generateFile(f.path, f.tmpl, cfg); err != nil {
			log.Fatalf("❌ Gagal generate %s: %v", f.path, err)
		}
		fmt.Printf("   ✅ %s\n", f.path)
	}

	printNextSteps(cfg)
}

// ─── File List ─────────────────────────────────────────────────────────────────

type fileEntry struct {
	path string
	tmpl string
}

func buildFileList(cfg ModuleConfig, base string) []fileEntry {
	return []fileEntry{
		{filepath.Join(base, "contracts", "interfaces.go"), tmplContracts},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_request.go", cfg.ModuleName)), tmplRequest},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_response.go", cfg.ModuleName)), tmplResponse},
		{filepath.Join(base, "models", fmt.Sprintf("%s.go", cfg.ModuleName)), tmplModel},
		{filepath.Join(base, "repositories", fmt.Sprintf("%s_repository.go", cfg.ModuleName)), tmplRepository},
		{filepath.Join(base, "services", fmt.Sprintf("%s_service.go", cfg.ModuleName)), tmplService},
		{filepath.Join(base, "handlers", fmt.Sprintf("%s_handler.go", cfg.ModuleName)), tmplHandler},
		{filepath.Join(base, "migrations", fmt.Sprintf("%s_migrate.go", cfg.ModuleName)), tmplMigration},
		{filepath.Join(base, "migrations", fmt.Sprintf("001_create_%s_table.sql", cfg.ModuleName)), tmplSQL},
		{filepath.Join(base, "tests", "factories", fmt.Sprintf("%s_factory.go", cfg.ModuleName)), tmplFactory},
		{filepath.Join(base, "tests", "seeders", fmt.Sprintf("%s_seeder.go", cfg.ModuleName)), tmplSeeder},
		{filepath.Join(base, "tests", "helpers", "db_helper.go"), tmplDBHelper},
		{filepath.Join(base, "module.go"), tmplModule},
		{filepath.Join(base, "routes.go"), tmplRoutes},
		{filepath.Join(base, "register.go"), tmplRegister},
	}
}

// ─── Generator ─────────────────────────────────────────────────────────────────

func generateFile(path, tmplStr string, cfg ModuleConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("   ⏭️  Skip (sudah ada): %s\n", path)
		return nil
	}

	tmpl, err := template.New(path).Parse(tmplStr)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cfg)
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}
	return result.String()
}

func toPackageName(s string) string {
	return strings.ReplaceAll(s, "_", "")
}

// ─── Next Steps ────────────────────────────────────────────────────────────────

func printNextSteps(cfg ModuleConfig) {
	fmt.Printf(`
────────────────────────────────────────────────────────
✅ Module '%s' berhasil dibuat!

📋 Langkah selanjutnya:

1. Tambahkan blank import di internal/apps/apps.go:
   _ "%s/internal/modules/%s"

2. Tambahkan blank import di cmd/migrate/main.go:
   _ "%s/internal/modules/%s"

3. Edit model di:
   internal/modules/%s/models/%s.go

4. Jalankan migrasi:
   make migrate-dev

5. Jalankan seeder:
   make seed
────────────────────────────────────────────────────────
`,
		cfg.ModuleName,
		cfg.ProjectModule, cfg.ModuleName,
		cfg.ProjectModule, cfg.ModuleName,
		cfg.ModuleName, cfg.ModuleName,
	)
}

// ─── Templates ─────────────────────────────────────────────────────────────────

var tmplContracts = `package contracts

import (
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.{{.ModuleTitle}}) error
	GetByID(id int64) (*models.{{.ModuleTitle}}, error)
	List(page, pageSize int) ([]models.{{.ModuleTitle}}, int64, error)
	Update(m *models.{{.ModuleTitle}}) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.Create{{.ModuleTitle}}Request, createdBy *int64) (*dto.{{.ModuleTitle}}Response, error)
	GetByID(id int64) (*dto.{{.ModuleTitle}}Response, error)
	List(page, pageSize int) ([]dto.{{.ModuleTitle}}Response, int64, error)
	Update(id int64, req *dto.Update{{.ModuleTitle}}Request, updatedBy *int64) (*dto.{{.ModuleTitle}}Response, error)
	Delete(id int64) error
}
`

var tmplRequest = `package dto

// Create{{.ModuleTitle}}Request request body untuk membuat {{.ModuleName}} baru
type Create{{.ModuleTitle}}Request struct {
	Name        string  ` + "`" + `json:"name" validate:"required,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}

// Update{{.ModuleTitle}}Request request body untuk update {{.ModuleName}}
type Update{{.ModuleTitle}}Request struct {
	Name        *string ` + "`" + `json:"name" validate:"omitempty,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}
`

var tmplResponse = `package dto

import (
	"time"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
)

// {{.ModuleTitle}}Response response untuk single {{.ModuleName}}
type {{.ModuleTitle}}Response struct {
	ID          int64     ` + "`" + `json:"id"` + "`" + `
	Name        string    ` + "`" + `json:"name"` + "`" + `
	Description *string   ` + "`" + `json:"description"` + "`" + `
	CreatedBy   *int64    ` + "`" + `json:"created_by"` + "`" + `
	UpdatedBy   *int64    ` + "`" + `json:"updated_by"` + "`" + `
	CreatedAt   time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

// To{{.ModuleTitle}}Response mengubah model menjadi response
func To{{.ModuleTitle}}Response(m *models.{{.ModuleTitle}}) *{{.ModuleTitle}}Response {
	return &{{.ModuleTitle}}Response{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// To{{.ModuleTitle}}ListResponse mengubah slice model menjadi slice response
func To{{.ModuleTitle}}ListResponse(items []models.{{.ModuleTitle}}) []{{.ModuleTitle}}Response {
	var responses []{{.ModuleTitle}}Response
	for _, m := range items {
		responses = append(responses, *To{{.ModuleTitle}}Response(&m))
	}
	return responses
}
`

var tmplModel = `package models

import (
	"time"

	"gorm.io/gorm"
)

// {{.ModuleTitle}} represents the {{.ModuleName}}s table in database
type {{.ModuleTitle}} struct {
	ID          int64          ` + "`" + `gorm:"primaryKey;autoIncrement;column:id" json:"id"` + "`" + `
	Name        string         ` + "`" + `gorm:"column:name;type:varchar(255);not null" json:"name"` + "`" + `
	Description *string        ` + "`" + `gorm:"column:description;type:text" json:"description"` + "`" + `
	CreatedBy   *int64         ` + "`" + `gorm:"column:created_by" json:"created_by"` + "`" + `
	UpdatedBy   *int64         ` + "`" + `gorm:"column:updated_by" json:"updated_by"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"` + "`" + `
	DeletedAt   gorm.DeletedAt ` + "`" + `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"` + "`" + `
}

func ({{.ModuleTitle}}) TableName() string {
	return "{{.ModuleName}}s"
}
`

var tmplRepository = `package repositories

import (
	"errors"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/contracts"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// New{{.ModuleTitle}}Repository membuat instance repository baru
func New{{.ModuleTitle}}Repository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *models.{{.ModuleTitle}}) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id int64) (*models.{{.ModuleTitle}}, error) {
	var m models.{{.ModuleTitle}}
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) List(page, pageSize int) ([]models.{{.ModuleTitle}}, int64, error) {
	var items []models.{{.ModuleTitle}}
	var total int64
	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.{{.ModuleTitle}}{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("deleted_at IS NULL").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *repository) Update(m *models.{{.ModuleTitle}}) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.{{.ModuleTitle}}{}).Error
}
`

var tmplService = `package services

import (
	"errors"
	"time"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/contracts"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
)

type service struct {
	repo contracts.Repository
}

// New{{.ModuleTitle}}Service membuat instance service baru
func New{{.ModuleTitle}}Service(repo contracts.Repository) contracts.Service {
	return &service{repo: repo}
}

func (s *service) Create(req *dto.Create{{.ModuleTitle}}Request, createdBy *int64) (*dto.{{.ModuleTitle}}Response, error) {
	m := &models.{{.ModuleTitle}}{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.To{{.ModuleTitle}}Response(m), nil
}

func (s *service) GetByID(id int64) (*dto.{{.ModuleTitle}}Response, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ModuleName}} tidak ditemukan")
	}
	return dto.To{{.ModuleTitle}}Response(m), nil
}

func (s *service) List(page, pageSize int) ([]dto.{{.ModuleTitle}}Response, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return dto.To{{.ModuleTitle}}ListResponse(items), total, nil
}

func (s *service) Update(id int64, req *dto.Update{{.ModuleTitle}}Request, updatedBy *int64) (*dto.{{.ModuleTitle}}Response, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ModuleName}} tidak ditemukan")
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
	return dto.To{{.ModuleTitle}}Response(m), nil
}

func (s *service) Delete(id int64) error {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("{{.ModuleName}} tidak ditemukan")
	}
	return s.repo.Delete(id)
}
`

// ✅ Fix: c echo.Context (bukan *echo.Context), c.PathParam (bukan c.Param)
var tmplHandler = `package handlers

import (
	"net/http"
	"strconv"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/contracts"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/dto"
	"{{.ProjectModule}}/internal/shared/response"
	"{{.ProjectModule}}/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// {{.ModuleTitle}}Handler defines HTTP handlers
type {{.ModuleTitle}}Handler struct {
	service contracts.Service
}

// New{{.ModuleTitle}}Handler membuat instance handler baru
func New{{.ModuleTitle}}Handler(service contracts.Service) *{{.ModuleTitle}}Handler {
	return &{{.ModuleTitle}}Handler{service: service}
}

// ─── Private Helpers ───────────────────────────────────────────────────────────

func parseID(c echo.Context) (int64, error) {
	return strconv.ParseInt(c.PathParam("id"), 10, 64)
}

func getActorID(c echo.Context) *int64 {
	if userID, ok := c.Get("userID").(int64); ok {
		return &userID
	}
	return nil
}

// ─── Handlers ──────────────────────────────────────────────────────────────────

// List handles GET /api/v1/{{.ModuleName}}s
func (h *{{.ModuleTitle}}Handler) List(c echo.Context) error {
	page, pageSize := 1, 10

	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	items, total, err := h.service.List(page, pageSize)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// GetByID handles GET /api/v1/{{.ModuleName}}s/:id
func (h *{{.ModuleTitle}}Handler) GetByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// Create handles POST /api/v1/{{.ModuleName}}s
func (h *{{.ModuleTitle}}Handler) Create(c echo.Context) error {
	var req dto.Create{{.ModuleTitle}}Request
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	item, err := h.service.Create(&req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// Update handles PUT /api/v1/{{.ModuleName}}s/:id
func (h *{{.ModuleTitle}}Handler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.Update{{.ModuleTitle}}Request
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	item, err := h.service.Update(id, &req, getActorID(c))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "{{.ModuleName}} tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// Delete handles DELETE /api/v1/{{.ModuleName}}s/:id
func (h *{{.ModuleTitle}}Handler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	if err := h.service.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "{{.ModuleName}} tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
`

var tmplMigration = `package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"

	"gorm.io/gorm"
)

//go:embed 001_create_{{.ModuleName}}_table.sql
var {{.PackageName}}SQL string

// Migrate{{.ModuleTitle}} menjalankan GORM auto-migration
func Migrate{{.ModuleTitle}}(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.{{.ModuleTitle}}{})
}

// Migrate{{.ModuleTitle}}WithSQL menjalankan migrasi via raw SQL
func Migrate{{.ModuleTitle}}WithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec({{.PackageName}}SQL)
	if err != nil {
		log.Printf("Error creating {{.ModuleName}}s table: %v", err)
		return err
	}
	log.Println("{{.ModuleTitle}}s table migrated successfully")
	return nil
}

// Drop{{.ModuleTitle}}Table menghapus tabel (gunakan dengan hati-hati!)
func Drop{{.ModuleTitle}}Table(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.{{.ModuleTitle}}{})
}
`

// ✅ Fix: hapus trailing comma di SQL
var tmplSQL = `-- Migration: Create {{.ModuleName}}s table
-- Timestamp: {{.Timestamp}}

CREATE TABLE IF NOT EXISTS {{.ModuleName}}s (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_{{.ModuleName}}s_deleted_at ON {{.ModuleName}}s(deleted_at);
`

var tmplFactory = `package factories

import (
	"fmt"
	"math/rand"
	"time"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// {{.ModuleTitle}}Factory membuat data {{.ModuleName}} untuk testing/seeding
type {{.ModuleTitle}}Factory struct {
	overrides map[string]interface{}
}

func New{{.ModuleTitle}}Factory() *{{.ModuleTitle}}Factory {
	return &{{.ModuleTitle}}Factory{overrides: make(map[string]interface{})}
}

func (f *{{.ModuleTitle}}Factory) With(field string, value interface{}) *{{.ModuleTitle}}Factory {
	f.overrides[field] = value
	return f
}

func (f *{{.ModuleTitle}}Factory) Make() *models.{{.ModuleTitle}} {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("{{.ModuleTitle}} %d", idx)
	desc := fmt.Sprintf("Deskripsi {{.ModuleName}} %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.{{.ModuleTitle}}{
		Name:        name,
		Description: &desc,
	}
}

func (f *{{.ModuleTitle}}Factory) MakeMany(count int) []*models.{{.ModuleTitle}} {
	items := make([]*models.{{.ModuleTitle}}, count)
	for i := 0; i < count; i++ {
		items[i] = New{{.ModuleTitle}}Factory().Make()
	}
	return items
}
`

var tmplSeeder = `package seeders

import (
	"log"

	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/tests/factories"

	"gorm.io/gorm"
)

// {{.ModuleTitle}}Seeder mengelola seeding data {{.ModuleName}}
type {{.ModuleTitle}}Seeder struct {
	db *gorm.DB
}

func New{{.ModuleTitle}}Seeder(db *gorm.DB) *{{.ModuleTitle}}Seeder {
	return &{{.ModuleTitle}}Seeder{db: db}
}

// Run menjalankan seeder
func (s *{{.ModuleTitle}}Seeder) Run() error {
	log.Println("🌱 Seeding {{.ModuleName}}s...")

	items := factories.New{{.ModuleTitle}}Factory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat {{.ModuleName}}: %v", err)
			continue
		}
		log.Printf("   ✅ {{.ModuleTitle}} '%s' dibuat.", item.Name)
	}

	log.Println("✅ {{.ModuleTitle}}s seeding selesai!")
	return nil
}

// Fresh menghapus semua data lalu seed ulang
func (s *{{.ModuleTitle}}Seeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data {{.ModuleName}}s...")

	if err := s.db.Exec("DELETE FROM {{.ModuleName}}s").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE {{.ModuleName}}s_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault menyimpan satu item jika belum ada
func (s *{{.ModuleTitle}}Seeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.{{.ModuleTitle}}{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}

	item := factories.New{{.ModuleTitle}}Factory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}

	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
`

var tmplDBHelper = `package helpers

import (
	"log"

	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"

	"gorm.io/gorm"
)

// SetupTestDB membuat koneksi DB untuk keperluan test
func SetupTestDB() *gorm.DB {
	cfg := config.LoadConfig("DEV")
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi DB untuk test:", err)
	}
	return db
}

// MigrateTestDB menjalankan migrasi untuk test DB
func MigrateTestDB(db *gorm.DB) {
	if err := db.AutoMigrate(&models.{{.ModuleTitle}}{}); err != nil {
		log.Fatal("Gagal migrasi test DB:", err)
	}
}

// TruncateTable menghapus semua record dan reset sequence
func TruncateTable(db *gorm.DB, tables ...string) {
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Printf("Warning: Gagal truncate table %s: %v", table, err)
		}
	}
}
`

// ✅ Fix: tambah jwtManager ke Module
var tmplModule = `package {{.PackageName}}

import (
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/handlers"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/repositories"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/services"
	"{{.ProjectModule}}/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module mewakili {{.ModuleName}} module
type Module struct {
	db         *gorm.DB
	handler    *handlers.{{.ModuleTitle}}Handler
	jwtManager *utils.JWTManager
}

// NewModule membuat instance module baru dan wire semua layer
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager) *Module {
	repo := repositories.New{{.ModuleTitle}}Repository(db)
	svc := services.New{{.ModuleTitle}}Service(repo)
	handler := handlers.New{{.ModuleTitle}}Handler(svc)

	return &Module{
		db:         db,
		handler:    handler,
		jwtManager: jwtManager,
	}
}

// InitRoutes mendaftarkan routes ke echo instance
func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager)
}
`

// ✅ Fix: tambah JWT middleware di routes
var tmplRoutes = `package {{.PackageName}}

import (
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/handlers"
	"{{.ProjectModule}}/internal/modules/auth/middlewares"
	"{{.ProjectModule}}/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mendaftarkan semua routes untuk module {{.ModuleName}}
func RegisterRoutes(e *echo.Echo, h *handlers.{{.ModuleTitle}}Handler, jwtManager *utils.JWTManager) {
	g := e.Group("/api/v1/{{.ModuleName}}s", middlewares.JWTMiddleware(jwtManager))
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
`

// ✅ Fix: tambah SetConfig dan build jwtManager dari config
var tmplRegister = `package {{.PackageName}}

import (
	"database/sql"

	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/apps"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/migrations"
	"{{.ProjectModule}}/internal/modules/{{.ModuleName}}/models"
	"{{.ProjectModule}}/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db  *gorm.DB
	cfg *config.Config
}

// init dipanggil otomatis saat package di-import (blank import)
func init() {
	apps.Register(&registryModule{})
}

func (r *registryModule) SetDB(db *gorm.DB)            { r.db = db }
func (r *registryModule) SetConfig(cfg *config.Config) { r.cfg = cfg }

func (r *registryModule) InitRoutes(e *echo.Echo) {
	jwtManager := utils.NewJWTManager(
		r.cfg.JWTSecret,
		r.cfg.JWTIssuer,
		r.cfg.JWTAccessTokenExpMinutes,
		r.cfg.JWTRefreshTokenExpDays,
	)
	NewModule(r.db, jwtManager).InitRoutes(e)
}

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.{{.ModuleTitle}}{},
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return nil
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	return migrations.Migrate{{.ModuleTitle}}WithSQL(sqlDB)
}
`
