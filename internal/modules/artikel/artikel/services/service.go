package services

import (
	"neosim_go/config"
	artikelContracts "neosim_go/internal/modules/artikel/artikel/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     artikelContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	cfg      *config.Config
}

// NewArtikelService membuat instance service baru
func NewArtikelService(
	repo artikelContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	cfg    *config.Config,
) artikelContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
		cfg:      cfg,
	}
}
