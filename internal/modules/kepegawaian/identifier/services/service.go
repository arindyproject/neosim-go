package services

import (
	identifierContracts "neosim_go/internal/modules/kepegawaian/identifier/contracts"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	userContracts "neosim_go/internal/modules/users/contracts"
)

type service struct {
	repo     identifierContracts.Repository
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
	userRepo userContracts.Repository
}

// NewKepegawaianIdentifierService membuat instance service baru
func NewKepegawaianIdentifierService(
	repo identifierContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
) identifierContracts.Service {
	return &service{
		repo:     repo,
		rbacRepo: rbacRepo,
		authRepo: authRepo,
		userRepo: userRepo,
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

func (s *service) buildAuditMaps(items []models.KepegawaianIdentifier) (map[int64]*he.UserData, map[int64]*he.UserData) {
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

func toIdentifierResponses(
	items []models.KepegawaianIdentifier,
	creatorsMap, updatersMap map[int64]*he.UserData,
) []dto.KepegawaianIdentifierResponse {
	responses := make([]dto.KepegawaianIdentifierResponse, 0, len(items))
	for _, item := range items {
		var creator, updater *he.UserData
		if item.CreatedBy != nil {
			creator = creatorsMap[*item.CreatedBy]
		}
		if item.UpdatedBy != nil {
			updater = updatersMap[*item.UpdatedBy]
		}
		responses = append(responses, *dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
			Identifier: &item,
			Creator:    creator,
			Updater:    updater,
		}))
	}
	return responses
}
