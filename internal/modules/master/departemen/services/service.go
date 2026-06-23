package services

import (
	departemenContracts "neosim_go/internal/modules/master/departemen/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
)

type service struct {
	repo     departemenContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

// NewMasterDepartemenService membuat instance service baru
func NewMasterDepartemenService(
	repo departemenContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) departemenContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
	}
}
