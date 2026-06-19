package services

import (
	masterContracts "neosim_go/internal/modules/master/master/contracts"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/shared/cache"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo     masterContracts.Repository
	rbacRepo rbacContracts.RBACRepository //RBAC
	authRepo authContracts.AuthRepository //AUTH
	cache    *cache.Manager               // <--- Gunakan Cache Manager
}

// NewMasterService membuat instance service baru
func NewMasterService(
	repo masterContracts.Repository,
	rbacRepo rbacContracts.RBACRepository, //RBAC
	authRepo authContracts.AuthRepository, //AUTH
	cacheManager *cache.Manager, // <--- Terima Cache Manager
) masterContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo, //RBAC
		authRepo: authRepo, //AUTH
		cache:    cacheManager,
	}
}
