package tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/config"
	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
	"neosim_go/internal/modules/master/departemen/services"
	"neosim_go/internal/modules/master/departemen/tests/factories"
	"neosim_go/internal/modules/master/departemen/tests/mocks"

	departemenContracts "neosim_go/internal/modules/master/departemen/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  MasterDepartemen Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/master/departemen")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/master/departemen")
	}

	os.Exit(code)
}

// MasterDepartemenServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type MasterDepartemenServiceTestSuite struct {
	suite.Suite
	repo     *mocks.MasterDepartemenRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      departemenContracts.Service
	cfg      *config.Config
}

func (s *MasterDepartemenServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.MasterDepartemenRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg      = &config.Config{}
	s.svc = services.NewMasterDepartemenService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
}

func TestMasterDepartemenService(t *testing.T) {
	suite.Run(t, new(MasterDepartemenServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *MasterDepartemenServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *MasterDepartemenServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *MasterDepartemenServiceTestSuite) Test_CreateDepartemen_Superadmin_Success() {
	req := &dto.CreateMasterDepartemenRequest{Name: "Test MasterDepartemen"}
	actor := superadminActor()

	s.repo.On("CreateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(nil)

	result, err := s.svc.CreateDepartemen(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterDepartemenServiceTestSuite) Test_CreateDepartemen_WithPermission_Success() {
	req := &dto.CreateMasterDepartemenRequest{Name: "Test MasterDepartemen"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(nil)

	result, err := s.svc.CreateDepartemen(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterDepartemenServiceTestSuite) Test_CreateDepartemen_WithManagePermission_Success() {
	req := &dto.CreateMasterDepartemenRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(nil)

	result, err := s.svc.CreateDepartemen(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterDepartemenServiceTestSuite) Test_CreateDepartemen_Forbidden() {
	req := &dto.CreateMasterDepartemenRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateDepartemen(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterDepartemenServiceTestSuite) Test_CreateDepartemen_RepoError() {
	req := &dto.CreateMasterDepartemenRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateDepartemen(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterDepartemenServiceTestSuite) Test_GetDepartemenByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewMasterDepartemenFactory().Make()
	item.ID = 1

	s.repo.On("GetDepartemenByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetDepartemenByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterDepartemenServiceTestSuite) Test_GetDepartemenByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewMasterDepartemenFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetDepartemenByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetDepartemenByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterDepartemenServiceTestSuite) Test_GetDepartemenByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetDepartemenByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterDepartemenServiceTestSuite) Test_GetDepartemenByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetDepartemenByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetDepartemenByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterDepartemenServiceTestSuite) Test_GetDepartemenByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetDepartemenByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetDepartemenByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterMasterDepartemenRequest{}
	items := []models.MasterDepartemen{
		*factories.NewMasterDepartemenFactory().Make(),
		*factories.NewMasterDepartemenFactory().Make(),
	}

	s.repo.On("ListDepartemen", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListDepartemen(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterMasterDepartemenRequest{}
	items := []models.MasterDepartemen{*factories.NewMasterDepartemenFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListDepartemen", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListDepartemen(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterMasterDepartemenRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListDepartemen(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterMasterDepartemenRequest{}

	s.repo.On("ListDepartemen", 1, 10, filter).Return([]models.MasterDepartemen{}, int64(0), nil)

	result, total, err := s.svc.ListDepartemen(0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterMasterDepartemenRequest{}

	s.repo.On("ListDepartemen", 1, 10, filter).Return([]models.MasterDepartemen{}, int64(0), nil)

	_, _, err := s.svc.ListDepartemen(1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListDepartemen", 1, 10, filter)
}

func (s *MasterDepartemenServiceTestSuite) Test_ListDepartemen_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterMasterDepartemenRequest{Name: "test"}
	items := []models.MasterDepartemen{*factories.NewMasterDepartemenFactory().Make()}

	s.repo.On("ListDepartemen", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListDepartemen(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterDepartemenRequest{Name: &newName}

	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(nil)

	result, err := s.svc.UpdateDepartemen(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateMasterDepartemenRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(nil)

	result, err := s.svc.UpdateDepartemen(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterDepartemenRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateDepartemen(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterDepartemenRequest{}

	s.repo.On("GetDepartemenByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateDepartemen(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_PartialFields() {
	actor := superadminActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateMasterDepartemenRequest{Description: &newDesc}

	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateDepartemen", mock.MatchedBy(func(m *models.MasterDepartemen) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateDepartemen(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *MasterDepartemenServiceTestSuite) Test_UpdateDepartemen_RepoError() {
	actor := superadminActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterDepartemenRequest{}

	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateDepartemen", mock.AnythingOfType("*models.MasterDepartemen")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateDepartemen(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterDepartemenServiceTestSuite) Test_DeleteDepartemen_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1

	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteDepartemen", int64(1)).Return(nil)

	err := s.svc.DeleteDepartemen(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterDepartemenServiceTestSuite) Test_DeleteDepartemen_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteDepartemen", int64(1)).Return(nil)

	err := s.svc.DeleteDepartemen(1, actor)

	s.NoError(err)
}

func (s *MasterDepartemenServiceTestSuite) Test_DeleteDepartemen_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteDepartemen(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterDepartemenServiceTestSuite) Test_DeleteDepartemen_NotFound() {
	actor := superadminActor()

	s.repo.On("GetDepartemenByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteDepartemen(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterDepartemenServiceTestSuite) Test_DeleteDepartemen_RepoError() {
	actor := superadminActor()
	existing := factories.NewMasterDepartemenFactory().Make()
	existing.ID = 1

	s.repo.On("GetDepartemenByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteDepartemen", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteDepartemen(1, actor)

	s.Error(err)
}
