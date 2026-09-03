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
	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/services"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/tests/factories"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/tests/mocks"

	kualifikasiContracts "neosim_go/internal/modules/kepegawaian/kualifikasi/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianKualifikasi Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/kualifikasi")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/kualifikasi")
	}

	os.Exit(code)
}

// KepegawaianKualifikasiServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini — karena hanya ada satu struct service/repository, satu suite
// ini sudah cukup untuk semuanya.
type KepegawaianKualifikasiServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianKualifikasiRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      kualifikasiContracts.Service
	cfg      *config.Config
}

func (s *KepegawaianKualifikasiServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianKualifikasiRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.svc = services.NewKepegawaianKualifikasiService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()
}

func TestKepegawaianKualifikasiService(t *testing.T) {
	suite.Run(t, new(KepegawaianKualifikasiServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *KepegawaianKualifikasiServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

// testTipe adalah master tipe kualifikasi dummy yang dikembalikan GetTipeByID.
func testTipe(id int64) *models.Tipe {
	return &models.Tipe{
		ID:    id,
		Code:  fmt.Sprintf("TIPE-%d", id),
		Label: fmt.Sprintf("Tipe %d", id),
	}
}

// ── Create ────────────────────────────────────────────────────────────────

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_Superadmin_Success() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID:     10,
		TipeID:        1,
		Nama:          "Sertifikasi BLS",
		Penyelenggara: "PMI",
	}
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(testTipe(req.TipeID), nil)
	s.repo.On("CreateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Nama, result.Nama)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_WithPermission_Success() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID:     10,
		TipeID:        1,
		Nama:          "Sertifikasi BLS",
		Penyelenggara: "PMI",
	}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(testTipe(req.TipeID), nil)
	s.repo.On("CreateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_WithManagePermission_Success() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID:     10,
		TipeID:        1,
		Nama:          "Sertifikasi BLS",
		Penyelenggara: "PMI",
	}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(testTipe(req.TipeID), nil)
	s.repo.On("CreateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_Forbidden() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID: 10, TipeID: 1, Nama: "Test", Penyelenggara: "Test",
	}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_TipeNotFound() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID: 10, TipeID: 999, Nama: "Test", Penyelenggara: "Test",
	}
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(nil, nil)

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_CreateKualifikasi_RepoError() {
	req := &dto.CreateKepegawaianKualifikasiRequest{
		PegawaiID: 10, TipeID: 1, Nama: "Test", Penyelenggara: "Test",
	}
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(testTipe(req.TipeID), nil)
	s.repo.On("CreateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKualifikasi(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ───────────────────────────────────────────────────────────────

func (s *KepegawaianKualifikasiServiceTestSuite) Test_GetKualifikasiByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianKualifikasiFactory().Make()
	item.ID = 1

	s.repo.On("GetKualifikasiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKualifikasiByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Nama, result.Nama)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_GetKualifikasiByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianKualifikasiFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetKualifikasiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetKualifikasiByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_GetKualifikasiByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetKualifikasiByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_GetKualifikasiByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKualifikasiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetKualifikasiByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_GetKualifikasiByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetKualifikasiByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetKualifikasiByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

// ── List ──────────────────────────────────────────────────────────────────

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{}
	items := []models.KepegawaianKualifikasi{
		*factories.NewKepegawaianKualifikasiFactory().Make(),
		*factories.NewKepegawaianKualifikasiFactory().Make(),
	}

	s.repo.On("ListKualifikasi", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListKualifikasi(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{}
	items := []models.KepegawaianKualifikasi{*factories.NewKepegawaianKualifikasiFactory().Make()}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListKualifikasi", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKualifikasi(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListKualifikasi(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{}

	s.repo.On("ListKualifikasi", 1, 10, filter).Return([]models.KepegawaianKualifikasi{}, int64(0), nil)

	result, total, err := s.svc.ListKualifikasi(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{}

	s.repo.On("ListKualifikasi", 1, 10, filter).Return([]models.KepegawaianKualifikasi{}, int64(0), nil)

	_, _, err := s.svc.ListKualifikasi(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListKualifikasi", 1, 10, filter)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_ListKualifikasi_WithNamaFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianKualifikasiRequest{Nama: "BLS"}
	items := []models.KepegawaianKualifikasi{*factories.NewKepegawaianKualifikasiFactory().Make()}

	s.repo.On("ListKualifikasi", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKualifikasi(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

// ── Update ────────────────────────────────────────────────────────────────

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1
	newNama := "Sertifikasi ACLS"
	req := &dto.UpdateKepegawaianKualifikasiRequest{Nama: &newNama}

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNomorSertifikatAndTipe", existing.TipeID, *existing.NomorSertifikat, int64(1)).Return(false, nil)
	s.repo.On("UpdateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newNama, result.Nama)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1
	newNama := "Updated"
	req := &dto.UpdateKepegawaianKualifikasiRequest{Nama: &newNama}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNomorSertifikatAndTipe", existing.TipeID, *existing.NomorSertifikat, int64(1)).Return(false, nil)
	s.repo.On("UpdateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianKualifikasiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianKualifikasiRequest{}

	s.repo.On("GetKualifikasiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_DuplicateNomorSertifikat() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1
	newNomor := "CERT-999999"
	req := &dto.UpdateKepegawaianKualifikasiRequest{NomorSertifikat: &newNomor}

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNomorSertifikatAndTipe", existing.TipeID, newNomor, int64(1)).Return(true, nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1
	originalNama := existing.Nama
	newPenyelenggara := "Penyelenggara Baru"
	req := &dto.UpdateKepegawaianKualifikasiRequest{Penyelenggara: &newPenyelenggara}

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNomorSertifikatAndTipe", existing.TipeID, *existing.NomorSertifikat, int64(1)).Return(false, nil)
	s.repo.On("UpdateKualifikasi", mock.MatchedBy(func(m *models.KepegawaianKualifikasi) bool {
		return m.Nama == originalNama && m.Penyelenggara == newPenyelenggara
	})).Return(nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalNama, result.Nama)
	s.Equal(newPenyelenggara, result.Penyelenggara)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_NoNomorSertifikat_SkipsDuplicateCheck() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().With("nomor_sertifikat", "").Make()
	existing.ID = 1
	existing.NomorSertifikat = nil
	newNama := "Tanpa Nomor Sertifikat"
	req := &dto.UpdateKepegawaianKualifikasiRequest{Nama: &newNama}

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(nil)

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	// Tidak boleh memanggil ExistsByNomorSertifikatAndTipe sama sekali
	s.repo.AssertNotCalled(s.T(), "ExistsByNomorSertifikatAndTipe", mock.Anything, mock.Anything, mock.Anything)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_UpdateKualifikasi_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianKualifikasiRequest{}

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNomorSertifikatAndTipe", existing.TipeID, *existing.NomorSertifikat, int64(1)).Return(false, nil)
	s.repo.On("UpdateKualifikasi", mock.AnythingOfType("*models.KepegawaianKualifikasi")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKualifikasi(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ────────────────────────────────────────────────────────────────

func (s *KepegawaianKualifikasiServiceTestSuite) Test_DeleteKualifikasi_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKualifikasi", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteKualifikasi(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_DeleteKualifikasi_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKualifikasi", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteKualifikasi(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_DeleteKualifikasi_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKualifikasi(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_DeleteKualifikasi_NotFound() {
	actor := superadminActor()

	s.repo.On("GetKualifikasiByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKualifikasi(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKualifikasiServiceTestSuite) Test_DeleteKualifikasi_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianKualifikasiFactory().Make()
	existing.ID = 1

	s.repo.On("GetKualifikasiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKualifikasi", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKualifikasi(context.Background(), 1, actor)

	s.Error(err)
}
