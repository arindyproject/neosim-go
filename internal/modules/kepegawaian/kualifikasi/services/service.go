package services

import (
	kualifikasiContracts "neosim_go/internal/modules/kepegawaian/kualifikasi/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     kualifikasiContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

// NewKepegawaianKualifikasiService membuat instance service baru
func NewKepegawaianKualifikasiService(
	repo kualifikasiContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) kualifikasiContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}
