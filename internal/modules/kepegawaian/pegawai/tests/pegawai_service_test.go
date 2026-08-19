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
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/services"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/mocks"

	pegawaiContracts "neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianPegawai Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/pegawai")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/pegawai")
	}

	os.Exit(code)
}

// KepegawaianPegawaiServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type KepegawaianPegawaiServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianPegawaiRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      pegawaiContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianPegawaiServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.KepegawaianPegawaiRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg      = &config.Config{}
	s.svc = services.NewKepegawaianPegawaiService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
}

func TestKepegawaianPegawaiService(t *testing.T) {
	suite.Run(t, new(KepegawaianPegawaiServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianPegawaiServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianPegawaiServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_Superadmin_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test KepegawaianPegawai"}
	actor := superadminActor()

	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_WithPermission_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test KepegawaianPegawai"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_Forbidden() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePegawai(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_RepoError() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePegawai(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPegawaiByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPegawaiByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetPegawaiByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetPegawaiByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetPegawaiByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	items := []models.KepegawaianPegawai{
		*factories.NewKepegawaianPegawaiFactory().Make(),
		*factories.NewKepegawaianPegawaiFactory().Make(),
	}

	s.repo.On("ListPegawai", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPegawai(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	items := []models.KepegawaianPegawai{*factories.NewKepegawaianPegawaiFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListPegawai", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListPegawai(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListPegawai(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}

	s.repo.On("ListPegawai", 1, 10, filter).Return([]models.KepegawaianPegawai{}, int64(0), nil)

	result, total, err := s.svc.ListPegawai(0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}

	s.repo.On("ListPegawai", 1, 10, filter).Return([]models.KepegawaianPegawai{}, int64(0), nil)

	_, _, err := s.svc.ListPegawai(1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPegawai", 1, 10, filter)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{Name: "test"}
	items := []models.KepegawaianPegawai{*factories.NewKepegawaianPegawaiFactory().Make()}

	s.repo.On("ListPegawai", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListPegawai(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianPegawaiRequest{Name: &newName}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.UpdatePegawai(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianPegawaiRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.UpdatePegawai(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePegawai(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianPegawaiRequest{}

	s.repo.On("GetPegawaiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePegawai(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateKepegawaianPegawaiRequest{Description: &newDesc}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.MatchedBy(func(m *models.KepegawaianPegawai) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdatePegawai(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianPegawaiRequest{}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdatePegawai(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1)).Return(nil)

	err := s.svc.DeletePegawai(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1)).Return(nil)

	err := s.svc.DeletePegawai(1, actor)

	s.NoError(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePegawai(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(999)).Return(nil, nil)

	err := s.svc.DeletePegawai(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeletePegawai(1, actor)

	s.Error(err)
}
