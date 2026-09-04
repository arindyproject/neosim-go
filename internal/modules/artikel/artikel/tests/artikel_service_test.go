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
	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	"neosim_go/internal/modules/artikel/artikel/services"
	"neosim_go/internal/modules/artikel/artikel/tests/factories"
	"neosim_go/internal/modules/artikel/artikel/tests/mocks"

	artikelContracts "neosim_go/internal/modules/artikel/artikel/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  Artikel Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/artikel/artikel")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/artikel/artikel")
	}

	os.Exit(code)
}

// ArtikelServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type ArtikelServiceTestSuite struct {
	suite.Suite
	repo     *mocks.ArtikelRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      artikelContracts.Service
	cfg      *config.Config
}

func (s *ArtikelServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.ArtikelRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewArtikelService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()
}

func TestArtikelService(t *testing.T) {
	suite.Run(t, new(ArtikelServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *ArtikelServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *ArtikelServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *ArtikelServiceTestSuite) Test_CreateArtikel_Superadmin_Success() {
	req := &dto.CreateArtikelRequest{Name: "Test Artikel"}
	actor := superadminActor()

	s.repo.On("CreateArtikel", mock.AnythingOfType("*models.Artikel")).Return(nil)

	result, err := s.svc.CreateArtikel(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelServiceTestSuite) Test_CreateArtikel_WithPermission_Success() {
	req := &dto.CreateArtikelRequest{Name: "Test Artikel"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreateArtikel", mock.AnythingOfType("*models.Artikel")).Return(nil)

	result, err := s.svc.CreateArtikel(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelServiceTestSuite) Test_CreateArtikel_WithManagePermission_Success() {
	req := &dto.CreateArtikelRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreateArtikel", mock.AnythingOfType("*models.Artikel")).Return(nil)

	result, err := s.svc.CreateArtikel(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelServiceTestSuite) Test_CreateArtikel_Forbidden() {
	req := &dto.CreateArtikelRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateArtikel(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelServiceTestSuite) Test_CreateArtikel_RepoError() {
	req := &dto.CreateArtikelRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateArtikel", mock.AnythingOfType("*models.Artikel")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateArtikel(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelServiceTestSuite) Test_GetArtikelByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewArtikelFactory().Make()
	item.ID = 1

	s.repo.On("GetArtikelByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetArtikelByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *ArtikelServiceTestSuite) Test_GetArtikelByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewArtikelFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetArtikelByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetArtikelByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelServiceTestSuite) Test_GetArtikelByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetArtikelByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelServiceTestSuite) Test_GetArtikelByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetArtikelByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetArtikelByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelServiceTestSuite) Test_GetArtikelByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetArtikelByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetArtikelByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterArtikelRequest{}
	items := []models.Artikel{
		*factories.NewArtikelFactory().Make(),
		*factories.NewArtikelFactory().Make(),
	}

	s.repo.On("ListArtikel", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListArtikel(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterArtikelRequest{}
	items := []models.Artikel{*factories.NewArtikelFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListArtikel", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListArtikel(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterArtikelRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListArtikel(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterArtikelRequest{}

	s.repo.On("ListArtikel", 1, 10, filter).Return([]models.Artikel{}, int64(0), nil)

	result, total, err := s.svc.ListArtikel(context.Background(),0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterArtikelRequest{}

	s.repo.On("ListArtikel", 1, 10, filter).Return([]models.Artikel{}, int64(0), nil)

	_, _, err := s.svc.ListArtikel(context.Background(),1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListArtikel", 1, 10, filter)
}

func (s *ArtikelServiceTestSuite) Test_ListArtikel_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterArtikelRequest{Name: "test"}
	items := []models.Artikel{*factories.NewArtikelFactory().Make()}

	s.repo.On("ListArtikel", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListArtikel(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateArtikelRequest{Name: &newName}

	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateArtikel", mock.AnythingOfType("*models.Artikel")).Return(nil)

	result, err := s.svc.UpdateArtikel(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateArtikelRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateArtikel", mock.AnythingOfType("*models.Artikel")).Return(nil)

	result, err := s.svc.UpdateArtikel(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateArtikelRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateArtikel(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateArtikelRequest{}

	s.repo.On("GetArtikelByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateArtikel(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_PartialFields() {
	actor := superadminActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateArtikelRequest{Description: &newDesc}

	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateArtikel", mock.MatchedBy(func(m *models.Artikel) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateArtikel(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *ArtikelServiceTestSuite) Test_UpdateArtikel_RepoError() {
	actor := superadminActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1
	req := &dto.UpdateArtikelRequest{}

	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateArtikel", mock.AnythingOfType("*models.Artikel")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateArtikel(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelServiceTestSuite) Test_DeleteArtikel_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1

	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteArtikel", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteArtikel(context.Background(),1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelServiceTestSuite) Test_DeleteArtikel_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteArtikel", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteArtikel(context.Background(),1, actor)

	s.NoError(err)
}

func (s *ArtikelServiceTestSuite) Test_DeleteArtikel_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteArtikel(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelServiceTestSuite) Test_DeleteArtikel_NotFound() {
	actor := superadminActor()

	s.repo.On("GetArtikelByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteArtikel(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelServiceTestSuite) Test_DeleteArtikel_RepoError() {
	actor := superadminActor()
	existing := factories.NewArtikelFactory().Make()
	existing.ID = 1

	s.repo.On("GetArtikelByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteArtikel", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteArtikel(context.Background(),1, actor)

	s.Error(err)
}
