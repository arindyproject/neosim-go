package services

import (
	"context"
	"neosim_go/config"
	identifierContracts "neosim_go/internal/modules/kepegawaian/identifier/contracts"

	authContracts "neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	userContracts "neosim_go/internal/modules/users/contracts"
	he "neosim_go/internal/shared/httputil"
)

// service adalah satu-satunya struct service untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat struct service baru — method
// CRUD & permission-nya ditempelkan langsung ke struct ini (mis.
// services/tag_service.go, services/tag_permission.go), dan repo field
// di bawah ini otomatis mencakup method item begitu contracts.Repository
// di-embed dengan interface repository item (lihat contracts/interfaces.go).
type service struct {
	repo     identifierContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	userRepo userContracts.Repository
	cfg      *config.Config
}

// NewKepegawaianIdentifierService membuat instance service baru
func NewKepegawaianIdentifierService(
	repo identifierContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	cfg *config.Config,
) identifierContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// buildCreator mengambil data creator user
func (s *service) buildCreator(ctx context.Context, createdBy *int64) *he.UserData {
	if createdBy == nil {
		return nil
	}
	creator, err := s.userRepo.GetByID(ctx, *createdBy)
	if err != nil || creator == nil {
		return nil
	}
	return &he.UserData{
		ID:       creator.ID,
		Username: creator.Username,
		Name:     creator.Name,
	}
}

// ── helper: build creator/updater maps ───────────────────────────────────────

func (s *service) buildAuditMaps(ctx context.Context, items []models.KepegawaianIdentifier) (map[int64]*he.UserData, map[int64]*he.UserData) {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			idSet[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			idSet[*item.UpdatedBy] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, ids) // ← 1 query total, bukan 40
	if err != nil {
		return map[int64]*he.UserData{}, map[int64]*he.UserData{}
	}

	userMap := make(map[int64]*he.UserData, len(users))
	for _, u := range users {
		userMap[u.ID] = &he.UserData{ID: u.ID, Username: u.Username, Name: u.Name}
	}
	// creator dan updater sekarang share map yang sama — reuse otomatis, kode lebih pendek juga
	return userMap, userMap
}
