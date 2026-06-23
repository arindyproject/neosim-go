package services

import (
	kontakContracts "neosim_go/internal/modules/kepegawaian/kontak/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     kontakContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

// NewKepegawaianKontakService membuat instance service baru
func NewKepegawaianKontakService(
	repo kontakContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) kontakContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}
