package services

import (
	lampiranContracts "neosim_go/internal/modules/kepegawaian/lampiran/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     lampiranContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

// NewKepegawaianLampiranService membuat instance service baru
func NewKepegawaianLampiranService(
	repo lampiranContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) lampiranContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}
