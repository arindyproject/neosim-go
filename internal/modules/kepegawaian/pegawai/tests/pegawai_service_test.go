package tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

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

type KepegawaianPegawaiServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianPegawaiRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	svc      pegawaiContracts.Service
}

func (s *KepegawaianPegawaiServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.KepegawaianPegawaiRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.svc = services.NewKepegawaianPegawaiService(s.repo, s.rbacRepo, s.authRepo)
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

func (s *KepegawaianPegawaiServiceTestSuite) Test_Create_Superadmin_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test KepegawaianPegawai"}
	actor := superadminActor()

	s.repo.On("Create", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Create_WithPermission_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test KepegawaianPegawai"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("Create", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Create_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("Create", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Create_Forbidden() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Create_RepoError() {
	req := &dto.CreateKepegawaianPegawaiRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("Create", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.repo.On("GetByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	items := []models.KepegawaianPegawai{
		*factories.NewKepegawaianPegawaiFactory().Make(),
		*factories.NewKepegawaianPegawaiFactory().Make(),
	}

	s.repo.On("List", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	items := []models.KepegawaianPegawai{*factories.NewKepegawaianPegawaiFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("List", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}

	s.repo.On("List", 1, 10, filter).Return([]models.KepegawaianPegawai{}, int64(0), nil)

	result, total, err := s.svc.List(0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}

	s.repo.On("List", 1, 10, filter).Return([]models.KepegawaianPegawai{}, int64(0), nil)

	_, _, err := s.svc.List(1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "List", 1, 10, filter)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_List_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{Name: "test"}
	items := []models.KepegawaianPegawai{*factories.NewKepegawaianPegawaiFactory().Make()}

	s.repo.On("List", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianPegawaiRequest{Name: &newName}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Update", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianPegawaiRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Update", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianPegawaiRequest{}

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Update(999, req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateKepegawaianPegawaiRequest{Description: &newDesc}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Update", mock.MatchedBy(func(m *models.KepegawaianPegawai) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Update_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianPegawaiRequest{}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Update", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Delete_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete", int64(1)).Return(nil)

	err := s.svc.Delete(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Delete_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete", int64(1)).Return(nil)

	err := s.svc.Delete(1, actor)

	s.NoError(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Delete_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.Delete(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Delete_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	err := s.svc.Delete(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_Delete_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.Delete(1, actor)

	s.Error(err)
}
