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
	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	"neosim_go/internal/modules/kepegawaian/alamat/services"
	"neosim_go/internal/modules/kepegawaian/alamat/tests/factories"
	"neosim_go/internal/modules/kepegawaian/alamat/tests/mocks"

	alamatContracts "neosim_go/internal/modules/kepegawaian/alamat/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianAlamat Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/alamat")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/alamat")
	}

	os.Exit(code)
}

// KepegawaianAlamatServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type KepegawaianAlamatServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianAlamatRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      alamatContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianAlamatServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.KepegawaianAlamatRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewKepegawaianAlamatService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()
}

func TestKepegawaianAlamatService(t *testing.T) {
	suite.Run(t, new(KepegawaianAlamatServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianAlamatServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianAlamatServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_Superadmin_Success() {
	req := &dto.CreateKepegawaianAlamatRequest{Name: "Test KepegawaianAlamat"}
	actor := superadminActor()

	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_WithPermission_Success() {
	req := &dto.CreateKepegawaianAlamatRequest{Name: "Test KepegawaianAlamat"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianAlamatRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_Forbidden() {
	req := &dto.CreateKepegawaianAlamatRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateAlamat(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_RepoError() {
	req := &dto.CreateKepegawaianAlamatRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateAlamat(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianAlamatFactory().Make()
	item.ID = 1

	s.repo.On("GetAlamatByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetAlamatByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianAlamatFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetAlamatByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetAlamatByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetAlamatByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetAlamatByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}
	items := []models.KepegawaianAlamat{
		*factories.NewKepegawaianAlamatFactory().Make(),
		*factories.NewKepegawaianAlamatFactory().Make(),
	}

	s.repo.On("ListAlamat", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListAlamat(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}
	items := []models.KepegawaianAlamat{*factories.NewKepegawaianAlamatFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListAlamat", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListAlamat(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListAlamat(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}

	s.repo.On("ListAlamat", 1, 10, filter).Return([]models.KepegawaianAlamat{}, int64(0), nil)

	result, total, err := s.svc.ListAlamat(context.Background(),0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}

	s.repo.On("ListAlamat", 1, 10, filter).Return([]models.KepegawaianAlamat{}, int64(0), nil)

	_, _, err := s.svc.ListAlamat(context.Background(),1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListAlamat", 1, 10, filter)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{Name: "test"}
	items := []models.KepegawaianAlamat{*factories.NewKepegawaianAlamatFactory().Make()}

	s.repo.On("ListAlamat", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListAlamat(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianAlamatRequest{Name: &newName}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianAlamatRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianAlamatRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateAlamat(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianAlamatRequest{}

	s.repo.On("GetAlamatByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateAlamat(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateKepegawaianAlamatRequest{Description: &newDesc}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateAlamat", mock.MatchedBy(func(m *models.KepegawaianAlamat) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianAlamatRequest{}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateAlamat(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteAlamat(context.Background(),1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteAlamat(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteAlamat(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_NotFound() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteAlamat(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteAlamat(context.Background(),1, actor)

	s.Error(err)
}
