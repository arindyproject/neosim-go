
package services

import (

	artikelContracts "neosim_go/internal/modules/artikel/artikel/contracts"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo artikelContracts.Repository
	rbacRepo rbacContracts.RBACRepository	//RBAC
	authRepo authContracts.AuthRepository	//AUTH
}

// NewArtikelService membuat instance service baru
func NewArtikelService(
	repo artikelContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,	//RBAC
	authRepo authContracts.AuthRepository,	//AUTH
) artikelContracts.Service {
	return &service{
		repo: repo,
		rbacRepo: rbacRepo,	//RBAC
		authRepo: authRepo,	//AUTH
	}
}
