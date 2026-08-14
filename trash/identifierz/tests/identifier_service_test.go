package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/config"
	identifierContracts "neosim_go/internal/modules/kepegawaian/identifier/contracts"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"neosim_go/internal/modules/kepegawaian/identifier/services"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/factories"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/mocks"
	userModels "neosim_go/internal/modules/users/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianIdentifier Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/identifier")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/identifier")
	}

	os.Exit(code)
}

type KepegawaianIdentifierServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	repo     *mocks.KepegawaianIdentifierRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	cfg      *config.Config
	svc      identifierContracts.Service
}

func (s *KepegawaianIdentifierServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = new(mocks.KepegawaianIdentifierRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 100,
	}

	s.svc = services.NewKepegawaianIdentifierService(
		s.repo,
		s.rbacRepo,
		s.authRepo,
		s.userRepo,
		s.cfg,
	)
}

func TestKepegawaianIdentifierService(t *testing.T) {
	suite.Run(t, new(KepegawaianIdentifierServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianIdentifierServiceTestSuite) mockUserCreator(userID int64) {
	s.userRepo.On("GetByID", userID).Return(&userModels.User{
		ID:       userID,
		Username: "admin",
		Name:     "Administrator",
	}, nil).Maybe()
}

func (s *KepegawaianIdentifierServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil).Maybe()
}

// ── Test Create ───────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_Create_Superadmin_Success() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		Tipe:      models.IdentifierNIK,
		Nilai:     "3171012345678901",
		IsPrimary: true,
		IsAktif:   true,
	}

	s.mockUserCreator(actor.UserID)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, req.Tipe, req.Nilai, int64(0)).Return(false, nil)
	s.repo.On("UnsetPrimaryByPegawaiDAndTipe", s.ctx, req.PegawaiID, req.Tipe, actor.UserID).Return(nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	res, err := s.svc.Create(s.ctx, req, actor)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Nilai, res.Nilai)
	s.Equal(string(req.Tipe), string(res.Tipe)) // ✅ Cast kedua nilai ke string
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_Create_DuplicateNilai_Conflict() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		Tipe:      models.IdentifierNIK,
		Nilai:     "3171012345678901",
	}

	s.repo.On("ExistsByNilaiAndTipe", s.ctx, req.Tipe, req.Nilai, int64(0)).Return(true, nil)

	res, err := s.svc.Create(s.ctx, req, actor)

	s.Nil(res)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_Create_Forbidden() {
	actor := regularActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		Tipe:      models.IdentifierNIK,
		Nilai:     "3171012345678901",
	}
	s.mockNoPermissions()

	res, err := s.svc.Create(s.ctx, req, actor)

	s.Nil(res)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── Test GetByID ──────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetByID_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianIdentifierFactory().Make()
	item.ID = 1

	s.mockUserCreator(actor.UserID)
	s.repo.On("FindByID", s.ctx, int64(1)).Return(item, nil)

	res, err := s.svc.GetByID(s.ctx, 1, actor)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(item.ID, res.ID)
	s.Equal(item.Nilai, res.Nilai)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetByID_NotFound() {
	actor := superadminActor()

	s.repo.On("FindByID", s.ctx, int64(999)).Return(nil, nil)

	res, err := s.svc.GetByID(s.ctx, 999, actor)

	s.Nil(res)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusNotFound, appErr.Code)
}

// ── Test List & ListByPegawai ─────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_List_Success() {
	actor := superadminActor()
	filter := dto.FilterKepegawaianIdentifierRequest{}
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().Make(),
		*factories.NewKepegawaianIdentifierFactory().Make(),
	}

	s.repo.On("FindAll", s.ctx, filter, 1, 10).Return(items, int64(2), nil)

	res, total, err := s.svc.List(s.ctx, 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(res, 2)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListByPegawai_Success() {
	actor := superadminActor()
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().Make(),
	}

	s.repo.On("FindByPegawaiD", s.ctx, int64(5)).Return(items, nil)

	res, err := s.svc.ListByPegawai(s.ctx, 5, actor)

	s.NoError(err)
	s.Len(res, 1)
}

// ── Test Update ───────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_Update_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	existing.IsPrimary = false

	newNilai := "9988776655"
	isPrimary := true
	req := &dto.UpdateKepegawaianIdentifierRequest{
		Nilai:     &newNilai,
		IsPrimary: &isPrimary,
	}

	s.mockUserCreator(actor.UserID)
	s.repo.On("FindByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, existing.Tipe, newNilai, int64(1)).Return(false, nil)
	s.repo.On("UnsetPrimaryByPegawaiDAndTipe", s.ctx, existing.PegawaiID, existing.Tipe, actor.UserID).Return(nil)
	s.repo.On("Update", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	res, err := s.svc.Update(s.ctx, 1, req, actor)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(newNilai, res.Nilai)
	s.True(res.IsPrimary)
}

// ── Test Delete ───────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_Delete_PrimaryWithOthers_UnprocessableEntity() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	existing.IsPrimary = true

	otherActive := factories.NewKepegawaianIdentifierFactory().Make()
	otherActive.ID = 2
	otherActive.IsAktif = true

	s.repo.On("FindByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("FindByPegawaiDAndTipe", s.ctx, existing.PegawaiID, existing.Tipe).
		Return([]models.KepegawaianIdentifier{*existing, *otherActive}, nil)

	err := s.svc.Delete(s.ctx, 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_Delete_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	existing.IsPrimary = false

	s.repo.On("FindByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("Delete", s.ctx, int64(1), actor.UserID).Return(nil)

	err := s.svc.Delete(s.ctx, 1, actor)

	s.NoError(err)
}

// ── Test Expiring & Dropdown Meta ─────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetExpiringSoon_Success() {
	actor := superadminActor()
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().Make(),
	}

	s.repo.On("FindExpiringSoon", s.ctx, 30).Return(items, nil)

	res, err := s.svc.GetExpiringSoon(s.ctx, 30, actor)

	s.NoError(err)
	s.Len(res, 1)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierTypes_Success() {
	types := s.svc.GetIdentifierTypes()
	s.NotEmpty(types)
}
