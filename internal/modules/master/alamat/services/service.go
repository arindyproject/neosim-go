package services

import (
	alamatContracts "neosim_go/internal/modules/master/alamat/contracts"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"

	//cache

	"neosim_go/internal/shared/cache"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo     alamatContracts.Repository
	rbacRepo rbacContracts.RBACRepository //RBAC
	authRepo authContracts.AuthRepository //AUTH
	cache    *cache.Manager               // <--- Gunakan Cache Manager
}

// NewMasterAlamatService membuat instance service baru
func NewMasterAlamatService(
	repo alamatContracts.Repository,
	rbacRepo rbacContracts.RBACRepository, //RBAC
	authRepo authContracts.AuthRepository, //AUTH
	cacheManager *cache.Manager, // <--- Terima Cache Manager
) alamatContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo, //RBAC
		authRepo: authRepo, //AUTH
		cache:    cacheManager,
	}
}
