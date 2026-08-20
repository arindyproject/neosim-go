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
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/modules/artikel/kategori/services"
	"neosim_go/internal/modules/artikel/kategori/tests/factories"
	"neosim_go/internal/modules/artikel/kategori/tests/mocks"

	kategoriContracts "neosim_go/internal/modules/artikel/kategori/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  ArtikelKategori Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/artikel/kategori")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/artikel/kategori")
	}

	os.Exit(code)
}

// ArtikelKategoriServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type ArtikelKategoriServiceTestSuite struct {
	suite.Suite
	repo     *mocks.ArtikelKategoriRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      kategoriContracts.Service
	cfg      *config.Config
}

func (s *ArtikelKategoriServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.ArtikelKategoriRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewArtikelKategoriService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()
}

func TestArtikelKategoriService(t *testing.T) {
	suite.Run(t, new(ArtikelKategoriServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *ArtikelKategoriServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *ArtikelKategoriServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateKategori_Superadmin_Success() {
	req := &dto.CreateArtikelKategoriRequest{Name: "Test ArtikelKategori"}
	actor := superadminActor()

	s.repo.On("CreateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(nil)

	result, err := s.svc.CreateKategori(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateKategori_WithPermission_Success() {
	req := &dto.CreateArtikelKategoriRequest{Name: "Test ArtikelKategori"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("CreateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(nil)

	result, err := s.svc.CreateKategori(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateKategori_WithManagePermission_Success() {
	req := &dto.CreateArtikelKategoriRequest{Name: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("CreateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(nil)

	result, err := s.svc.CreateKategori(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateKategori_Forbidden() {
	req := &dto.CreateArtikelKategoriRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateKategori_RepoError() {
	req := &dto.CreateArtikelKategoriRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetKategoriByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewArtikelKategoriFactory().Make()
	item.ID = 1

	s.repo.On("GetKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKategoriByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetKategoriByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewArtikelKategoriFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKategoriByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetKategoriByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetKategoriByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetKategoriByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetKategoriByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetKategoriByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetKategoriByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetKategoriByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterArtikelKategoriRequest{}
	items := []models.ArtikelKategori{
		*factories.NewArtikelKategoriFactory().Make(),
		*factories.NewArtikelKategoriFactory().Make(),
	}

	s.repo.On("ListKategori", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterArtikelKategoriRequest{}
	items := []models.ArtikelKategori{*factories.NewArtikelKategoriFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListKategori", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterArtikelKategoriRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListKategori(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterArtikelKategoriRequest{}

	s.repo.On("ListKategori", 1, 10, filter).Return([]models.ArtikelKategori{}, int64(0), nil)

	result, total, err := s.svc.ListKategori(context.Background(),0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterArtikelKategoriRequest{}

	s.repo.On("ListKategori", 1, 10, filter).Return([]models.ArtikelKategori{}, int64(0), nil)

	_, _, err := s.svc.ListKategori(context.Background(),1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListKategori", 1, 10, filter)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListKategori_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterArtikelKategoriRequest{Name: "test"}
	items := []models.ArtikelKategori{*factories.NewArtikelKategoriFactory().Make()}

	s.repo.On("ListKategori", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateArtikelKategoriRequest{Name: &newName}

	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(nil)

	result, err := s.svc.UpdateKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateArtikelKategoriRequest{Name: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(nil)

	result, err := s.svc.UpdateKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateArtikelKategoriRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKategori(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateArtikelKategoriRequest{}

	s.repo.On("GetKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKategori(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_PartialFields() {
	actor := superadminActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"
	req := &dto.UpdateArtikelKategoriRequest{Description: &newDesc}

	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKategori", mock.MatchedBy(func(m *models.ArtikelKategori) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateKategori_RepoError() {
	actor := superadminActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1
	req := &dto.UpdateArtikelKategoriRequest{}

	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKategori", mock.AnythingOfType("*models.ArtikelKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKategori(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteKategori_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKategori", int64(1)).Return(nil)

	err := s.svc.DeleteKategori(context.Background(),1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteKategori_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKategori", int64(1)).Return(nil)

	err := s.svc.DeleteKategori(context.Background(),1, actor)

	s.NoError(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKategori(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteKategori_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKategoriByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKategori(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteKategori_RepoError() {
	actor := superadminActor()
	existing := factories.NewArtikelKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKategori", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKategori(context.Background(),1, actor)

	s.Error(err)
}
