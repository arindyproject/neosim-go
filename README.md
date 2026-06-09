# Neosim Go

Backend REST API berbasis Go menggunakan Echo v5 dan GORM.

---

## Ringkasan

Neosim Go adalah backend API microservice dengan modul: `auth`, `users`, dan `rbac`.
Aplikasi menyediakan autentikasi JWT, manajemen user, serta role dan permission untuk kontrol akses.

---

## Fitur Utama

- Auth: Login, register, logout, refresh token, forgot password, reset password.
- User management: CRUD user, change password, settings, upload foto profil.
- RBAC: Permission CRUD, role CRUD, assign/revoke permissions ke role, assign/revoke role ke user, direct permission.
- Protected routes dengan JWT dan middleware role/permission.
- Swagger API documentation (di environment `development`).
- PostgreSQL + Redis untuk session/login limit.

---

## Struktur Folder Utama

```
neosim-go/
├── cmd/
│   ├── api/        # entry point server
│   ├── migrate/    # entry point migrasi DB
│   └── seed/       # entry point seeder
├── config/         # konfigurasi environment, database, Redis
├── internal/       # domain aplikasi dan modul
│   ├── apps/       # registry module, inisialisasi route/module
│   ├── modules/
│   │   ├── auth/   # autentikasi JWT, auth routes
│   │   ├── users/  # user CRUD, settings, upload foto
   │   │   └── rbac/   # role-based access control
│   └── shared/     # utilities, middleware, response, errors
├── docs/           # generated Swagger docs
└── config/.env.dev  # contoh environment dev
```

---

## Prasyarat

- Go 1.21+
- PostgreSQL
- Redis
- Make

---

## Setup

1. Clone repository:

```bash
git clone <repo-url>
cd neosim-go
```

2. Install dependency:

```bash
go mod tidy
```

3. Salin environment dev:

```bash
cp config/.env.dev .env
```

4. Sesuaikan `.env` jika diperlukan.

---

## Menjalankan Aplikasi

### Jalankan server

```bash
make run
```

atau:

```bash
go run ./cmd/api/main.go
```

### Swagger UI (development only)

Saat `ENV=development`, Swagger tersedia di:

```
http://localhost:1323/swagger/index.html
```

---

## Perintah Make

| Perintah                  | Deskripsi                        |
| ------------------------- | -------------------------------- |
| `make run`                | Jalankan API server              |
| `make build`              | Build binary ke `bin/api`        |
| `make clean`              | Hapus build artifacts            |
| `make test`               | Jalankan semua tests             |
| `make migrate-dev`        | GORM migrate database untuk DEV  |
| `make migrate-prod`       | GORM migrate database untuk PROD |
| `make migrate-sql`        | SQL migration untuk DEV          |
| `make migrate-sql-prod`   | SQL migration untuk PROD         |
| `make seed`               | Jalankan seeder DEV              |
| `make seed-prod`          | Jalankan seeder PROD             |
| `make migrate-seed`       | Migrate + seed DEV               |
| `make migrate-fresh-seed` | Fresh migrate + seed DEV         |

> ⚠️ `migrate-fresh-*` akan menghapus data lama sebelum deploy ulang.

---

## Konfigurasi Environment

Gunakan `config/.env.dev` sebagai contoh. Variabel penting:

| Variabel                       | Default                 | Deskripsi                               |
| ------------------------------ | ----------------------- | --------------------------------------- |
| `BASE_URL`                     | `http://localhost:1323` | URL base aplikasi                       |
| `SERVER_PORT`                  | `1323`                  | Port server                             |
| `ENV`                          | `development`           | Environment                             |
| `LOG_LEVEL`                    | `debug`                 | Level logging                           |
| `DATABASE_URL`                 | `postgres://...`        | DSN PostgreSQL                          |
| `DB_HOST`                      | `localhost`             | Host DB                                 |
| `DB_PORT`                      | `5432`                  | Port DB                                 |
| `DB_USER`                      | `postgres`              | Username DB                             |
| `DB_PASSWORD`                  |                         | Password DB                             |
| `DB_NAME`                      | `neosim`                | Nama database                           |
| `DB_SSL_MODE`                  | `disable`               | SSL mode postgres                       |
| `REDIS_HOST`                   | `localhost`             | Host Redis                              |
| `REDIS_PORT`                   | `6379`                  | Port Redis                              |
| `JWT_SECRET`                   |                         | Secret JWT                              |
| `JWT_ISSUER`                   | `neosim`                | Issuer JWT                              |
| `JWT_ACCESS_TOKEN_EXP_MINUTES` | `15`                    | Expired access token dalam menit        |
| `JWT_REFRESH_TOKEN_EXP_DAYS`   | `7`                     | Expired refresh token dalam hari        |
| `PASSWORD_MIN_LENGTH`          | `6`                     | Minimum panjang password                |
| `IS_REGISTRATION_ACTIVE`       | `true`                  | Aktifkan registrasi user                |
| `AUTO_ACTIVE_USER`             | `true`                  | Aktifkan user langsung setelah register |

---

## Arsitektur

Aplikasi dibangun sebagai modul: `auth`, `users`, dan `rbac`.
Setiap modul memiliki lapisan:

- `contracts` untuk interface service/repository
- `dto` untuk request/response
- `handlers` untuk HTTP handler
- `services` untuk business logic
- `repositories` untuk akses database
- `migrations` untuk skrip DB
- `middlewares` untuk proteksi route

Modul terdaftar lewat `internal/apps/apps.go` dan blank import di `cmd/api/main.go`.

---

## Auth

### Endpoint Auth

| Method | Endpoint                       | Deskripsi                                                      |
| ------ | ------------------------------ | -------------------------------------------------------------- |
| `POST` | `/api/v1/auth/login`           | Login menggunakan `identifier` (username/email) dan `password` |
| `POST` | `/api/v1/auth/register`        | Register user baru                                             |
| `POST` | `/api/v1/auth/refresh`         | Refresh access token dengan refresh token                      |
| `POST` | `/api/v1/auth/forgot-password` | Minta reset password                                           |
| `POST` | `/api/v1/auth/reset-password`  | Reset password menggunakan token                               |
| `POST` | `/api/v1/auth/logout`          | Logout dari session saat ini                                   |
| `POST` | `/api/v1/auth/logout-all`      | Logout dari semua device                                       |

### Request Body

#### Login

```json
{
  "identifier": "user@example.com",
  "password": "secret123"
}
```

#### Register

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

#### Refresh Token

```json
{
  "refresh_token": "<refresh_token>"
}
```

#### Logout

```json
{
  "refresh_token": "<refresh_token>"
}
```

#### Forgot Password

```json
{
  "identifier": "user@example.com"
}
```

#### Reset Password

```json
{
  "token": "<reset_token>",
  "new_password": "newPassword123",
  "confirm_password": "newPassword123"
}
```

### Response Token

Respons sukses login/refresh mengembalikan:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "name": "John Doe",
    "is_superadmin": false,
    "is_staff": false,
    "is_verified": true
  }
}
```

---

## Users

### Endpoint Users

| Method   | Endpoint                            | Deskripsi                       |
| -------- | ----------------------------------- | ------------------------------- |
| `GET`    | `/api/v1/users`                     | List semua user                 |
| `GET`    | `/api/v1/users/:id`                 | Ambil user berdasarkan ID       |
| `GET`    | `/api/v1/users/username/:username`  | Ambil user berdasarkan username |
| `POST`   | `/api/v1/users`                     | Buat user baru                  |
| `PUT`    | `/api/v1/users/:id`                 | Update user                     |
| `DELETE` | `/api/v1/users/:id`                 | Hapus user                      |
| `GET`    | `/api/v1/users/deleted`             | List user yang dihapus          |
| `PUT`    | `/api/v1/users/:id/change-password` | Ganti password user             |
| `POST`   | `/api/v1/users/:id/reset-password`  | Reset password user             |
| `GET`    | `/api/v1/users/:id/settings`        | Ambil settings user             |
| `PUT`    | `/api/v1/users/:id/settings`        | Update settings user            |
| `PUT`    | `/api/v1/users/:id/photo`           | Upload foto profil              |
| `DELETE` | `/api/v1/users/:id/photo`           | Hapus foto profil               |

> Semua endpoint user membutuhkan header `Authorization: Bearer <access_token>`.

### Catatan User

- `username`, `email`, dan `password` divalidasi.
- `password` default minimal 8 karakter untuk register.
- User bisa mengakses foto dari folder `/uploads` jika token JWT valid.

---

## RBAC (Role & Permission)

RBAC berfungsi sebagai kontrol akses: pengguna boleh melakukan aksi hanya jika:

1. Memiliki role yang benar,
2. Memiliki permission yang dibutuhkan, atau
3. Berstatus `superadmin`.

### Konsep Utama

- Permission: tindakan spesifik yang dapat diberikan ke role.
- Role: grup hak akses. Role bisa memiliki banyak permission.
- User role: relasi user ke role.
- Direct permission: hak akses langsung ke user (jika ada) — saat ini tersedia endpoint untuk assign direct permission.
- Superadmin: akses penuh ke semua route RBAC dan bypass cek permission/role.

### Endpoint RBAC

#### Permissions

| Method   | Endpoint                       | Deskripsi         |
| -------- | ------------------------------ | ----------------- |
| `GET`    | `/api/v1/rbac/permissions`     | List permission   |
| `POST`   | `/api/v1/rbac/permissions`     | Buat permission   |
| `GET`    | `/api/v1/rbac/permissions/:id` | Ambil permission  |
| `PUT`    | `/api/v1/rbac/permissions/:id` | Update permission |
| `DELETE` | `/api/v1/rbac/permissions/:id` | Hapus permission  |

#### Roles

| Method   | Endpoint                 | Deskripsi   |
| -------- | ------------------------ | ----------- |
| `GET`    | `/api/v1/rbac/roles`     | List role   |
| `POST`   | `/api/v1/rbac/roles`     | Buat role   |
| `GET`    | `/api/v1/rbac/roles/:id` | Ambil role  |
| `PUT`    | `/api/v1/rbac/roles/:id` | Update role |
| `DELETE` | `/api/v1/rbac/roles/:id` | Hapus role  |

#### Role ↔ Permission

| Method   | Endpoint                             | Deskripsi                           |
| -------- | ------------------------------------ | ----------------------------------- |
| `POST`   | `/api/v1/rbac/roles/:id/permissions` | Assign permission ke role           |
| `PUT`    | `/api/v1/rbac/roles/:id/permissions` | Sync permission role (replace list) |
| `DELETE` | `/api/v1/rbac/roles/:id/permissions` | Revoke permission dari role         |

#### User ↔ Role

| Method   | Endpoint                            | Deskripsi           |
| -------- | ----------------------------------- | ------------------- |
| `GET`    | `/api/v1/rbac/users/:user_id/roles` | List roles user     |
| `POST`   | `/api/v1/rbac/users/:user_id/roles` | Assign role ke user |
| `PUT`    | `/api/v1/rbac/users/:user_id/roles` | Sync roles user     |
| `DELETE` | `/api/v1/rbac/users/:user_id/roles` | Revoke role user    |

#### User Permissions

| Method | Endpoint                                  | Deskripsi                        |
| ------ | ----------------------------------------- | -------------------------------- |
| `GET`  | `/api/v1/rbac/users/:user_id/permissions` | List semua permissions user      |
| `POST` | `/api/v1/rbac/users/:user_id/permissions` | Assign direct permission ke user |

> Semua endpoint RBAC membutuhkan JWT dan umumnya `superadmin` atau permission khusus.

### Request Payload RBAC

#### Create Permission

```json
{
  "name": "manage_users",
  "display_name": "Manage Users",
  "description": "Mengatur user dan role",
  "resource": "users",
  "action": "manage"
}
```

#### Update Permission

```json
{
  "display_name": "Kelola Pengguna",
  "description": "Akses untuk mengelola user"
}
```

#### Create Role

```json
{
  "name": "admin",
  "display_name": "Administrator",
  "description": "Role dengan hak akses administratif"
}
```

#### Assign Permissions ke Role

```json
{
  "permission_ids": [1, 2, 3]
}
```

#### Assign Roles ke User

```json
{
  "role_ids": [1, 2]
}
```

#### Assign Direct Permission ke User

```json
{
  "permission_id": 1,
  "is_granted": true
}
```

---

## Middleware Authorization Rules

Akses RBAC dijaga dengan middleware berikut:

- `RequirePermission(repo, permission)`: pastikan user punya permission tertentu.
- `RequireAnyPermission(repo, permissions...)`: pastikan user punya salah satu permission.
- `RequireRole(repo, roleName)`: pastikan user punya role tertentu.
- `RequireSuperadmin()`: hanya superadmin.
- `RequireSelf()`: hanya pemilik data saja.
- `RequireSelfOrPermission(repo, permission)`: pemilik data atau pemilik permission.
- `RequireSelfOrRole(repo, roleName)`: pemilik data atau pemilik role.

> Superadmin otomatis melewati semua pengecekan permission/role.

---

## Token & Headers

Semua endpoint privat memerlukan header:

```
Authorization: Bearer <access_token>
```

Jika token tidak valid atau expired, server mengembalikan status `401 Unauthorized`.

---

## Cara Menambah Module Baru

1. Buat file `internal/modules/<module>/register.go` dengan struktur module.
2. Tambahkan blank import di `cmd/api/main.go` dan `cmd/migrate/main.go`.
3. Buat route, handler, service, repository, dto, dan migration sesuai kebutuhan.

> Modul baru akan otomatis didaftarkan oleh mekanisme registry di `internal/apps`.

---

## Contoh Alur RBAC

1. Buat permission.
2. Buat role.
3. Assign permission ke role.
4. Assign role ke user.
5. User login, lalu akses endpoint berdasarkan permission/role.

Contoh: jika route membutuhkan permission `manage_users`, maka user harus memiliki role yang memuat permission tersebut atau menjadi `superadmin`.

---

## Catatan

- Jika Anda ingin membuka akses static file `/uploads`, pastikan server memeriksa JWT.
- Untuk debugging, route terdaftar dicetak di console saat server start.
- Gunakan `config/.env.dev` untuk environment development.

---

## Lisensi

MIT
