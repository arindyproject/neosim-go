# Neosim Go

Backend REST API berbasis Go menggunakan Echo v5 dan GORM.

---

## Tech Stack

- **Language**: Go
- **Framework**: Echo v5
- **ORM**: GORM
- **Database**: PostgreSQL
- **Auth**: JWT

---

## Struktur Folder

```
neosim_go/
├── cmd/
│   ├── api/
│   │   └── main.go              # Entry point server
│   ├── migrate/
│   │   └── main.go              # Entry point migrasi
│   └── seed/
│       └── main.go              # Entry point seeder
│
├── config/
│   ├── config.go                # Konfigurasi aplikasi
│   └── database.go              # Konfigurasi database
│
└── internal/
    ├── apps/
    │   ├── apps.go              # Entry point registrasi module
    │   └── registry.go          # Registry module (routes + migration)
    │
    └── modules/
        └── users/
            ├── contracts/       # Interface Service & Repository
            ├── dto/             # Request & Response structs
            │   ├── user_request.go
            │   └── user_response.go
            ├── handlers/        # HTTP handlers
            │   └── user_handler.go
            ├── migrations/      # File migrasi SQL
            │   ├── migrate.go
            │   └── 001_create_users_table.sql
            ├── models/          # GORM models
            │   └── users.go
            ├── repositories/    # Database queries
            │   └── user_repository.go
            ├── services/        # Business logic
            │   └── user_service.go
            ├── tests/
            │   ├── factories/   # Data generator untuk test
            │   │   └── user_factory.go
            │   ├── seeders/     # Seeder ke database
            │   │   └── user_seeder.go
            │   └── helpers/     # Helper untuk test
            │       └── db_helper.go
            ├── module.go        # Wire semua layer
            ├── register.go      # Auto-register ke registry
            └── routes.go        # Definisi routes
```

---

## Cara Menjalankan

### Prasyarat

- Go 1.21+
- PostgreSQL
- Make

### Clone & Setup

```bash
git clone <repo-url>
cd neosim_go
cp .env.example .env   # Sesuaikan konfigurasi database
go mod tidy
```

### Jalankan Server

```bash
make run
# atau
go run ./cmd/api/main.go
```

---

## Perintah Make

### 🚀 Server

| Perintah     | Deskripsi                 |
| ------------ | ------------------------- |
| `make run`   | Jalankan API server       |
| `make build` | Build binary ke `bin/api` |
| `make clean` | Hapus build artifacts     |
| `make test`  | Jalankan semua tests      |

---

### 🗄️ Migrasi

| Perintah                | Deskripsi                  |
| ----------------------- | -------------------------- |
| `make migrate-dev`      | GORM auto-migration (DEV)  |
| `make migrate-prod`     | GORM auto-migration (PROD) |
| `make migrate-sql`      | SQL migration (DEV)        |
| `make migrate-sql-prod` | SQL migration (PROD)       |

#### Fresh Migration (Drop All + Re-migrate)

> ⚠️ **Hati-hati** — Semua data akan hilang!

| Perintah                      | Deskripsi                       | Konfirmasi                    |
| ----------------------------- | ------------------------------- | ----------------------------- |
| `make migrate-fresh-dev`      | Drop + migrate ulang (DEV)      | `yes/no`                      |
| `make migrate-fresh-dev-sql`  | Drop + SQL migrate ulang (DEV)  | `yes/no`                      |
| `make migrate-fresh-prod`     | Drop + migrate ulang (PROD)     | Ketik `PRODUCTION` + `yes/no` |
| `make migrate-fresh-prod-sql` | Drop + SQL migrate ulang (PROD) | Ketik `PRODUCTION` + `yes/no` |

---

### 🌱 Seeder

| Perintah                  | Deskripsi                                      |
| ------------------------- | ---------------------------------------------- |
| `make seed`               | Jalankan seeder DEV (skip jika data sudah ada) |
| `make seed-prod`          | Jalankan seeder PROD                           |
| `make seed-fresh`         | Hapus semua data lalu seed ulang (DEV)         |
| `make migrate-seed`       | Migrate + seed sekaligus (DEV)                 |
| `make migrate-fresh-seed` | Fresh migrate + seed (DEV)                     |

#### Data yang di-seed

| Tipe         | Jumlah | Username                        | Password      |
| ------------ | ------ | ------------------------------- | ------------- |
| Superuser    | 1      | `superadmin`                    | `password123` |
| Staff        | 3      | `staff_1`, `staff_2`, `staff_3` | `password123` |
| Regular User | 10     | random `user_XXXXX`             | `password123` |

---

## API Endpoints

### Users

| Method   | Endpoint                            | Deskripsi                   |
| -------- | ----------------------------------- | --------------------------- |
| `GET`    | `/api/v1/users`                     | List semua user (paginated) |
| `GET`    | `/api/v1/users/:id`                 | Detail user by ID           |
| `GET`    | `/api/v1/users/username/:username`  | Detail user by username     |
| `POST`   | `/api/v1/users`                     | Buat user baru              |
| `PUT`    | `/api/v1/users/:id`                 | Update user                 |
| `DELETE` | `/api/v1/users/:id`                 | Hapus user                  |
| `PUT`    | `/api/v1/users/:id/change-password` | Ganti password              |
| `GET`    | `/api/v1/users/:id/settings`        | Ambil settings user         |
| `PUT`    | `/api/v1/users/:id/settings`        | Update settings user        |

### Query Params (List)

| Param       | Default | Deskripsi                         |
| ----------- | ------- | --------------------------------- |
| `page`      | `1`     | Nomor halaman                     |
| `page_size` | `10`    | Jumlah data per halaman (max 100) |

### Contoh Response

**Success:**

```json
{
    "status": true,
    "message": "Berhasil mengambil data user",
    "data": { ... }
}
```

**Paginated:**

```json
{
    "status": true,
    "message": "Berhasil mengambil data user",
    "data": {
        "items": [ ... ],
        "pagination": {
            "total_items": 100,
            "total_pages": 10,
            "current_page": 1,
            "per_page": 10,
            "next_page": 2,
            "prev_page": null
        }
    }
}
```

**Validation Error:**

```json
{
  "status": false,
  "message": "Validasi gagal",
  "errors": [
    { "field": "email", "message": "email harus berupa email yang valid" },
    { "field": "password", "message": "password minimal 8 karakter" }
  ]
}
```

---

## Menambah Module Baru

Cukup 2 langkah saat ingin tambah module baru (misal `roles`):

**1.** Buat `internal/modules/roles/register.go`:

```go
package roles

import (
    "database/sql"
    "neosim_go/internal/apps"
    "neosim_go/internal/modules/roles/models"
    "github.com/labstack/echo/v5"
    "gorm.io/gorm"
)

type module struct{ db *gorm.DB }

func init() { apps.Register(&module{}) }

func (m *module) SetDB(db *gorm.DB)          { m.db = db }
func (m *module) InitRoutes(e *echo.Echo)     { NewModule(m.db).InitRoutes(e) }
func (m *module) Models() []interface{}       { return []interface{}{&models.Role{}} }
func (m *module) SeedData(db *gorm.DB) error  { return nil }
func (m *module) MigrateSQL(sqlDB *sql.DB) error { return nil }
```

**2.** Tambah blank import di `internal/apps/apps.go` dan `cmd/migrate/main.go`:

```go
_ "neosim_go/internal/modules/roles"
```

Tidak ada file lain yang perlu diubah.
