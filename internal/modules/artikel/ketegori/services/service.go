
package services

import (

	ketegoriContracts "neosim_go/internal/modules/artikel/ketegori/contracts"

	//RBAC AUTH----------------------------------------
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

// ─── Init ───────────────────────────────────────────────────────────────────────
type service struct {
	repo ketegoriContracts.Repository
	rbacRepo rbacContracts.RBACRepository	//RBAC
	authRepo authContracts.AuthRepository	//AUTH
}

// NewArtikelKetegoriService membuat instance service baru
func NewArtikelKetegoriService(
	repo ketegoriContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,	//RBAC
	authRepo authContracts.AuthRepository,	//AUTH
) ketegoriContracts.Service {
	return &service{
		repo: repo,
		rbacRepo: rbacRepo,	//RBAC
		authRepo: authRepo,	//AUTH
	}
}
