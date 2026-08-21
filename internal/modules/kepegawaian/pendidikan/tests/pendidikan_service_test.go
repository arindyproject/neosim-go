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
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	"neosim_go/internal/modules/kepegawaian/pendidikan/services"
	"neosim_go/internal/modules/kepegawaian/pendidikan/tests/factories"
	"neosim_go/internal/modules/kepegawaian/pendidikan/tests/mocks"

	pendidikanContracts "neosim_go/internal/modules/kepegawaian/pendidikan/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianPendidikan Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/pendidikan")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/pendidikan")
	}

	os.Exit(code)
}

// KepegawaianPendidikanServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type KepegawaianPendidikanServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianPendidikanRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      pendidikanContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianPendidikanServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianPendidikanRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewKepegawaianPendidikanService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()
}

func TestKepegawaianPendidikanService(t *testing.T) {
	suite.Run(t, new(KepegawaianPendidikanServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianPendidikanServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianPendidikanServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreatePendidikan_Superadmin_Success() {
	req := &dto.CreateKepegawaianPendidikanRequest{NamaInstitusi: "Test KepegawaianPendidikan", PegawaiID: 1, JenjangID: 1}
	actor := superadminActor()

	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.NamaInstitusi, result.NamaInstitusi)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreatePendidikan_WithPermission_Success() {
	req := &dto.CreateKepegawaianPendidikanRequest{NamaInstitusi: "Test KepegawaianPendidikan", PegawaiID: 1, JenjangID: 1}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreatePendidikan_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianPendidikanRequest{NamaInstitusi: "Test", PegawaiID: 1, JenjangID: 1}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreatePendidikan_Forbidden() {
	req := &dto.CreateKepegawaianPendidikanRequest{NamaInstitusi: "Test", PegawaiID: 1, JenjangID: 1}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreatePendidikan_RepoError() {
	req := &dto.CreateKepegawaianPendidikanRequest{NamaInstitusi: "Test", PegawaiID: 1, JenjangID: 1}
	actor := superadminActor()

	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetPendidikanByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianPendidikanFactory().Make()
	item.ID = 1

	s.repo.On("GetPendidikanByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPendidikanByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.NamaInstitusi, result.NamaInstitusi)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetPendidikanByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianPendidikanFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetPendidikanByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPendidikanByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetPendidikanByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetPendidikanByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetPendidikanByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPendidikanByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetPendidikanByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetPendidikanByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetPendidikanByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetPendidikanByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{}
	items := []models.KepegawaianPendidikan{
		*factories.NewKepegawaianPendidikanFactory().Make(),
		*factories.NewKepegawaianPendidikanFactory().Make(),
	}

	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{}
	items := []models.KepegawaianPendidikan{*factories.NewKepegawaianPendidikanFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{}

	s.repo.On("ListPendidikan", 1, 10, filter).Return([]models.KepegawaianPendidikan{}, int64(0), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{}

	s.repo.On("ListPendidikan", 1, 10, filter).Return([]models.KepegawaianPendidikan{}, int64(0), nil)

	_, _, err := s.svc.ListPendidikan(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPendidikan", 1, 10, filter)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListPendidikan_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPendidikanRequest{NamaInstitusi: "test"}
	items := []models.KepegawaianPendidikan{*factories.NewKepegawaianPendidikanFactory().Make()}

	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianPendidikanRequest{NamaInstitusi: &newName}

	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.NamaInstitusi)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianPendidikanRequest{NamaInstitusi: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianPendidikanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianPendidikanRequest{}

	s.repo.On("GetPendidikanByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1
	originalName := existing.NamaInstitusi
	newBidangStudi := "New study"
	req := &dto.UpdateKepegawaianPendidikanRequest{BidangStudi: &newBidangStudi}

	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePendidikan", mock.MatchedBy(func(m *models.KepegawaianPendidikan) bool {
		return m.NamaInstitusi == originalName && *m.BidangStudi == newBidangStudi
	})).Return(nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.NamaInstitusi)
	s.Equal(newBidangStudi, *result.BidangStudi)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdatePendidikan_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianPendidikanRequest{}

	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePendidikan", mock.AnythingOfType("*models.KepegawaianPendidikan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeletePendidikan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1

	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePendidikan", int64(1)).Return(nil)

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeletePendidikan_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePendidikan", int64(1)).Return(nil)

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeletePendidikan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeletePendidikan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPendidikanByID", int64(999)).Return(nil, nil)

	err := s.svc.DeletePendidikan(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeletePendidikan_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPendidikanFactory().Make()
	existing.ID = 1

	s.repo.On("GetPendidikanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePendidikan", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.Error(err)
}
