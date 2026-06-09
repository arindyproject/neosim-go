# Contoh Penggunaan RBAC (Routes & Service)

Berisi contoh snippet penggunaan middleware RBAC pada route dan contoh pengecekan permission di layer service.

## Contoh: Proteksi Route (Allow Self Or Permission)

Contoh ini memperlihatkan cara mendaftarkan route `PUT /api/v1/users/:id` sehingga:

- Pemilik data (self) boleh mengupdate profilnya sendiri.
- Selain itu, user yang punya permission `users:update` juga boleh mengupdate user lain.

```go
// Route example (snippet)
func RegisterUserRoutes(e *echo.Echo, h *users.Handler, rbacRepo rbacContracts.RBACRepository, jwtManager *utils.JWTManager, db *gorm.DB) {
    // JWT middleware dibangun dari jwtManager + db (untuk cek isSuperadmin realtime)
    jwt := authMiddlewares.JWTMiddleware(jwtManager, db)

    // Group route yang butuh autentikasi
    protected := e.Group("/api/v1/users", jwt)

    // Update user: boleh jika pemilik data (RequireSelf) atau punya permission users:update
    protected.PUT("/:id", h.UpdateUserHandler,
        rbacMiddlewares.RequireSelfOrPermission(rbacRepo, rbacModels.PermUsersUpdate),
    )

    // Delete user: hanya untuk yang punya permission users:delete (atau superadmin)
    protected.DELETE("/:id", h.DeleteUserHandler,
        rbacMiddlewares.RequirePermission(rbacRepo, rbacModels.PermUsersDelete),
    )
}
```

Catatan:

- `rbacMiddlewares.RequireSelfOrPermission` akan mengizinkan akses jika user adalah pemilik data (self) atau memiliki permission yang diberikan.
- `RequirePermission` hanya mengizinkan jika memiliki permission tersebut; `RequireSuperadmin` bisa ditambahkan jika ingin membatasi lebih kuat.

## Contoh: Pengecekan Permission di Service

Kadang Anda perlu melakukan pengecekan permission di servis (bukan di middleware) — misal ketika aksi dipicu dari background job atau proses internal.

```go
// Service example (snippet)
type UserService struct {
    userRepo userContracts.Repository
    rbacRepo rbacContracts.RBACRepository
}

// Hapus user dengan cek permission: actorID adalah ID user yang melakukan request
func (s *UserService) DeleteUser(targetUserID, actorID int64) error {
    // 1) Jika actor adalah superadmin, izinkan langsung
    isSuper, err := s.rbacRepo.IsSuperadmin(actorID)
    if err != nil {
        return err
    }
    if !isSuper {
        // 2) Jika bukan superadmin, cek apakah actor punya permission `users:delete`
        ok, err := s.rbacRepo.HasPermission(actorID, rbacModels.PermUsersDelete)
        if err != nil {
            return err
        }
        if !ok {
            return appErrors.Forbidden("permission 'users:delete' diperlukan")
        }
    }

    // 3) Lanjutkan business logic untuk menghapus user
    return s.userRepo.DeleteByID(targetUserID)
}
```

Catatan:

- Gunakan `rbacRepo.IsSuperadmin` bila Anda butuh pengecualian superadmin.
- `HasPermission` menggabungkan permission dari role dan direct permission sehingga pengecekan cukup pada satu method.
- Kode di atas memakai `appErrors.Forbidden` (internal/shared/errors) untuk konsistensi response.

## Tip Praktis

- Untuk endpoint HTTP prefer melindungi lewat middleware (lebih konsisten dan terpusat).
- Gunakan pengecekan di service jika aksi dipanggil di luar konteks request HTTP atau ada kebutuhan validasi bisnis ekstra.

---

File ini hanya contoh: salin snippet ke module Anda (`routes.go` atau `service.go`) dan sesuaikan import serta nama package/handler/repo yang ada di proyek Anda.
