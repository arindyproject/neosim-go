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

// ─── Configs ───────────────────────────────────────────────────────────────────

type ModuleConfig struct {
	MainModule       string
	SubModule        string
	ModuleName       string
	ModuleTitle      string
	MethodSuffix     string // Pascal(ModuleName) — suffix method: Create{{MethodSuffix}}, Get{{MethodSuffix}}ByID, dst.
	ModulePlural     string
	PackageName      string
	ProjectModule    string
	Timestamp        string
	URLPrefix        string
	URLPrefixOpenAPI string
	TableName        string
}

// AddItemConfig dipakai pada Mode 3 (-sub=... -add=...).
// Item BARU TIDAK memiliki struct/handler/repository/service sendiri.
// Semua method item ditempelkan ke struct `service`, `repository`, dan
// `{{SubModuleTitle}}Handler` yang sudah dibuat pada Mode 1/2 untuk sub-module
// terkait. Interface tetap didefinisikan di file terpisah (mis. tag_interfaces.go)
// lalu di-embed otomatis ke dalam interface Repository/Service utama.
type AddItemConfig struct {
	MainModule       string
	SubModule        string
	SubModuleTitle   string // Pascal title milik sub-module existing, mis. "ArtikelKategori" -> dipakai untuk nama Handler/Mock/Suite yang sudah ada
	ItemName         string
	ItemTitle        string
	PackageName      string
	ProjectModule    string
	Timestamp        string
	URLPrefix        string
	URLPrefixOpenAPI string
	TableName        string
}

// ─── Main ──────────────────────────────────────────────────────────────────────

func main() {
	name := flag.String("name", "", "Nama module utama (snake_case, contoh: artikel)")
	add := flag.String("add", "", "Sub-module ATAU item baru (tergantung apakah -sub diisi)")
	sub := flag.String("sub", "", "Folder sub-module existing (hanya untuk mode add-item)")
	project := flag.String("project", "neosim_go", "Nama Go module project (dari go.mod)")
	flag.Parse()

	if *name == "" {
		log.Fatal("❌ Flag -name wajib diisi.\n" +
			"   Mode 1 (main module):  go run generator.go -name=artikel\n" +
			"   Mode 2 (sub module):   go run generator.go -name=artikel -add=kategori\n" +
			"   Mode 3 (add item):     go run generator.go -name=artikel -sub=kategori -add=tag")
	}

	// ── Mode 3: add item ke sub-module existing ─────────────────────────────
	if *sub != "" {
		if *add == "" {
			log.Fatal("❌ Flag -add wajib diisi saat menggunakan -sub.\n" +
				"   Contoh: go run generator.go -name=artikel -sub=kategori -add=tag")
		}
		runAddItem(*name, *sub, *add, *project)
		return
	}

	// ── Mode 1 & 2: generate module / sub-module ────────────────────────────
	runGenerateModule(*name, *add, *project)
}

// ═══════════════════════════════════════════════════════════════════════════════
// MODE 1 & 2 — Generate Module / Sub-Module
// ═══════════════════════════════════════════════════════════════════════════════

func runGenerateModule(name, add, project string) {
	subModule := name
	if add != "" {
		subModule = add
	}

	mainPascal := toPascalCase(name)
	subPascal := toPascalCase(subModule)

	entityTitle := mainPascal
	if add != "" {
		entityTitle = mainPascal + subPascal
	}

	urlPrefix := fmt.Sprintf("/api/v1/%s", name)
	urlPrefixOpenAPI := fmt.Sprintf("/%s", name)
	if add != "" {
		urlPrefix = fmt.Sprintf("/api/v1/%s/%s", name, add)
		urlPrefixOpenAPI = fmt.Sprintf("/%s/%s", name, add)
	}

	tableName := name + "s"
	if add != "" {
		tableName = name + "_" + add + "s"
	}

	cfg := ModuleConfig{
		MainModule:       name,
		SubModule:        subModule,
		ModuleName:       subModule,
		ModuleTitle:      entityTitle,
		MethodSuffix:     subPascal, // Create{{MethodSuffix}}, Get{{MethodSuffix}}ByID, dst — konsisten dgn pola item (Tag)
		ModulePlural:     entityTitle + "s",
		PackageName:      toPackageName(subModule),
		ProjectModule:    project,
		Timestamp:        time.Now().Format("20060102150405"),
		URLPrefix:        urlPrefix,
		URLPrefixOpenAPI: urlPrefixOpenAPI,
		TableName:        tableName,
	}

	basePath := filepath.Join("internal", "modules", cfg.MainModule, cfg.SubModule)

	fmt.Printf("\n🚀 Membuat module: %s/%s\n", cfg.MainModule, cfg.SubModule)
	fmt.Printf("   Path: %s\n\n", basePath)

	for _, f := range buildModuleFileList(cfg, basePath) {
		if err := generateModuleFile(f.path, f.tmpl, cfg); err != nil {
			log.Fatalf("❌ Gagal generate %s: %v", f.path, err)
		}
		fmt.Printf("   ✅ %s\n", f.path)
	}

	printModuleNextSteps(cfg)
}

type fileEntry struct {
	path string
	tmpl string
}

func buildModuleFileList(cfg ModuleConfig, base string) []fileEntry {
	return []fileEntry{
		{filepath.Join(base, "contracts", "interfaces.go"), tmplContracts},
		{filepath.Join(base, "contracts", fmt.Sprintf("%s_interfaces.go", cfg.ModuleName)), tmplModuleInterfaces},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_request.go", cfg.ModuleName)), tmplRequest},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_response.go", cfg.ModuleName)), tmplResponse},
		{filepath.Join(base, "models", fmt.Sprintf("%s.go", cfg.ModuleName)), tmplModel},
		{filepath.Join(base, "repositories", fmt.Sprintf("%s_repository.go", cfg.ModuleName)), tmplRepository},
		{filepath.Join(base, "repositories", "repository.go"), tmplMainRepository},
		{filepath.Join(base, "services", fmt.Sprintf("%s_service.go", cfg.ModuleName)), tmplService},
		{filepath.Join(base, "services", "service.go"), tmplMainService},
		{filepath.Join(base, "services", "permission.go"), tmplPermissionService},
		{filepath.Join(base, "handlers", fmt.Sprintf("%s_handler.go", cfg.ModuleName)), tmplHandler},
		{filepath.Join(base, "handlers", "handler.go"), tmplMainHandler},
		{filepath.Join(base, "migrations", fmt.Sprintf("%s_migrate.go", cfg.ModuleName)), tmplMigration},
		{filepath.Join(base, "migrations", fmt.Sprintf("001_create_%s_table.sql", cfg.TableName)), tmplSQL},
		{filepath.Join(base, "tests", "factories", fmt.Sprintf("%s_factory.go", cfg.ModuleName)), tmplFactory},
		{filepath.Join(base, "tests", "seeders", fmt.Sprintf("%s_seeder.go", cfg.ModuleName)), tmplSeeder},
		{filepath.Join(base, "tests", "helpers", "db_helper.go"), tmplDBHelper},
		{filepath.Join(base, "tests", "mocks", fmt.Sprintf("%s_repository_mock.go", cfg.ModuleName)), tmplModuleServiceMock},
		{filepath.Join(base, "tests", "mocks", "rbac_repository_mock.go"), tmplRBACMock},
		{filepath.Join(base, "tests", "mocks", "auth_repository_mock.go"), tmplAuthMock},
		{filepath.Join(base, "tests", "mocks", "user_repository_mock.go"), tmplUserMock},
		{filepath.Join(base, "tests", fmt.Sprintf("%s_service_test.go", cfg.ModuleName)), tmplModuleServiceTest},
		{filepath.Join(base, "module.go"), tmplModule},
		{filepath.Join(base, "routes.go"), tmplRoutes},
		{filepath.Join(base, "register.go"), tmplRegister},
	}
}

func generateModuleFile(path, tmplStr string, cfg ModuleConfig) error {
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

func printModuleNextSteps(cfg ModuleConfig) {
	fmt.Printf(`
────────────────────────────────────────────────────────
✅ Module '%s/%s' berhasil dibuat!

📋 Langkah selanjutnya:

1. Tambahkan blank import di internal/apps/apps.go:
   _ "%s/internal/modules/%s/%s"

2. Tambahkan blank import di cmd/migrate/main.go:
   _ "%s/internal/modules/%s/%s"

3. Edit model di:
   internal/modules/%s/%s/models/%s.go

4. Jalankan migrasi:
   make migrate-dev

5. Jalankan seeder:
   make seed

💡 Ingin menambah entitas anak (mis. tag) di dalam module ini tanpa
   membuat struct/handler/service terpisah? Gunakan mode add-item:
   go run generator.go -name=%s -sub=%s -add=<item>
────────────────────────────────────────────────────────
`,
		cfg.MainModule, cfg.SubModule,
		cfg.ProjectModule, cfg.MainModule, cfg.SubModule,
		cfg.ProjectModule, cfg.MainModule, cfg.SubModule,
		cfg.MainModule, cfg.SubModule, cfg.ModuleName,
		cfg.MainModule, cfg.SubModule,
	)
}

// ═══════════════════════════════════════════════════════════════════════════════
// MODE 3 — Add Item ke Sub-Module Existing (SHARED STRUCT, tanpa struct baru)
// ═══════════════════════════════════════════════════════════════════════════════
//
// Prinsip:
//   - Interface item (mis. TagRepository / TagService) didefinisikan di file
//     terpisah (contracts/tag_interfaces.go) dengan method bersuffix nama item
//     (CreateTag, GetTagByID, dst) agar tidak bentrok ketika di-embed.
//   - Method CRUD item ditempelkan ke struct `repository` & `service` yang
//     SAMA dengan sub-module induknya (bukan struct baru). Jadi tidak perlu
//     tagRepo/tagService terpisah — cukup satu instance yang sudah ada.
//   - Interface Repository & Service utama (contracts/interfaces.go) di-embed
//     otomatis dengan {{ModuleTitle}}Repository/{{ModuleTitle}}Service (method
//     utama) dan TagRepository/TagService (item tambahan) lewat marker komentar
//     "// GEN:ITEM_REPOSITORY_INTERFACE" & "// GEN:ITEM_SERVICE_INTERFACE".
//   - routes.go, register.go (Models & MigrateSQL) juga di-update otomatis
//     lewat marker, sehingga TIDAK ADA lagi langkah wiring manual.

func runAddItem(name, sub, add, project string) {
	mainPascal := toPascalCase(name)
	subPascal := toPascalCase(sub)
	itemPascal := toPascalCase(add)

	subModuleTitle := mainPascal
	if sub != name {
		subModuleTitle = mainPascal + subPascal
	}

	cfg := AddItemConfig{
		MainModule:       name,
		SubModule:        sub,
		SubModuleTitle:   subModuleTitle,
		ItemName:         add,
		ItemTitle:        itemPascal,
		PackageName:      toPackageName(sub),
		ProjectModule:    project,
		Timestamp:        time.Now().Format("20060102150405"),
		URLPrefix:        fmt.Sprintf("/api/v1/%s/%s/%ss", name, sub, add),
		URLPrefixOpenAPI: fmt.Sprintf("/%s/%s/%ss", name, sub, add),
		TableName:        fmt.Sprintf("%s_%s_%ss", name, sub, add),
	}

	basePath := filepath.Join("internal", "modules", cfg.MainModule, cfg.SubModule)

	fmt.Printf("\n🚀 Menambahkan item '%s' ke dalam %s/%s (shared struct, tanpa file service/repository baru)\n", cfg.ItemName, cfg.MainModule, cfg.SubModule)
	fmt.Printf("   Path: %s\n\n", basePath)

	for _, f := range buildItemFileList(cfg, basePath) {
		if err := generateItemFile(f.path, f.tmpl, cfg); err != nil {
			log.Fatalf("❌ Gagal generate %s: %v", f.path, err)
		}
		fmt.Printf("   ✅ %s\n", f.path)
	}

	fmt.Println()
	if err := updateSharedFiles(cfg, basePath); err != nil {
		log.Fatalf("❌ Gagal update file bersama: %v", err)
	}

	printItemNextSteps(cfg)
}

func buildItemFileList(cfg AddItemConfig, base string) []fileEntry {
	ts := cfg.Timestamp
	return []fileEntry{
		{filepath.Join(base, "contracts", fmt.Sprintf("%s_interfaces.go", cfg.ItemName)), tmplItemContracts},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_request.go", cfg.ItemName)), tmplItemRequest},
		{filepath.Join(base, "dto", fmt.Sprintf("%s_response.go", cfg.ItemName)), tmplItemResponse},
		{filepath.Join(base, "models", fmt.Sprintf("%s.go", cfg.ItemName)), tmplItemModel},
		{filepath.Join(base, "repositories", fmt.Sprintf("%s_repository.go", cfg.ItemName)), tmplItemRepository},
		{filepath.Join(base, "services", fmt.Sprintf("%s_service.go", cfg.ItemName)), tmplItemService},
		{filepath.Join(base, "services", fmt.Sprintf("%s_permission.go", cfg.ItemName)), tmplItemPermission},
		{filepath.Join(base, "handlers", fmt.Sprintf("%s_handler.go", cfg.ItemName)), tmplItemHandler},
		{filepath.Join(base, "migrations", fmt.Sprintf("%s_migrate.go", cfg.ItemName)), tmplItemMigration},
		{filepath.Join(base, "migrations", fmt.Sprintf("%s_create_%s_table.sql", ts, cfg.TableName)), tmplItemSQL},
		{filepath.Join(base, "tests", "factories", fmt.Sprintf("%s_factory.go", cfg.ItemName)), tmplItemFactory},
		{filepath.Join(base, "tests", "seeders", fmt.Sprintf("%s_seeder.go", cfg.ItemName)), tmplItemSeeder},
		{filepath.Join(base, "tests", "mocks", fmt.Sprintf("%s_repository_mock.go", cfg.ItemName)), tmplItemRepositoryMock},
		{filepath.Join(base, "tests", fmt.Sprintf("%s_service_test.go", cfg.ItemName)), tmplItemServiceTest},
	}
}

func generateItemFile(path, tmplStr string, cfg AddItemConfig) error {
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

// updateSharedFiles menyuntikkan embed interface + wiring routes/models/migration
// ke file-file yang sudah ada, lewat marker komentar "// GEN:...".
// Ini menggantikan seluruh langkah manual yang sebelumnya dicetak di
// printItemNextSteps.
func updateSharedFiles(cfg AddItemConfig, basePath string) error {
	type patch struct {
		file   string
		marker string
		insert string
	}

	itemRoutesBlock := fmt.Sprintf(
		"g%s := e.Group(\"%s\", jwt)\n"+
			"\tg%s.GET(\"\", h.List%s)\n"+
			"\tg%s.GET(\"/:id\", h.Get%sByID)\n"+
			"\tg%s.POST(\"\", h.Create%s)\n"+
			"\tg%s.PUT(\"/:id\", h.Update%s)\n"+
			"\tg%s.DELETE(\"/:id\", h.Delete%s)\n"+
			"\t// GEN:ITEM_ROUTES",
		cfg.ItemTitle, cfg.URLPrefix,
		cfg.ItemTitle, cfg.ItemTitle,
		cfg.ItemTitle, cfg.ItemTitle,
		cfg.ItemTitle, cfg.ItemTitle,
		cfg.ItemTitle, cfg.ItemTitle,
		cfg.ItemTitle, cfg.ItemTitle,
	)

	migrationBlock := fmt.Sprintf(
		"if err := migrations.Migrate%sWithSQL(sqlDB); err != nil {\n"+
			"\t\treturn err\n"+
			"\t}\n"+
			"\t// GEN:ITEM_MIGRATIONS",
		cfg.ItemTitle,
	)

	patches := []patch{
		{
			file:   filepath.Join(basePath, "contracts", "interfaces.go"),
			marker: "// GEN:ITEM_REPOSITORY_INTERFACE",
			insert: fmt.Sprintf("%sRepository\n\t// GEN:ITEM_REPOSITORY_INTERFACE", cfg.ItemTitle),
		},
		{
			file:   filepath.Join(basePath, "contracts", "interfaces.go"),
			marker: "// GEN:ITEM_SERVICE_INTERFACE",
			insert: fmt.Sprintf("%sService\n\t// GEN:ITEM_SERVICE_INTERFACE", cfg.ItemTitle),
		},
		{
			file:   filepath.Join(basePath, "routes.go"),
			marker: "// GEN:ITEM_ROUTES",
			insert: itemRoutesBlock,
		},
		{
			file:   filepath.Join(basePath, "register.go"),
			marker: "// GEN:ITEM_MODELS",
			insert: fmt.Sprintf("&models.%s{},\n\t\t// GEN:ITEM_MODELS", cfg.ItemTitle),
		},
		{
			file:   filepath.Join(basePath, "register.go"),
			marker: "// GEN:ITEM_MIGRATIONS",
			insert: migrationBlock,
		},
	}

	anyFailed := false
	for _, p := range patches {
		if err := insertBeforeMarker(p.file, p.marker, p.insert); err != nil {
			anyFailed = true
			fmt.Printf("   ⚠️  %v\n", err)
		} else {
			fmt.Printf("   ✅ %s (updated)\n", p.file)
		}
	}
	if anyFailed {
		fmt.Println("   ⚠️  Sebagian file bersama tidak ter-update otomatis. Cek marker \"// GEN:...\" di atas secara manual.")
	}
	return nil
}

// insertBeforeMarker mengganti baris marker dengan (insertion + marker),
// sehingga marker tetap ada untuk pemanggilan add-item berikutnya.
func insertBeforeMarker(path, marker, insertion string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("gagal baca %s: %w", path, err)
	}
	content := string(data)
	if strings.Contains(content, insertion) {
		return fmt.Errorf("skip (sudah ter-update sebelumnya): %s", path)
	}
	if !strings.Contains(content, marker) {
		return fmt.Errorf("marker %q tidak ditemukan di %s — tambahkan manual", marker, path)
	}
	content = strings.Replace(content, marker, insertion, 1)
	return os.WriteFile(path, []byte(content), 0644)
}

func printItemNextSteps(cfg AddItemConfig) {
	fmt.Printf(`
────────────────────────────────────────────────────────
✅ Item '%s' berhasil ditambahkan ke %s/%s!

Semua wiring (interface embed, route, model, migration) sudah
di-update otomatis lewat marker "// GEN:...". Tidak ada lagi
struct/service/repository terpisah untuk '%s' — semuanya menempel
ke struct 'service' & 'repository' milik %s/%s, dan ke handler
%sHandler yang sudah ada.

📋 Yang tersisa (opsional):

1. Jalankan migrasi:
   make migrate-dev

2. Panggil seeders.New%sSeeder(db).Run() dari seed runner-mu bila perlu,
   lalu jalankan test:
   go test ./internal/modules/%s/%s/...

3. Review validasi di dto/%s_request.go sesuai kebutuhan bisnis.
────────────────────────────────────────────────────────
`,
		cfg.ItemName, cfg.MainModule, cfg.SubModule,
		cfg.ItemName,
		cfg.MainModule, cfg.SubModule,
		cfg.SubModuleTitle,
		cfg.ItemTitle,
		cfg.MainModule, cfg.SubModule,
		cfg.ItemName,
	)
}

// ═══════════════════════════════════════════════════════════════════════════════
// HELPERS
// ═══════════════════════════════════════════════════════════════════════════════

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

// ═══════════════════════════════════════════════════════════════════════════════
// TEMPLATES — MODE 1 & 2
// ═══════════════════════════════════════════════════════════════════════════════

// tmplContracts (interfaces.go) HANYA berfungsi sebagai agregator: menggabungkan
// interface method utama ({{ModuleTitle}}Repository/{{ModuleTitle}}Service, yang
// didefinisikan di {{ModuleName}}_interfaces.go) dengan interface item tambahan
// hasil mode add-item (mis. TagRepository/TagService) via marker GEN:.
var tmplContracts = `package contracts

// Repository defines database operations.
// Method utama {{.ModuleTitle}} didefinisikan di {{.ModuleName}}_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	{{.ModuleTitle}}Repository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama {{.ModuleTitle}} didefinisikan di {{.ModuleName}}_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	{{.ModuleTitle}}Service
	// GEN:ITEM_SERVICE_INTERFACE
}
`

// tmplModuleInterfaces mendefinisikan interface method utama entitas
// {{.ModuleTitle}}, dengan method diberi suffix {{.MethodSuffix}} — pola yang
// SAMA dengan interface item hasil mode add-item (mis. TagRepository/TagService),
// supaya konsisten di seluruh sub-module.
var tmplModuleInterfaces = `package contracts

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// {{.ModuleTitle}}Repository defines database operations for {{.ModuleTitle}}.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/{{.ModuleName}}_repository.go).
type {{.ModuleTitle}}Repository interface {
	Create{{.MethodSuffix}}(m *models.{{.ModuleTitle}}) error
	Get{{.MethodSuffix}}ByID(id int64) (*models.{{.ModuleTitle}}, error)
	List{{.MethodSuffix}}(page, pageSize int, filter *dto.Filter{{.ModuleTitle}}Request) ([]models.{{.ModuleTitle}}, int64, error)
	Update{{.MethodSuffix}}(m *models.{{.ModuleTitle}}) error
	Delete{{.MethodSuffix}}(id int64) error
}

// {{.ModuleTitle}}Service defines business logic operations for {{.ModuleTitle}}.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/{{.ModuleName}}_service.go).
type {{.ModuleTitle}}Service interface {
	Create{{.MethodSuffix}}(req *dto.Create{{.ModuleTitle}}Request, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error)
	Get{{.MethodSuffix}}ByID(id int64, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error)
	List{{.MethodSuffix}}(page, pageSize int, filter *dto.Filter{{.ModuleTitle}}Request, actor he.AuthContext) ([]dto.{{.ModuleTitle}}Response, int64, error)
	Update{{.MethodSuffix}}(id int64, req *dto.Update{{.ModuleTitle}}Request, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error)
	Delete{{.MethodSuffix}}(id int64, actor he.AuthContext) error
}
`

var tmplRequest = `package dto

// Create{{.ModuleTitle}}Request request body untuk membuat {{.ModuleTitle}} baru
type Create{{.ModuleTitle}}Request struct {
	Name        string  ` + "`" + `json:"name" validate:"required,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}

// Update{{.ModuleTitle}}Request request body untuk update {{.ModuleTitle}}
type Update{{.ModuleTitle}}Request struct {
	Name        *string ` + "`" + `json:"name" validate:"omitempty,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}

// Filter{{.ModuleTitle}}Request request body untuk filter {{.ModuleTitle}}
type Filter{{.ModuleTitle}}Request struct {
	Name string ` + "`" + `query:"name"` + "`" + `
}
`

var tmplResponse = `package dto

import (

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/shared/types"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// {{.ModuleTitle}}Response response untuk single {{.ModuleTitle}}
type {{.ModuleTitle}}Response struct {
	ID          int64     ` + "`" + `json:"id"` + "`" + `
	Name        string    ` + "`" + `json:"name"` + "`" + `
	Description *string   ` + "`" + `json:"description"` + "`" + `
	CreatedBy   *he.UserData ` + "`" + `json:"created_by"` + "`" + `
	UpdatedBy   *he.UserData ` + "`" + `json:"updated_by"` + "`" + `
	CreatedAt   types.CustomTime ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   types.CustomTime ` + "`" + `json:"updated_at"` + "`" + `
}

type {{.ModuleTitle}}ResponseParams struct {
	{{.ModuleTitle}} *models.{{.ModuleTitle}}
	Creator         *he.UserData
	Updater         *he.UserData
}

// To{{.ModuleTitle}}Response mengubah model menjadi response
func To{{.ModuleTitle}}Response(params {{.ModuleTitle}}ResponseParams) *{{.ModuleTitle}}Response {
	return &{{.ModuleTitle}}Response{
		ID:          params.{{.ModuleTitle}}.ID,
		Name:        params.{{.ModuleTitle}}.Name,
		Description: params.{{.ModuleTitle}}.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.{{.ModuleTitle}}.CreatedAt),
		UpdatedAt:   types.CustomTime(params.{{.ModuleTitle}}.UpdatedAt),
	}
}

// To{{.ModuleTitle}}ListResponse mengubah slice model menjadi slice response
func To{{.ModuleTitle}}ListResponse(
	items []models.{{.ModuleTitle}},
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []{{.ModuleTitle}}Response {
	responses := make([]{{.ModuleTitle}}Response, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *To{{.ModuleTitle}}Response({{.ModuleTitle}}ResponseParams{
			{{.ModuleTitle}}: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
`

var tmplModel = `package models

import (
	"time"

	"gorm.io/gorm"
)

// {{.ModuleTitle}} represents the {{.TableName}} table in database
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
	return "{{.TableName}}"
}
`

var tmplMainRepository = `package repositories

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
	"gorm.io/gorm"
)

// repository adalah satu-satunya struct repository untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat struct repository baru — method
// CRUD-nya ditempelkan langsung ke struct ini di file terpisah
// (mis. repositories/tag_repository.go), sehingga satu instance struct ini
// otomatis memenuhi contracts.Repository maupun interface item (TagRepository, dst).
type repository struct {
	db *gorm.DB
}

// New{{.ModuleTitle}}Repository membuat instance repository baru
func New{{.ModuleTitle}}Repository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
`

var tmplRepository = `package repositories

import (
	"errors"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) Create{{.MethodSuffix}}(m *models.{{.ModuleTitle}}) error {
	return r.db.Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) Get{{.MethodSuffix}}ByID(id int64) (*models.{{.ModuleTitle}}, error) {
	var m models.{{.ModuleTitle}}
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) List{{.MethodSuffix}}(page, pageSize int, filter *dto.Filter{{.ModuleTitle}}Request) ([]models.{{.ModuleTitle}}, int64, error) {
	var items []models.{{.ModuleTitle}}
	var total int64

	query := r.db.Model(&models.{{.ModuleTitle}}{}).Where("deleted_at IS NULL")

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) Update{{.MethodSuffix}}(m *models.{{.ModuleTitle}}) error {
	return r.db.Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) Delete{{.MethodSuffix}}(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.{{.ModuleTitle}}{}).Error
}
`

var tmplMainService = `package services

import (
	"{{.ProjectModule}}/config"
	{{.ModuleName}}Contracts "{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
	

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	authContracts "{{.ProjectModule}}/internal/modules/auth/contracts"
	rbacContracts "{{.ProjectModule}}/internal/modules/rbac/contracts"
	userContracts "{{.ProjectModule}}/internal/modules/users/contracts"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// service adalah satu-satunya struct service untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat struct service baru — method
// CRUD & permission-nya ditempelkan langsung ke struct ini (mis.
// services/tag_service.go, services/tag_permission.go), dan repo field
// di bawah ini otomatis mencakup method item begitu contracts.Repository
// di-embed dengan interface repository item (lihat contracts/interfaces.go).
type service struct {
	repo     {{.ModuleName}}Contracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	userRepo userContracts.Repository
	cfg      *config.Config
}

// New{{.ModuleTitle}}Service membuat instance service baru
func New{{.ModuleTitle}}Service(
	repo {{.ModuleName}}Contracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	cfg    *config.Config,
) {{.ModuleName}}Contracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// buildCreator mengambil data creator user
func (s *service) buildCreator(createdBy *int64) *he.UserData {
	if createdBy == nil {
		return nil
	}
	creator, err := s.userRepo.GetByID(*createdBy)
	if err != nil || creator == nil {
		return nil
	}
	return &he.UserData{
		ID:       creator.ID,
		Username: creator.Username,
		Name:     creator.Name,
	}
}

// ── helper: build creator/updater maps ───────────────────────────────────────
func (s *service) buildAuditMaps(items []models.{{.ModuleTitle}}) (map[int64]*he.UserData, map[int64]*he.UserData) {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			idSet[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			idSet[*item.UpdatedBy] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users, err := s.userRepo.GetByIDs(ids) // ← 1 query total, bukan 40
	if err != nil {
		return map[int64]*he.UserData{}, map[int64]*he.UserData{}
	}

	userMap := make(map[int64]*he.UserData, len(users))
	for _, u := range users {
		userMap[u.ID] = &he.UserData{ID: u.ID, Username: u.Username, Name: u.Name}
	}
	// creator dan updater sekarang share map yang sama — reuse otomatis, kode lebih pendek juga
	return userMap, userMap
}

`

var tmplPermissionService = `package services

import (
	rbacMiddlewares "{{.ProjectModule}}/internal/modules/rbac/middlewares"
	rbacModels "{{.ProjectModule}}/internal/modules/rbac/models"
	he "{{.ProjectModule}}/internal/shared/httputil"
)


// ── canRead ───────────────────────────────────────────────────────────────────
func (s *service) canRead{{.ModuleTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyRead); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canCreate ─────────────────────────────────────────────────────────────────
func (s *service) canCreate{{.ModuleTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canUpdate ─────────────────────────────────────────────────────────────────
func (s *service) canUpdate{{.ModuleTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canDelete ─────────────────────────────────────────────────────────────────
func (s *service) canDelete{{.ModuleTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}
`

var tmplService = `package services

import (
	"errors"
	"net/http"
	"time"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	appErrors "{{.ProjectModule}}/internal/shared/errors"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) Create{{.MethodSuffix}}(req *dto.Create{{.ModuleTitle}}Request, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error) {
	can, err := s.canCreate{{.ModuleTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat {{.ModuleTitle}} baru.", nil)
	}

	m := &models.{{.ModuleTitle}}{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.Create{{.MethodSuffix}}(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)

	return dto.To{{.ModuleTitle}}Response(dto.{{.ModuleTitle}}ResponseParams{
		{{.ModuleTitle}}: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) Get{{.MethodSuffix}}ByID(id int64, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error) {
	can, err := s.canRead{{.ModuleTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat {{.ModuleTitle}}.", nil)
	}

	m, err := s.repo.Get{{.MethodSuffix}}ByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ModuleTitle}} tidak ditemukan")
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.To{{.ModuleTitle}}Response(dto.{{.ModuleTitle}}ResponseParams{
		{{.ModuleTitle}}: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) List{{.MethodSuffix}}(page, pageSize int, filter *dto.Filter{{.ModuleTitle}}Request, actor he.AuthContext) ([]dto.{{.ModuleTitle}}Response, int64, error) {
	can, err := s.canRead{{.ModuleTitle}}(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar {{.ModuleTitle}}.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.List{{.MethodSuffix}}(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.To{{.ModuleTitle}}ListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) Update{{.MethodSuffix}}(id int64, req *dto.Update{{.ModuleTitle}}Request, actor he.AuthContext) (*dto.{{.ModuleTitle}}Response, error) {
	can, err := s.canUpdate{{.ModuleTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah {{.ModuleTitle}}.", nil)
	}

	m, err := s.repo.Get{{.MethodSuffix}}ByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ModuleTitle}} tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.Update{{.MethodSuffix}}(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.To{{.ModuleTitle}}Response(dto.{{.ModuleTitle}}ResponseParams{
		{{.ModuleTitle}}: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) Delete{{.MethodSuffix}}(id int64, actor he.AuthContext) error {
	can, err := s.canDelete{{.ModuleTitle}}(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus {{.ModuleTitle}}.", nil)
	}

	m, err := s.repo.Get{{.MethodSuffix}}ByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("{{.ModuleTitle}} tidak ditemukan")
	}
	return s.repo.Delete{{.MethodSuffix}}(id)
}
`

var tmplMainHandler = `package handlers

import (
	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
)

// {{.ModuleTitle}}Handler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (List{{.MethodSuffix}}, Create{{.MethodSuffix}}, dst).
type {{.ModuleTitle}}Handler struct {
	service contracts.Service
	cfg     *config.Config
}

// New{{.ModuleTitle}}Handler membuat instance handler baru
func New{{.ModuleTitle}}Handler(service contracts.Service, cfg *config.Config) *{{.ModuleTitle}}Handler {
	return &{{.ModuleTitle}}Handler{service: service, cfg: cfg}
}
`

var tmplHandler = `package handlers

import (
	"io"
	"net/http"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/shared/response"
	"{{.ProjectModule}}/internal/shared/validator"
	"{{.ProjectModule}}/internal/shared/binding"
	he "{{.ProjectModule}}/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── List{{.MethodSuffix}} ─────────────────────────────────────────────────────
//
//	@Summary		Get list of {{.ModuleTitle}}
//	@Description	Get paginated list of {{.ModuleTitle}}
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.{{.ModuleTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}} [get]
func (h *{{.ModuleTitle}}Handler) List{{.MethodSuffix}}(c *echo.Context) error {

	filter := dto.Filter{{.ModuleTitle}}Request{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.List{{.MethodSuffix}}(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── Get{{.MethodSuffix}}ByID ───────────────────────────────────────────────────
//
//	@Summary		Get {{.ModuleTitle}}
//	@Description	Get {{.ModuleTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"{{.ModuleTitle}} ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.{{.ModuleTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [get]
func (h *{{.ModuleTitle}}Handler) Get{{.MethodSuffix}}ByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Get{{.MethodSuffix}}ByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── Create{{.MethodSuffix}} ────────────────────────────────────────────────────
//
//	@Summary		Create {{.ModuleTitle}}
//	@Description	Create New {{.ModuleTitle}}
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.Create{{.ModuleTitle}}Request	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.{{.ModuleTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}} [post]
func (h *{{.ModuleTitle}}Handler) Create{{.MethodSuffix}}(c *echo.Context) error {
	var req dto.Create{{.ModuleTitle}}Request
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Create{{.MethodSuffix}}(&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── Update{{.MethodSuffix}} ────────────────────────────────────────────────────
//
//	@Summary		Update {{.ModuleTitle}}
//	@Description	Update {{.ModuleTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"{{.ModuleTitle}} ID"
//	@Param			body	body		dto.Update{{.ModuleTitle}}Request	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.{{.ModuleTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [put]
func (h *{{.ModuleTitle}}Handler) Update{{.MethodSuffix}}(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.Update{{.ModuleTitle}}Request
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Update{{.MethodSuffix}}(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "{{.ModuleTitle}} tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete{{.MethodSuffix}} ────────────────────────────────────────────────────
//
//	@Summary		Delete {{.ModuleTitle}}
//	@Description	Delete {{.ModuleTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"{{.ModuleTitle}} ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [delete]
func (h *{{.ModuleTitle}}Handler) Delete{{.MethodSuffix}}(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.Delete{{.MethodSuffix}}(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "{{.ModuleTitle}} tidak ditemukan" {
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

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"

	"gorm.io/gorm"
)

//go:embed 001_create_{{.TableName}}_table.sql
var {{.PackageName}}SQL string

// Migrate{{.ModuleTitle}} menjalankan GORM auto-migration
func Migrate{{.ModuleTitle}}(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.{{.ModuleTitle}}{})
}

// Migrate{{.ModuleTitle}}WithSQL menjalankan migrasi via raw SQL
func Migrate{{.ModuleTitle}}WithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec({{.PackageName}}SQL)
	if err != nil {
		log.Printf("Error creating {{.TableName}} table: %v", err)
		return err
	}
	log.Println("{{.TableName}} table migrated successfully")
	return nil
}

// Drop{{.ModuleTitle}}Table menghapus tabel (gunakan dengan hati-hati!)
func Drop{{.ModuleTitle}}Table(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.{{.ModuleTitle}}{})
}
`

var tmplSQL = `-- Migration: Create {{.TableName}} table
-- Timestamp: {{.Timestamp}}

CREATE TABLE IF NOT EXISTS {{.TableName}} (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_{{.TableName}}_deleted_at ON {{.TableName}}(deleted_at);
`

var tmplFactory = `package factories

import (
	"fmt"
	"math/rand"
	"time"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// {{.ModuleTitle}}Factory membuat data {{.ModuleTitle}} untuk testing/seeding
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
	desc := fmt.Sprintf("Deskripsi {{.ModuleTitle}} %d", idx)

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

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/tests/factories"

	"gorm.io/gorm"
)

// {{.ModuleTitle}}Seeder mengelola seeding data {{.ModuleTitle}}
type {{.ModuleTitle}}Seeder struct {
	db *gorm.DB
}

func New{{.ModuleTitle}}Seeder(db *gorm.DB) *{{.ModuleTitle}}Seeder {
	return &{{.ModuleTitle}}Seeder{db: db}
}

func (s *{{.ModuleTitle}}Seeder) Run() error {
	log.Println("🌱 Seeding {{.TableName}}...")

	items := factories.New{{.ModuleTitle}}Factory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat {{.ModuleTitle}}: %v", err)
			continue
		}
		log.Printf("   ✅ {{.ModuleTitle}} '%s' dibuat.", item.Name)
	}

	log.Println("✅ {{.TableName}} seeding selesai!")
	return nil
}

func (s *{{.ModuleTitle}}Seeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data {{.TableName}}...")
	if err := s.db.Exec("DELETE FROM {{.TableName}}").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE {{.TableName}}_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

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
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"

	"gorm.io/gorm"
)

// SetupTestDB membuat koneksi DB untuk keperluan test
func SetupTestDB() *gorm.DB {
	cfg := config.LoadConfig()
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi DB untuk test:", err)
	}
	return db
}

func MigrateTestDB(db *gorm.DB) {
	if err := db.AutoMigrate(&models.{{.ModuleTitle}}{}); err != nil {
		log.Fatal("Gagal migrasi test DB:", err)
	}
}

func TruncateTable(db *gorm.DB, tables ...string) {
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Printf("Warning: Gagal truncate table %s: %v", table, err)
		}
	}
}
`

var tmplModule = `package {{.PackageName}}

import (
	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/handlers"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/repositories"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/services"
	"{{.ProjectModule}}/internal/shared/utils"

	authContracts "{{.ProjectModule}}/internal/modules/auth/contracts"
	rbacContracts "{{.ProjectModule}}/internal/modules/rbac/contracts"
	userContracts "{{.ProjectModule}}/internal/modules/users/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.{{.ModuleTitle}}Handler
	jwtManager *utils.JWTManager
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository
}

func NewModule(
	db *gorm.DB,
	jwtManager *utils.JWTManager,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	cfg *config.Config,
) *Module {
	repo := repositories.New{{.ModuleTitle}}Repository(db)
	svc := services.New{{.ModuleTitle}}Service(repo, rbacRepo, authRepo,userRepo, cfg)
	handler := handlers.New{{.ModuleTitle}}Handler(svc, cfg)

	return &Module{
		db:         db,
		handler:    handler,
		jwtManager: jwtManager,
		repo:       repo,
		rbacRepo:   rbacRepo,
	}
}

func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager, m.db)
}
`

var tmplRoutes = `package {{.PackageName}}

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/handlers"
	authMiddlewares "{{.ProjectModule}}/internal/modules/auth/middlewares"
	"{{.ProjectModule}}/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.{{.ModuleTitle}}Handler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("{{.URLPrefix}}", jwt)
	g.GET("", h.List{{.MethodSuffix}})
	g.GET("/:id", h.Get{{.MethodSuffix}}ByID)
	g.POST("", h.Create{{.MethodSuffix}})
	g.PUT("/:id", h.Update{{.MethodSuffix}})
	g.DELETE("/:id", h.Delete{{.MethodSuffix}})
	// GEN:ITEM_ROUTES
}
`

var tmplRegister = `package {{.PackageName}}

import (
	"database/sql"

	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/apps"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/migrations"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/shared/utils"

	authContracts "{{.ProjectModule}}/internal/modules/auth/contracts"
	authRepositories "{{.ProjectModule}}/internal/modules/auth/repositories"
	rbacContracts "{{.ProjectModule}}/internal/modules/rbac/contracts"
	rbacRepositories "{{.ProjectModule}}/internal/modules/rbac/repositories"
	userRepositories "neosim_go/internal/modules/users/repositories"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db       *gorm.DB
	cfg      *config.Config
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

func init() {
	apps.Register(&registryModule{})
}

func (r *registryModule) SetDB(db *gorm.DB) {
	r.db = db
	r.rbacRepo = rbacRepositories.NewRBACRepository(db)
	r.authRepo = authRepositories.NewAuthRepository(db)
}

func (r *registryModule) SetConfig(cfg *config.Config) {
	r.cfg = cfg
}

func (r *registryModule) InitRoutes(e *echo.Echo) {
	jwtManager := utils.NewJWTManager(
		r.cfg.JWTSecret,
		r.cfg.JWTIssuer,
		r.cfg.JWTAccessTokenExpMinutes,
		r.cfg.JWTRefreshTokenExpDays,
	)
	userRepo := userRepositories.NewRepository(r.db)
	NewModule(r.db, jwtManager, r.rbacRepo, r.authRepo,userRepo, r.cfg).InitRoutes(e)
}

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.{{.ModuleTitle}}{},
		// GEN:ITEM_MODELS
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return nil
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	if err := migrations.Migrate{{.ModuleTitle}}WithSQL(sqlDB); err != nil {
		return err
	}
	// GEN:ITEM_MIGRATIONS
	return nil
}
`

var tmplModuleServiceMock = `package mocks

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"github.com/stretchr/testify/mock"
)

// {{.ModuleTitle}}RepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type {{.ModuleTitle}}RepositoryMock struct {
	mock.Mock
}

func (m *{{.ModuleTitle}}RepositoryMock) Create{{.MethodSuffix}}(item *models.{{.ModuleTitle}}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *{{.ModuleTitle}}RepositoryMock) Get{{.MethodSuffix}}ByID(id int64) (*models.{{.ModuleTitle}}, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.{{.ModuleTitle}}), args.Error(1)
}

func (m *{{.ModuleTitle}}RepositoryMock) List{{.MethodSuffix}}(page, pageSize int, filter *dto.Filter{{.ModuleTitle}}Request) ([]models.{{.ModuleTitle}}, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.{{.ModuleTitle}}), args.Get(1).(int64), args.Error(2)
}

func (m *{{.ModuleTitle}}RepositoryMock) Update{{.MethodSuffix}}(item *models.{{.ModuleTitle}}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *{{.ModuleTitle}}RepositoryMock) Delete{{.MethodSuffix}}(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
`

var tmplRBACMock = `package mocks

import (
	rbacModels "{{.ProjectModule}}/internal/modules/rbac/models"
	"github.com/stretchr/testify/mock"
)

// RBACRepositoryMock is a mock implementation of rbacContracts.RBACRepository
type RBACRepositoryMock struct {
	mock.Mock
}

func (m *RBACRepositoryMock) IsSuperadmin(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *RBACRepositoryMock) CreatePermission(p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetPermissionByID(id int64) (*rbacModels.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) GetPermissionByName(name string) (*rbacModels.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) ListPermissions(page, pageSize int) ([]rbacModels.Permission, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdatePermission(p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeletePermission(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) CreateRole(r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRoleByID(id int64) (*rbacModels.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) GetRoleByName(name string) (*rbacModels.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) ListRoles(page, pageSize int) ([]rbacModels.Role, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Role), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdateRole(r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeleteRole(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUsersRoles(userIDs []int64) (map[int64][]rbacModels.Role, error) {
	args := m.Called(userIDs)
	return args.Get(0).(map[int64][]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) AssignPermissionsToRole(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokePermissionsFromRole(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRolePermissions(roleID int64) ([]rbacModels.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) SyncRolePermissions(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignRolesToUser(userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeRolesFromUser(userID int64, roleIDs []int64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserRoles(userID int64) ([]rbacModels.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) SyncUserRoles(userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignDirectPermission(userID, permissionID int64, isGranted bool, assignedBy *int64) error {
	args := m.Called(userID, permissionID, isGranted, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeDirectPermission(userID, permissionID int64) error {
	args := m.Called(userID, permissionID)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserDirectPermissions(userID int64) ([]rbacModels.UserPermission, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.UserPermission), args.Error(1)
}

func (m *RBACRepositoryMock) GetUserAllPermissions(userID int64) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *RBACRepositoryMock) HasPermission(userID int64, permission string) (bool, error) {
	args := m.Called(userID, permission)
	return args.Bool(0), args.Error(1)
}
`

var tmplAuthMock = `package mocks

import (
	authModels "{{.ProjectModule}}/internal/modules/auth/models"
	"github.com/stretchr/testify/mock"
)

// AuthRepositoryMock is a mock implementation of authContracts.AuthRepository
type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) SaveToken(token *authModels.AuthToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetTokenByJTI(jti string) (*authModels.AuthToken, error) {
	args := m.Called(jti)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModels.AuthToken), args.Error(1)
}

func (m *AuthRepositoryMock) BlacklistToken(jti string) error {
	args := m.Called(jti)
	return args.Error(0)
}

func (m *AuthRepositoryMock) BlacklistAllUserTokens(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *AuthRepositoryMock) CountActiveTokens(userID int64) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AuthRepositoryMock) SaveLoginHistory(history *authModels.LoginHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetUserLoginHistories(userID int64, limit int) ([]authModels.LoginHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]authModels.LoginHistory), args.Error(1)
}

func (m *AuthRepositoryMock) SavePasswordHistory(history *authModels.PasswordHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetPasswordHistories(userID int64, limit int) ([]authModels.PasswordHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]authModels.PasswordHistory), args.Error(1)
}
`

var tmplUserMock = `package mocks

import (
	userDto    "{{.ProjectModule}}/internal/modules/users/dto"
	userModels "{{.ProjectModule}}/internal/modules/users/models"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock adalah mock dari userContracts.Repository (modul users).
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *userModels.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(id int64) (*userModels.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByUsername(username string) (*userModels.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*userModels.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) List(page, pageSize int, filter *userDto.UserFilter) ([]userModels.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	var users []userModels.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userModels.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *UserRepositoryMock) Update(user *userModels.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(id int64, deletedBy int64, reason string) error {
	args := m.Called(id, deletedBy, reason)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeletedList(page, pageSize int, filter *userDto.UserDeletedFilter) ([]userModels.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	var users []userModels.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userModels.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *UserRepositoryMock) GetSettings(id int64) ([]userModels.UserSetting, error) {
	args := m.Called(id)
	var settings []userModels.UserSetting
	if args.Get(0) != nil {
		settings = args.Get(0).([]userModels.UserSetting)
	}
	return settings, args.Error(1)
}

func (m *UserRepositoryMock) UpdateSettings(id int64, settings []userModels.UserSetting) error {
	args := m.Called(id, settings)
	return args.Error(0)
}
`

var tmplModuleServiceTest = `package tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"{{.ProjectModule}}/config"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/services"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/tests/factories"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/tests/mocks"

	{{.ModuleName}}Contracts "{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
	rbacModels "{{.ProjectModule}}/internal/modules/rbac/models"
	appErrors "{{.ProjectModule}}/internal/shared/errors"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  {{.ModuleTitle}} Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  {{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  {{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}")
	}

	os.Exit(code)
}

// {{.ModuleTitle}}ServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type {{.ModuleTitle}}ServiceTestSuite struct {
	suite.Suite
	repo     *mocks.{{.ModuleTitle}}RepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      {{.ModuleName}}Contracts.Service
	cfg      *config.Config
}

func (s *{{.ModuleTitle}}ServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.{{.ModuleTitle}}RepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg      = &config.Config{}
	s.svc = services.New{{.ModuleTitle}}Service(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
}

func Test{{.ModuleTitle}}Service(t *testing.T) {
	suite.Run(t, new({{.ModuleTitle}}ServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *{{.ModuleTitle}}ServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *{{.ModuleTitle}}ServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Create{{.MethodSuffix}}_Superadmin_Success() {
	req := &dto.Create{{.ModuleTitle}}Request{Name: "Test {{.ModuleTitle}}"}
	actor := superadminActor()

	s.repo.On("Create{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(nil)

	result, err := s.svc.Create{{.MethodSuffix}}(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Create{{.MethodSuffix}}_WithPermission_Success() {
	req := &dto.Create{{.ModuleTitle}}Request{Name: "Test {{.ModuleTitle}}"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("Create{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(nil)

	result, err := s.svc.Create{{.MethodSuffix}}(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Create{{.MethodSuffix}}_WithManagePermission_Success() {
	req := &dto.Create{{.ModuleTitle}}Request{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("Create{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(nil)

	result, err := s.svc.Create{{.MethodSuffix}}(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Create{{.MethodSuffix}}_Forbidden() {
	req := &dto.Create{{.ModuleTitle}}Request{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.Create{{.MethodSuffix}}(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Create{{.MethodSuffix}}_RepoError() {
	req := &dto.Create{{.ModuleTitle}}Request{Name: "Test"}
	actor := superadminActor()

	s.repo.On("Create{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Create{{.MethodSuffix}}(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Get{{.MethodSuffix}}ByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.New{{.ModuleTitle}}Factory().Make()
	item.ID = 1

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(item, nil)

	result, err := s.svc.Get{{.MethodSuffix}}ByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Get{{.MethodSuffix}}ByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.New{{.ModuleTitle}}Factory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(item, nil)

	result, err := s.svc.Get{{.MethodSuffix}}ByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Get{{.MethodSuffix}}ByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.Get{{.MethodSuffix}}ByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Get{{.MethodSuffix}}ByID_NotFound() {
	actor := superadminActor()

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Get{{.MethodSuffix}}ByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Get{{.MethodSuffix}}ByID_RepoError() {
	actor := superadminActor()

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.Get{{.MethodSuffix}}ByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{}
	items := []models.{{.ModuleTitle}}{
		*factories.New{{.ModuleTitle}}Factory().Make(),
		*factories.New{{.ModuleTitle}}Factory().Make(),
	}

	s.repo.On("List{{.MethodSuffix}}", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.List{{.MethodSuffix}}(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{}
	items := []models.{{.ModuleTitle}}{*factories.New{{.ModuleTitle}}Factory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("List{{.MethodSuffix}}", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.List{{.MethodSuffix}}(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_Forbidden() {
	actor := regularActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{}
	s.mockNoPermissions()

	result, total, err := s.svc.List{{.MethodSuffix}}(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{}

	s.repo.On("List{{.MethodSuffix}}", 1, 10, filter).Return([]models.{{.ModuleTitle}}{}, int64(0), nil)

	result, total, err := s.svc.List{{.MethodSuffix}}(0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{}

	s.repo.On("List{{.MethodSuffix}}", 1, 10, filter).Return([]models.{{.ModuleTitle}}{}, int64(0), nil)

	_, _, err := s.svc.List{{.MethodSuffix}}(1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "List{{.MethodSuffix}}", 1, 10, filter)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_List{{.MethodSuffix}}_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.Filter{{.ModuleTitle}}Request{Name: "test"}
	items := []models.{{.ModuleTitle}}{*factories.New{{.ModuleTitle}}Factory().Make()}

	s.repo.On("List{{.MethodSuffix}}", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.List{{.MethodSuffix}}(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.Update{{.ModuleTitle}}Request{Name: &newName}

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Update{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(nil)

	result, err := s.svc.Update{{.MethodSuffix}}(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_WithPermission_Success() {
	actor := regularActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.Update{{.ModuleTitle}}Request{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Update{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(nil)

	result, err := s.svc.Update{{.MethodSuffix}}(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_Forbidden() {
	actor := regularActor()
	req := &dto.Update{{.ModuleTitle}}Request{}
	s.mockNoPermissions()

	result, err := s.svc.Update{{.MethodSuffix}}(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_NotFound() {
	actor := superadminActor()
	req := &dto.Update{{.ModuleTitle}}Request{}

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Update{{.MethodSuffix}}(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_PartialFields() {
	actor := superadminActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.Update{{.ModuleTitle}}Request{Description: &newDesc}

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Update{{.MethodSuffix}}", mock.MatchedBy(func(m *models.{{.ModuleTitle}}) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.Update{{.MethodSuffix}}(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Update{{.MethodSuffix}}_RepoError() {
	actor := superadminActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1
	req := &dto.Update{{.ModuleTitle}}Request{}

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Update{{.MethodSuffix}}", mock.AnythingOfType("*models.{{.ModuleTitle}}")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Update{{.MethodSuffix}}(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Delete{{.MethodSuffix}}_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete{{.MethodSuffix}}", int64(1)).Return(nil)

	err := s.svc.Delete{{.MethodSuffix}}(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Delete{{.MethodSuffix}}_WithPermission_Success() {
	actor := regularActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete{{.MethodSuffix}}", int64(1)).Return(nil)

	err := s.svc.Delete{{.MethodSuffix}}(1, actor)

	s.NoError(err)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Delete{{.MethodSuffix}}_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.Delete{{.MethodSuffix}}(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Delete{{.MethodSuffix}}_NotFound() {
	actor := superadminActor()

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(999)).Return(nil, nil)

	err := s.svc.Delete{{.MethodSuffix}}(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.ModuleTitle}}ServiceTestSuite) Test_Delete{{.MethodSuffix}}_RepoError() {
	actor := superadminActor()
	existing := factories.New{{.ModuleTitle}}Factory().Make()
	existing.ID = 1

	s.repo.On("Get{{.MethodSuffix}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete{{.MethodSuffix}}", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.Delete{{.MethodSuffix}}(1, actor)

	s.Error(err)
}
`

// ═══════════════════════════════════════════════════════════════════════════════
// TEMPLATES — MODE 3 (Add Item, SHARED STRUCT)
// ═══════════════════════════════════════════════════════════════════════════════

var tmplItemContracts = `package contracts

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// {{.ItemTitle}}Repository defines database operations for {{.ItemTitle}}.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type {{.ItemTitle}}Repository interface {
	Create{{.ItemTitle}}(m *models.{{.ItemTitle}}) error
	Get{{.ItemTitle}}ByID(id int64) (*models.{{.ItemTitle}}, error)
	List{{.ItemTitle}}(page, pageSize int, filter *dto.Filter{{.ItemTitle}}Request) ([]models.{{.ItemTitle}}, int64, error)
	Update{{.ItemTitle}}(m *models.{{.ItemTitle}}) error
	Delete{{.ItemTitle}}(id int64) error
}

// {{.ItemTitle}}Service defines business logic operations for {{.ItemTitle}}.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type {{.ItemTitle}}Service interface {
	Create{{.ItemTitle}}(req *dto.Create{{.ItemTitle}}Request, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error)
	Get{{.ItemTitle}}ByID(id int64, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error)
	List{{.ItemTitle}}(page, pageSize int, filter *dto.Filter{{.ItemTitle}}Request, actor he.AuthContext) ([]dto.{{.ItemTitle}}Response, int64, error)
	Update{{.ItemTitle}}(id int64, req *dto.Update{{.ItemTitle}}Request, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error)
	Delete{{.ItemTitle}}(id int64, actor he.AuthContext) error
}
`

var tmplItemRequest = `package dto

// Create{{.ItemTitle}}Request request body untuk membuat {{.ItemTitle}} baru
type Create{{.ItemTitle}}Request struct {
	Name        string  ` + "`" + `json:"name" validate:"required,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}

// Update{{.ItemTitle}}Request request body untuk update {{.ItemTitle}}
type Update{{.ItemTitle}}Request struct {
	Name        *string ` + "`" + `json:"name" validate:"omitempty,min=1,max=255"` + "`" + `
	Description *string ` + "`" + `json:"description" validate:"omitempty,max=500"` + "`" + `
}

// Filter{{.ItemTitle}}Request request body untuk filter {{.ItemTitle}}
type Filter{{.ItemTitle}}Request struct {
	Name string ` + "`" + `query:"name"` + "`" + `
}
`

var tmplItemResponse = `package dto

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	he "{{.ProjectModule}}/internal/shared/httputil"
	"{{.ProjectModule}}/internal/shared/types"
)

// {{.ItemTitle}}Response response untuk single {{.ItemTitle}}
type {{.ItemTitle}}Response struct {
	ID          int64            ` + "`" + `json:"id"` + "`" + `
	Name        string           ` + "`" + `json:"name"` + "`" + `
	Description *string          ` + "`" + `json:"description"` + "`" + `
	CreatedBy   *he.UserData     ` + "`" + `json:"created_by"` + "`" + `
	UpdatedBy   *he.UserData     ` + "`" + `json:"updated_by"` + "`" + `
	CreatedAt   types.CustomTime ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt   types.CustomTime ` + "`" + `json:"updated_at"` + "`" + `
}

type {{.ItemTitle}}ResponseParams struct {
	{{.ItemTitle}} *models.{{.ItemTitle}}
	Creator       *he.UserData
	Updater       *he.UserData
}

// To{{.ItemTitle}}Response mengubah model menjadi response
func To{{.ItemTitle}}Response(params {{.ItemTitle}}ResponseParams) *{{.ItemTitle}}Response {
	return &{{.ItemTitle}}Response{
		ID:          params.{{.ItemTitle}}.ID,
		Name:        params.{{.ItemTitle}}.Name,
		Description: params.{{.ItemTitle}}.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.{{.ItemTitle}}.CreatedAt),
		UpdatedAt:   types.CustomTime(params.{{.ItemTitle}}.UpdatedAt),
	}
}

// To{{.ItemTitle}}ListResponse mengubah slice model menjadi slice response
func To{{.ItemTitle}}ListResponse(
	items []models.{{.ItemTitle}},
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []{{.ItemTitle}}Response {
	responses := make([]{{.ItemTitle}}Response, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *To{{.ItemTitle}}Response({{.ItemTitle}}ResponseParams{
			{{.ItemTitle}}: &m,
			Creator:         creator,
			Updater:         updater,
		}))
	}

	return responses
}
`

var tmplItemModel = `package models

import (
	"time"

	"gorm.io/gorm"
)

// {{.ItemTitle}} represents the {{.TableName}} table in database
type {{.ItemTitle}} struct {
	ID          int64          ` + "`" + `gorm:"primaryKey;autoIncrement;column:id" json:"id"` + "`" + `
	Name        string         ` + "`" + `gorm:"column:name;type:varchar(255);not null" json:"name"` + "`" + `
	Description *string        ` + "`" + `gorm:"column:description;type:text" json:"description"` + "`" + `
	CreatedBy   *int64         ` + "`" + `gorm:"column:created_by" json:"created_by"` + "`" + `
	UpdatedBy   *int64         ` + "`" + `gorm:"column:updated_by" json:"updated_by"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"` + "`" + `
	DeletedAt   gorm.DeletedAt ` + "`" + `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"` + "`" + `
}

func ({{.ItemTitle}}) TableName() string {
	return "{{.TableName}}"
}
`

var tmplItemRepository = `package repositories

import (
	"errors"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/contracts"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"

	"gorm.io/gorm"
)

// New{{.ItemTitle}}Repository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.{{.ItemTitle}}Repository.
// Berguna untuk test {{.ItemTitle}} yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat New{{.SubModuleTitle}}Repository(db).
func New{{.ItemTitle}}Repository(db *gorm.DB) contracts.{{.ItemTitle}}Repository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) Create{{.ItemTitle}}(m *models.{{.ItemTitle}}) error {
	return r.db.Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) Get{{.ItemTitle}}ByID(id int64) (*models.{{.ItemTitle}}, error) {
	var m models.{{.ItemTitle}}
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) List{{.ItemTitle}}(page, pageSize int, filter *dto.Filter{{.ItemTitle}}Request) ([]models.{{.ItemTitle}}, int64, error) {
	var items []models.{{.ItemTitle}}
	var total int64

	query := r.db.Model(&models.{{.ItemTitle}}{}).Where("deleted_at IS NULL")
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) Update{{.ItemTitle}}(m *models.{{.ItemTitle}}) error {
	return r.db.Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) Delete{{.ItemTitle}}(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.{{.ItemTitle}}{}).Error
}
`

var tmplItemService = `package services

import (
	"errors"
	"net/http"
	"time"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	appErrors "{{.ProjectModule}}/internal/shared/errors"
	he "{{.ProjectModule}}/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) Create{{.ItemTitle}}(req *dto.Create{{.ItemTitle}}Request, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error) {
	can, err := s.canCreate{{.ItemTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat {{.ItemTitle}} baru.", nil)
	}

	m := &models.{{.ItemTitle}}{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.Create{{.ItemTitle}}(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.To{{.ItemTitle}}Response(dto.{{.ItemTitle}}ResponseParams{
		{{.ItemTitle}}: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) Get{{.ItemTitle}}ByID(id int64, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error) {
	can, err := s.canRead{{.ItemTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat {{.ItemTitle}}.", nil)
	}

	m, err := s.repo.Get{{.ItemTitle}}ByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ItemTitle}} tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.To{{.ItemTitle}}Response(dto.{{.ItemTitle}}ResponseParams{
		{{.ItemTitle}}: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────	
func (s *service) List{{.ItemTitle}}(page, pageSize int, filter *dto.Filter{{.ItemTitle}}Request, actor he.AuthContext) ([]dto.{{.ItemTitle}}Response, int64, error) {
	can, err := s.canRead{{.ItemTitle}}(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar {{.ItemTitle}}.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.List{{.ItemTitle}}(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsFor{{.ItemTitle}}(items)
	return dto.To{{.ItemTitle}}ListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) Update{{.ItemTitle}}(id int64, req *dto.Update{{.ItemTitle}}Request, actor he.AuthContext) (*dto.{{.ItemTitle}}Response, error) {
	can, err := s.canUpdate{{.ItemTitle}}(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah {{.ItemTitle}}.", nil)
	}

	m, err := s.repo.Get{{.ItemTitle}}ByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("{{.ItemTitle}} tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.Update{{.ItemTitle}}(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.To{{.ItemTitle}}Response(dto.{{.ItemTitle}}ResponseParams{
		{{.ItemTitle}}: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) Delete{{.ItemTitle}}(id int64, actor he.AuthContext) error {
	can, err := s.canDelete{{.ItemTitle}}(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus {{.ItemTitle}}.", nil)
	}

	m, err := s.repo.Get{{.ItemTitle}}ByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("{{.ItemTitle}} tidak ditemukan")
	}
	return s.repo.Delete{{.ItemTitle}}(id)
}

// ── helper khusus {{.ItemTitle}} (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsFor{{.ItemTitle}}(items []models.{{.ItemTitle}}) (map[int64]*he.UserData, map[int64]*he.UserData) {
	fetchUser := func(id int64) (*he.UserData, error) {
		user, err := s.userRepo.GetByID(id)
		if err != nil || user == nil {
			return nil, err
		}
		return &he.UserData{ID: user.ID, Username: user.Username, Name: user.Name}, nil
	}

	creatorIDs := make(map[int64]struct{})
	updaterIDs := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			creatorIDs[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			updaterIDs[*item.UpdatedBy] = struct{}{}
		}
	}

	creatorsMap := make(map[int64]*he.UserData)
	for id := range creatorIDs {
		if data, err := fetchUser(id); err == nil && data != nil {
			creatorsMap[id] = data
		}
	}

	updatersMap := make(map[int64]*he.UserData)
	for id := range updaterIDs {
		if data, ok := creatorsMap[id]; ok {
			updatersMap[id] = data
		} else if data, err := fetchUser(id); err == nil && data != nil {
			updatersMap[id] = data
		}
	}

	return creatorsMap, updatersMap
}


`

var tmplItemPermission = `package services

import (
	rbacMiddlewares "{{.ProjectModule}}/internal/modules/rbac/middlewares"
	rbacModels "{{.ProjectModule}}/internal/modules/rbac/models"
	he "{{.ProjectModule}}/internal/shared/httputil"
)


// ── CanRead ───────────────────────────────────────────────────────────────────
func (s *service) canRead{{.ItemTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyRead); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canCreate ─────────────────────────────────────────────────────────────────
func (s *service) canCreate{{.ItemTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canUpdate ─────────────────────────────────────────────────────────────────
func (s *service) canUpdate{{.ItemTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canDelete ─────────────────────────────────────────────────────────────────
func (s *service) canDelete{{.ItemTitle}}(actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}
`

var tmplItemHandler = `package handlers

import (
	"io"
	"net/http"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	he "{{.ProjectModule}}/internal/shared/httputil"
	"{{.ProjectModule}}/internal/shared/response"
	"{{.ProjectModule}}/internal/shared/validator"
	"{{.ProjectModule}}/internal/shared/binding"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct {{.SubModuleTitle}}Handler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix {{.ItemTitle}} agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── List{{.ItemTitle}} ──────────────────────────────────────────────────────
//
//	@Summary		Get list of {{.ItemTitle}}
//	@Description	Get paginated list of {{.ItemTitle}}
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.{{.ItemTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}} [get]
func (h *{{.SubModuleTitle}}Handler) List{{.ItemTitle}}(c *echo.Context) error {
	filter := dto.Filter{{.ItemTitle}}Request{Name: c.QueryParam("name")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.List{{.ItemTitle}}(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── Get{{.ItemTitle}}ByID ───────────────────────────────────────────────────
//
//	@Summary		Get {{.ItemTitle}}
//	@Description	Get {{.ItemTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"{{.ItemTitle}} ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.{{.ItemTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [get]
func (h *{{.SubModuleTitle}}Handler) Get{{.ItemTitle}}ByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Get{{.ItemTitle}}ByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── Create{{.ItemTitle}} ────────────────────────────────────────────────────
//
//	@Summary		Create {{.ItemTitle}}
//	@Description	Create New {{.ItemTitle}}
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.Create{{.ItemTitle}}Request	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.{{.ItemTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}} [post]
func (h *{{.SubModuleTitle}}Handler) Create{{.ItemTitle}}(c *echo.Context) error {
	var req dto.Create{{.ItemTitle}}Request
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Create{{.ItemTitle}}(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── Update{{.ItemTitle}} ────────────────────────────────────────────────────
//
//	@Summary		Update {{.ItemTitle}}
//	@Description	Update {{.ItemTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"{{.ItemTitle}} ID"
//	@Param			body	body		dto.Update{{.ItemTitle}}Request	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.{{.ItemTitle}}Response}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [put]
func (h *{{.SubModuleTitle}}Handler) Update{{.ItemTitle}}(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.Update{{.ItemTitle}}Request
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Update{{.ItemTitle}}(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "{{.ItemTitle}} tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete{{.ItemTitle}} ────────────────────────────────────────────────────
//
//	@Summary		Delete {{.ItemTitle}}
//	@Description	Delete {{.ItemTitle}} by :id
//	@Tags			{{.MainModule}}/{{.SubModule}}
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"{{.ItemTitle}} ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			{{.URLPrefixOpenAPI}}/{id} [delete]
func (h *{{.SubModuleTitle}}Handler) Delete{{.ItemTitle}}(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.Delete{{.ItemTitle}}(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "{{.ItemTitle}} tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
`

var tmplItemMigration = `package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"

	"gorm.io/gorm"
)

//go:embed {{.Timestamp}}_create_{{.TableName}}_table.sql
var {{.ItemName}}SQL string

// Migrate{{.ItemTitle}} menjalankan GORM auto-migration
func Migrate{{.ItemTitle}}(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.{{.ItemTitle}}{})
}

// Migrate{{.ItemTitle}}WithSQL menjalankan migrasi via raw SQL
func Migrate{{.ItemTitle}}WithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec({{.ItemName}}SQL)
	if err != nil {
		log.Printf("Error creating {{.TableName}} table: %v", err)
		return err
	}
	log.Println("{{.TableName}} table migrated successfully")
	return nil
}

// Drop{{.ItemTitle}}Table menghapus tabel (gunakan dengan hati-hati!)
func Drop{{.ItemTitle}}Table(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.{{.ItemTitle}}{})
}
`

var tmplItemSQL = `-- Migration: Create {{.TableName}} table
-- Timestamp: {{.Timestamp}}

CREATE TABLE IF NOT EXISTS {{.TableName}} (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_{{.TableName}}_deleted_at ON {{.TableName}}(deleted_at);
`

var tmplItemFactory = `package factories

import (
	"fmt"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
)

// {{.ItemTitle}}Factory membuat data {{.ItemTitle}} untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type {{.ItemTitle}}Factory struct {
	overrides map[string]interface{}
}

func New{{.ItemTitle}}Factory() *{{.ItemTitle}}Factory {
	return &{{.ItemTitle}}Factory{overrides: make(map[string]interface{})}
}

func (f *{{.ItemTitle}}Factory) With(field string, value interface{}) *{{.ItemTitle}}Factory {
	f.overrides[field] = value
	return f
}

func (f *{{.ItemTitle}}Factory) Make() *models.{{.ItemTitle}} {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("{{.ItemTitle}} %d", idx)
	desc := fmt.Sprintf("Deskripsi {{.ItemTitle}} %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.{{.ItemTitle}}{
		Name:        name,
		Description: &desc,
	}
}

func (f *{{.ItemTitle}}Factory) MakeMany(count int) []*models.{{.ItemTitle}} {
	items := make([]*models.{{.ItemTitle}}, count)
	for i := 0; i < count; i++ {
		items[i] = New{{.ItemTitle}}Factory().Make()
	}
	return items
}
`

var tmplItemSeeder = `package seeders

import (
	"log"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/tests/factories"

	"gorm.io/gorm"
)

// {{.ItemTitle}}Seeder mengelola seeding data {{.ItemTitle}}.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type {{.ItemTitle}}Seeder struct {
	db *gorm.DB
}

func New{{.ItemTitle}}Seeder(db *gorm.DB) *{{.ItemTitle}}Seeder {
	return &{{.ItemTitle}}Seeder{db: db}
}

func (s *{{.ItemTitle}}Seeder) Run() error {
	log.Println("🌱 Seeding {{.TableName}}...")

	items := factories.New{{.ItemTitle}}Factory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat {{.ItemTitle}}: %v", err)
			continue
		}
		log.Printf("   ✅ {{.ItemTitle}} '%s' dibuat.", item.Name)
	}

	log.Println("✅ {{.TableName}} seeding selesai!")
	return nil
}

func (s *{{.ItemTitle}}Seeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data {{.TableName}}...")
	if err := s.db.Exec("DELETE FROM {{.TableName}}").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE {{.TableName}}_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *{{.ItemTitle}}Seeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.{{.ItemTitle}}{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.New{{.ItemTitle}}Factory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
`

var tmplItemRepositoryMock = `package mocks

import (
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
)

// Method di bawah ini ditempelkan ke {{.SubModuleTitle}}RepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/{{.SubModule}}_repository_mock.go).

func (m *{{.SubModuleTitle}}RepositoryMock) Create{{.ItemTitle}}(item *models.{{.ItemTitle}}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *{{.SubModuleTitle}}RepositoryMock) Get{{.ItemTitle}}ByID(id int64) (*models.{{.ItemTitle}}, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.{{.ItemTitle}}), args.Error(1)
}

func (m *{{.SubModuleTitle}}RepositoryMock) List{{.ItemTitle}}(page, pageSize int, filter *dto.Filter{{.ItemTitle}}Request) ([]models.{{.ItemTitle}}, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.{{.ItemTitle}}), args.Get(1).(int64), args.Error(2)
}

func (m *{{.SubModuleTitle}}RepositoryMock) Update{{.ItemTitle}}(item *models.{{.ItemTitle}}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *{{.SubModuleTitle}}RepositoryMock) Delete{{.ItemTitle}}(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
`

var tmplItemServiceTest = `package tests

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/dto"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/models"
	"{{.ProjectModule}}/internal/modules/{{.MainModule}}/{{.SubModule}}/tests/factories"

	appErrors "{{.ProjectModule}}/internal/shared/errors"
)

// Catatan: suite {{.SubModuleTitle}}ServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di {{.SubModule}}_service_test.go. File ini HANYA menambah skenario test untuk
// {{.ItemTitle}}, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module {{.SubModule}}).

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Create{{.ItemTitle}}_Superadmin_Success() {
	req := &dto.Create{{.ItemTitle}}Request{Name: "Test {{.ItemTitle}}"}
	actor := superadminActor()

	s.repo.On("Create{{.ItemTitle}}", mock.AnythingOfType("*models.{{.ItemTitle}}")).Return(nil)

	result, err := s.svc.Create{{.ItemTitle}}(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Create{{.ItemTitle}}_Forbidden() {
	req := &dto.Create{{.ItemTitle}}Request{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.Create{{.ItemTitle}}(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Create{{.ItemTitle}}_RepoError() {
	req := &dto.Create{{.ItemTitle}}Request{Name: "Test"}
	actor := superadminActor()

	s.repo.On("Create{{.ItemTitle}}", mock.AnythingOfType("*models.{{.ItemTitle}}")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Create{{.ItemTitle}}(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Get{{.ItemTitle}}ByID_Success() {
	actor := superadminActor()
	item := factories.New{{.ItemTitle}}Factory().Make()
	item.ID = 1

	s.repo.On("Get{{.ItemTitle}}ByID", int64(1)).Return(item, nil)

	result, err := s.svc.Get{{.ItemTitle}}ByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Get{{.ItemTitle}}ByID_NotFound() {
	actor := superadminActor()

	s.repo.On("Get{{.ItemTitle}}ByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Get{{.ItemTitle}}ByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Get{{.ItemTitle}}ByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.Get{{.ItemTitle}}ByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_List{{.ItemTitle}}_Success() {
	actor := superadminActor()
	filter := &dto.Filter{{.ItemTitle}}Request{}
	items := []models.{{.ItemTitle}}{
		*factories.New{{.ItemTitle}}Factory().Make(),
		*factories.New{{.ItemTitle}}Factory().Make(),
	}

	s.repo.On("List{{.ItemTitle}}", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.List{{.ItemTitle}}(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_List{{.ItemTitle}}_Forbidden() {
	actor := regularActor()
	filter := &dto.Filter{{.ItemTitle}}Request{}
	s.mockNoPermissions()

	result, total, err := s.svc.List{{.ItemTitle}}(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Update{{.ItemTitle}}_Success() {
	actor := superadminActor()
	existing := factories.New{{.ItemTitle}}Factory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.Update{{.ItemTitle}}Request{Name: &newName}

	s.repo.On("Get{{.ItemTitle}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Update{{.ItemTitle}}", mock.AnythingOfType("*models.{{.ItemTitle}}")).Return(nil)

	result, err := s.svc.Update{{.ItemTitle}}(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Update{{.ItemTitle}}_NotFound() {
	actor := superadminActor()
	req := &dto.Update{{.ItemTitle}}Request{}

	s.repo.On("Get{{.ItemTitle}}ByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Update{{.ItemTitle}}(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Update{{.ItemTitle}}_Forbidden() {
	actor := regularActor()
	req := &dto.Update{{.ItemTitle}}Request{}
	s.mockNoPermissions()

	result, err := s.svc.Update{{.ItemTitle}}(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Delete{{.ItemTitle}}_Success() {
	actor := superadminActor()
	existing := factories.New{{.ItemTitle}}Factory().Make()
	existing.ID = 1

	s.repo.On("Get{{.ItemTitle}}ByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete{{.ItemTitle}}", int64(1)).Return(nil)

	err := s.svc.Delete{{.ItemTitle}}(1, actor)

	s.NoError(err)
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Delete{{.ItemTitle}}_NotFound() {
	actor := superadminActor()

	s.repo.On("Get{{.ItemTitle}}ByID", int64(999)).Return(nil, nil)

	err := s.svc.Delete{{.ItemTitle}}(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *{{.SubModuleTitle}}ServiceTestSuite) Test_Delete{{.ItemTitle}}_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.Delete{{.ItemTitle}}(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
`
