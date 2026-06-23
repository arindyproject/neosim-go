package services

import (
	identifikasiContracts "neosim_go/internal/modules/kepegawaian/identifikasi/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     identifikasiContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

// NewKepegawaianIdentifikasiService membuat instance service baru
func NewKepegawaianIdentifikasiService(
	repo identifikasiContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) identifikasiContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}
