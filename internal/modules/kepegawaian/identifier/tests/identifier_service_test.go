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
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"neosim_go/internal/modules/kepegawaian/identifier/services"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/factories"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/mocks"

	identifierContracts "neosim_go/internal/modules/kepegawaian/identifier/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
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

// KepegawaianIdentifierServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type KepegawaianIdentifierServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianIdentifierRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      identifierContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianIdentifierServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianIdentifierRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{}
	s.svc = services.NewKepegawaianIdentifierService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
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

func (s *KepegawaianIdentifierServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianIdentifierServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_Superadmin_Success() {
	req := &dto.CreateKepegawaianIdentifierRequest{Name: "Test KepegawaianIdentifier"}
	actor := superadminActor()

	s.repo.On("CreateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_WithPermission_Success() {
	req := &dto.CreateKepegawaianIdentifierRequest{Name: "Test KepegawaianIdentifier"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianIdentifierRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_Forbidden() {
	req := &dto.CreateKepegawaianIdentifierRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateIdentifier(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_RepoError() {
	req := &dto.CreateKepegawaianIdentifierRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateIdentifier(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianIdentifierFactory().Make()
	item.ID = 1

	s.repo.On("GetIdentifierByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetIdentifierByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianIdentifierFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetIdentifierByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetIdentifierByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetIdentifierByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetIdentifierByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetIdentifierByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().Make(),
		*factories.NewKepegawaianIdentifierFactory().Make(),
	}

	s.repo.On("ListIdentifier", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListIdentifier(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}
	items := []models.KepegawaianIdentifier{*factories.NewKepegawaianIdentifierFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListIdentifier", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListIdentifier(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListIdentifier(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}

	s.repo.On("ListIdentifier", 1, 10, filter).Return([]models.KepegawaianIdentifier{}, int64(0), nil)

	result, total, err := s.svc.ListIdentifier(0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}

	s.repo.On("ListIdentifier", 1, 10, filter).Return([]models.KepegawaianIdentifier{}, int64(0), nil)

	_, _, err := s.svc.ListIdentifier(1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListIdentifier", 1, 10, filter)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{Name: "test"}
	items := []models.KepegawaianIdentifier{*factories.NewKepegawaianIdentifierFactory().Make()}

	s.repo.On("ListIdentifier", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListIdentifier(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianIdentifierRequest{Name: &newName}

	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.UpdateIdentifier(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianIdentifierRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.UpdateIdentifier(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianIdentifierRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateIdentifier(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianIdentifierRequest{}

	s.repo.On("GetIdentifierByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateIdentifier(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateKepegawaianIdentifierRequest{Description: &newDesc}

	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateIdentifier", mock.MatchedBy(func(m *models.KepegawaianIdentifier) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateIdentifier(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianIdentifierRequest{}

	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateIdentifier", mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateIdentifier(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1

	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteIdentifier", int64(1)).Return(nil)

	err := s.svc.DeleteIdentifier(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteIdentifier", int64(1)).Return(nil)

	err := s.svc.DeleteIdentifier(1, actor)

	s.NoError(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteIdentifier(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_NotFound() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteIdentifier(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().Make()
	existing.ID = 1

	s.repo.On("GetIdentifierByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteIdentifier", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteIdentifier(1, actor)

	s.Error(err)
}
