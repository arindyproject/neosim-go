package services

import (
	"neosim_go/config"
	kontakContracts "neosim_go/internal/modules/kepegawaian/kontak/contracts"
	"neosim_go/internal/modules/kepegawaian/kontak/dto"

	"neosim_go/internal/modules/kepegawaian/kontak/models"
	authContracts "neosim_go/internal/modules/auth/contracts"
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
	repo     kontakContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	userRepo userContracts.Repository
	cfg      *config.Config
}

// NewKepegawaianKontakService membuat instance service baru
func NewKepegawaianKontakService(
	repo kontakContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	cfg    *config.Config,
) kontakContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// buildCreator mengambil data creator user
func (s *service) buildCreator(createdBy *int64) *he.UserData {
	if createdBy == nil {
		return nil
	}
	creator, err := s.userRepo.GetByID(*createdBy)
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

func (s *service) buildAuditMaps(items []models.KepegawaianKontak) (map[int64]*he.UserData, map[int64]*he.UserData) {
	fetchUser := func(id int64) (*he.UserData, error) {
		user, err := s.userRepo.GetByID(id)
		if err != nil || user == nil {
			return nil, err
		}
		return &he.UserData{ID: user.ID, Username: user.Username, Name: user.Name}, nil
	}

	creatorIDs := make(map[int64]struct{})
	updaterIDs := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			creatorIDs[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			updaterIDs[*item.UpdatedBy] = struct{}{}
		}
	}

	creatorsMap := make(map[int64]*he.UserData)
	for id := range creatorIDs {
		if data, err := fetchUser(id); err == nil && data != nil {
			creatorsMap[id] = data
		}
	}

	updatersMap := make(map[int64]*he.UserData)
	for id := range updaterIDs {
		if data, ok := creatorsMap[id]; ok {
			updatersMap[id] = data
		} else if data, err := fetchUser(id); err == nil && data != nil {
			updatersMap[id] = data
		}
	}

	return creatorsMap, updatersMap
}
// ── helper: convert items to responses ───────────────────────────────────────
func toKepegawaianKontakResponses(
	items []models.KepegawaianKontak,
	creatorsMap, updatersMap map[int64]*he.UserData,
) []dto.KepegawaianKontakResponse {
	responses := make([]dto.KepegawaianKontakResponse, 0, len(items))
	for _, item := range items {
		var creator, updater *he.UserData
		if item.CreatedBy != nil {
			creator = creatorsMap[*item.CreatedBy]
		}
		if item.UpdatedBy != nil {
			updater = updatersMap[*item.UpdatedBy]
		}
		responses = append(responses, *dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
			KepegawaianKontak: &item,
			Creator: creator,
			Updater: updater,
		}))
	}
	return responses
}

