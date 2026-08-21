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
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	"neosim_go/internal/modules/kepegawaian/kontak/services"
	"neosim_go/internal/modules/kepegawaian/kontak/tests/factories"
	"neosim_go/internal/modules/kepegawaian/kontak/tests/mocks"

	kontakContracts "neosim_go/internal/modules/kepegawaian/kontak/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianKontak Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/kontak")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/kontak")
	}

	os.Exit(code)
}

// KepegawaianKontakServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini (lihat mis. tag_service_test.go) — karena hanya ada satu
// struct service/repository, satu suite ini sudah cukup untuk semuanya.
type KepegawaianKontakServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianKontakRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      kontakContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianKontakServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianKontakRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewKepegawaianKontakService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetTipeByCode", mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetTipeByLabel", mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetKontakByPegawaiID", mock.Anything, mock.Anything, mock.Anything).Return(nil, int64(0), nil).Maybe()
	s.repo.On("ExistsByNilaiAndTipe", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Maybe()
	s.repo.On("UnsetPrimaryByPegawaiIDAndTipe", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
}

func TestKepegawaianKontakService(t *testing.T) {
	suite.Run(t, new(KepegawaianKontakServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianKontakServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianKontakServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateKontak_Superadmin_Success() {
	req := &dto.CreateKepegawaianKontakRequest{PegawaiID: 1, TipeID: 1, Nilai: "Test KepegawaianKontak"}
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(&models.Tipe{ID: req.TipeID}, nil)
	s.repo.On("CreateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(nil)

	result, err := s.svc.CreateKontak(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Nilai, result.Nilai)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateKontak_WithPermission_Success() {
	req := &dto.CreateKepegawaianKontakRequest{PegawaiID: 1, TipeID: 1, Nilai: "Test KepegawaianKontak"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(&models.Tipe{ID: req.TipeID}, nil)
	s.repo.On("CreateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(nil)

	result, err := s.svc.CreateKontak(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateKontak_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianKontakRequest{PegawaiID: 1, TipeID: 1, Nilai: "Test"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(&models.Tipe{ID: req.TipeID}, nil)
	s.repo.On("CreateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(nil)

	result, err := s.svc.CreateKontak(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateKontak_Forbidden() {
	req := &dto.CreateKepegawaianKontakRequest{PegawaiID: 1, TipeID: 1, Nilai: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKontak(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateKontak_RepoError() {
	req := &dto.CreateKepegawaianKontakRequest{PegawaiID: 1, TipeID: 1, Nilai: "Test"}
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(&models.Tipe{ID: req.TipeID}, nil)
	s.repo.On("CreateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKontak(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetKontakByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianKontakFactory().Make()
	item.ID = 1

	s.repo.On("GetKontakByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKontakByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Nilai, result.Nilai)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetKontakByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianKontakFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetKontakByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKontakByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetKontakByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetKontakByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetKontakByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKontakByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetKontakByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetKontakByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetKontakByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetKontakByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKontakRequest{}
	items := []models.KepegawaianKontak{
		*factories.NewKepegawaianKontakFactory().Make(),
		*factories.NewKepegawaianKontakFactory().Make(),
	}

	s.repo.On("ListKontak", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListKontak(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianKontakRequest{}
	items := []models.KepegawaianKontak{*factories.NewKepegawaianKontakFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListKontak", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKontak(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianKontakRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListKontak(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKontakRequest{}

	s.repo.On("ListKontak", 1, 10, filter).Return([]models.KepegawaianKontak{}, int64(0), nil)

	result, total, err := s.svc.ListKontak(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKontakRequest{}

	s.repo.On("ListKontak", 1, 10, filter).Return([]models.KepegawaianKontak{}, int64(0), nil)

	_, _, err := s.svc.ListKontak(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListKontak", 1, 10, filter)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListKontak_WithNameFilter() {
	actor := superadminActor()
	filterValue := "test"
	filter := &dto.FilterKepegawaianKontakRequest{Nilai: &filterValue}
	items := []models.KepegawaianKontak{*factories.NewKepegawaianKontakFactory().Make()}

	s.repo.On("ListKontak", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKontak(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianKontakRequest{Nilai: &newName}

	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(nil)

	result, err := s.svc.UpdateKontak(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Nilai)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianKontakRequest{Nilai: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(nil)

	result, err := s.svc.UpdateKontak(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianKontakRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKontak(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianKontakRequest{}

	s.repo.On("GetKontakByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKontak(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1
	originalName := existing.Nilai
	newDesc := "New description"
	req := &dto.UpdateKepegawaianKontakRequest{Description: &newDesc}

	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKontak", mock.MatchedBy(func(m *models.KepegawaianKontak) bool {
		return m.Nilai == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateKontak(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Nilai)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateKontak_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianKontakRequest{}

	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKontak", mock.AnythingOfType("*models.KepegawaianKontak")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKontak(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteKontak_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1

	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKontak", int64(1)).Return(nil)

	err := s.svc.DeleteKontak(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteKontak_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKontak", int64(1)).Return(nil)

	err := s.svc.DeleteKontak(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteKontak_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKontak(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteKontak_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKontakByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKontak(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteKontak_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKontakFactory().Make()
	existing.ID = 1

	s.repo.On("GetKontakByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKontak", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKontak(context.Background(), 1, actor)

	s.Error(err)
}
