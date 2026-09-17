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
	pegawaiModels "neosim_go/internal/modules/kepegawaian/pegawai/models"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/shared/cache"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"

	masterAlamatMock "neosim_go/internal/modules/master/alamat/tests/mocks"
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
// sub-module ini — karena hanya ada satu struct service/repository, satu suite
// ini sudah cukup untuk semuanya.
type KepegawaianAlamatServiceTestSuite struct {
	suite.Suite
	repo             *mocks.KepegawaianAlamatRepositoryMock
	rbacRepo         *mocks.RBACRepositoryMock
	authRepo         *mocks.AuthRepositoryMock
	userRepo         *mocks.UserRepositoryMock
	pegawaiRepo      *mocks.KepegawaianPegawaiRepositoryMock
	masterAlamatRepo *masterAlamatMock.MasterAlamatRepositoryMock
	svc              alamatContracts.Service
	cfg              *config.Config
}

func (s *KepegawaianAlamatServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianAlamatRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.pegawaiRepo = new(mocks.KepegawaianPegawaiRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	cacheManager := cache.NewManager(nil, false, 0)
	s.masterAlamatRepo = new(masterAlamatMock.MasterAlamatRepositoryMock)
	s.svc = services.NewKepegawaianAlamatService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.pegawaiRepo, s.masterAlamatRepo, s.cfg, cacheManager)

	// Stub default agar buildCreator/buildAuditMaps tidak panic.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()

	// Pegawai selalu dianggap ada (ID 10), kecuali skenario butuh lain.
	s.pegawaiRepo.On("GetPegawaiByID", mock.Anything, mock.Anything).
		Return(&pegawaiModels.KepegawaianPegawai{ID: 10}, nil).Maybe()

	// buildWilayahSingle / buildWilayahMaps hanya proses pengayaan response —
	// tidak ada test di bawah ini yang meng-assert nama wilayah, jadi cukup
	// distub kosong (nil) agar tidak panic, apa pun ID yang dikirim.
	s.masterAlamatRepo.On("GetByIDProvinsi", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetByIDKotaKabupaten", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetByIDKecamatan", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetByIDKelurahanDesa", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetByIDNegara", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	// buildWilayahSingle / buildWilayahMaps hanya proses pengayaan response —
	// distub kosong (nil) agar tidak panic, apa pun ID yang dikirim.
	s.masterAlamatRepo.On("GetSimpelByIDNegara", mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetSimpelByIDProvinsi", mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetSimpelByIDKotaKabupaten", mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetSimpelByIDKecamatan", mock.Anything).Return(nil, nil).Maybe()
	s.masterAlamatRepo.On("GetSimpelByIDKelurahanDesa", mock.Anything).Return(nil, nil).Maybe()
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

// validCreateReq mengembalikan request Create minimal yang valid (tanpa
// hirarki wilayah, tanpa primary), dipakai sebagai basis di banyak test.
func validCreateReq() *dto.CreateKepegawaianAlamatRequest {
	return &dto.CreateKepegawaianAlamatRequest{
		PegawaiID: 10,
		TipeID:    1,
		Jalan:     "Jl. Merdeka No. 1",
	}
}

func tipeDomisili() *models.Tipe {
	return &models.Tipe{ID: 1, Code: "domisili", Label: "Domisili"}
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_Superadmin_Success() {
	req := validCreateReq()
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Jalan, result.Jalan)
	s.Equal(req.PegawaiID, result.PegawaiID)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_WithPermission_Success() {
	req := validCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_WithManagePermission_Success() {
	req := validCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_Forbidden() {
	req := validCreateReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_TipeNotFound() {
	req := validCreateReq()
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(nil, nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_DuplicateAlamat() {
	req := validCreateReq()
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(true, nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateAlamat", mock.Anything)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_SetPrimary_UnsetsOldPrimary() {
	req := validCreateReq()
	req.IsPrimary = true
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UnsetPrimaryAlamatByPegawaiID", req.PegawaiID, actor.UserID).Return(nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.True(result.IsPrimary)
	s.repo.AssertCalled(s.T(), "UnsetPrimaryAlamatByPegawaiID", req.PegawaiID, actor.UserID)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_NotPrimary_DoesNotUnset() {
	req := validCreateReq() // IsPrimary default false
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	_, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.NoError(err)
	s.repo.AssertNotCalled(s.T(), "UnsetPrimaryAlamatByPegawaiID", mock.Anything, mock.Anything)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_CreateAlamat_RepoError() {
	req := validCreateReq()
	actor := superadminActor()

	s.repo.On("GetTipeByID", req.TipeID).Return(tipeDomisili(), nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("CreateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateAlamat(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianAlamatFactory().Make()
	item.ID = 1

	s.repo.On("GetAlamatByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetAlamatByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Jalan, result.Jalan)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianAlamatFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetAlamatByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetAlamatByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetAlamatByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_GetAlamatByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetAlamatByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}
	items := []models.KepegawaianAlamat{
		*factories.NewKepegawaianAlamatFactory().Make(),
		*factories.NewKepegawaianAlamatFactory().Make(),
	}

	s.repo.On("ListAlamat", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListAlamat(context.Background(), 1, 10, filter, actor)

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

	result, total, err := s.svc.ListAlamat(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListAlamat(context.Background(), 1, 10, filter, actor)

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

	result, total, err := s.svc.ListAlamat(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianAlamatRequest{}

	s.repo.On("ListAlamat", 1, 10, filter).Return([]models.KepegawaianAlamat{}, int64(0), nil)

	_, _, err := s.svc.ListAlamat(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListAlamat", 1, 10, filter)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_ListAlamat_WithJalanFilter() {
	actor := superadminActor()
	jalan := "merdeka"
	filter := &dto.FilterKepegawaianAlamatRequest{Jalan: &jalan}
	items := []models.KepegawaianAlamat{*factories.NewKepegawaianAlamatFactory().Make()}

	s.repo.On("ListAlamat", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListAlamat(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_Superadmin_Success() {
	actor := superadminActor()
	existing := alamatFactoryNoWilayah() // ← bukan factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	newJalan := "Jl. Updated No. 2"
	req := &dto.UpdateKepegawaianAlamatRequest{Jalan: &newJalan}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Require().NoError(err)   // ← Require, bukan s.NoError biasa
	s.Require().NotNil(result) // ← Require, bukan s.NotNil biasa
	s.Equal(newJalan, result.Jalan)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_WithPermission_Success() {
	actor := regularActor()
	existing := alamatFactoryNoWilayah()
	existing.ID = 1
	newJalan := "Jl. Updated"
	req := &dto.UpdateKepegawaianAlamatRequest{Jalan: &newJalan}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Require().NoError(err)
	s.Require().NotNil(result)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianAlamatRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

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

	result, err := s.svc.UpdateAlamat(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_PartialFields() {
	actor := superadminActor()
	existing := alamatFactoryNoWilayah()
	existing.ID = 1
	originalJalan := existing.Jalan
	newDesc := "Deskripsi baru"
	req := &dto.UpdateKepegawaianAlamatRequest{Description: &newDesc}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.MatchedBy(func(m *models.KepegawaianAlamat) bool {
		return m.Jalan == originalJalan && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Equal(originalJalan, result.Jalan)
	s.Equal(newDesc, *result.Description)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_SetPrimary_UnsetsOldPrimary() {
	actor := superadminActor()
	existing := alamatFactoryNoWilayah()
	existing.ID = 1
	existing.PegawaiID = 10
	existing.IsPrimary = false
	isPrimary := true
	req := &dto.UpdateKepegawaianAlamatRequest{IsPrimary: &isPrimary}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("UnsetPrimaryAlamatByPegawaiID", existing.PegawaiID, actor.UserID).Return(nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.MatchedBy(func(m *models.KepegawaianAlamat) bool {
		return m.IsPrimary == true
	})).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.True(result.IsPrimary)
	s.repo.AssertCalled(s.T(), "UnsetPrimaryAlamatByPegawaiID", existing.PegawaiID, actor.UserID)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_IsPrimaryNil_DoesNotChange() {
	actor := superadminActor()
	existing := alamatFactoryNoWilayah()
	existing.ID = 1
	existing.IsPrimary = true
	req := &dto.UpdateKepegawaianAlamatRequest{} // IsPrimary tidak dikirim

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.NoError(err)
	s.True(result.IsPrimary) // tetap true, tidak berubah
	s.repo.AssertNotCalled(s.T(), "UnsetPrimaryAlamatByPegawaiID", mock.Anything, mock.Anything)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_DuplicateAlamat() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.NegaraID = nil
	existing.ProvinsiID = nil
	existing.KotaKabupatenID = nil
	existing.KecamatanID = nil
	existing.KelurahanDesaID = nil

	newJalan := "Jl. Duplikat"
	req := &dto.UpdateKepegawaianAlamatRequest{Jalan: &newJalan}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(true, nil)

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateAlamat", mock.Anything)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_UpdateAlamat_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianAlamatRequest{}

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckDuplicateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat"), mock.Anything).Return(false, nil)
	s.repo.On("UpdateAlamat", mock.AnythingOfType("*models.KepegawaianAlamat")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateAlamat(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.IsPrimary = false

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.IsPrimary = false

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_NotFound() {
	actor := superadminActor()

	s.repo.On("GetAlamatByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteAlamat(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.IsPrimary = false

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.Error(err)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Primary_BlockedWhenOthersExist() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.PegawaiID = 10
	existing.IsPrimary = true

	other := factories.NewKepegawaianAlamatFactory().Make()
	other.ID = 2
	other.PegawaiID = 10

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("FindAlamatByPegawaiID", existing.PegawaiID).
		Return([]models.KepegawaianAlamat{*existing, *other}, nil)

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "DeleteAlamat", mock.Anything, mock.Anything)
}

func (s *KepegawaianAlamatServiceTestSuite) Test_DeleteAlamat_Primary_AllowedWhenOnlyAddress() {
	actor := superadminActor()
	existing := factories.NewKepegawaianAlamatFactory().Make()
	existing.ID = 1
	existing.PegawaiID = 10
	existing.IsPrimary = true

	s.repo.On("GetAlamatByID", int64(1)).Return(existing, nil)
	s.repo.On("FindAlamatByPegawaiID", existing.PegawaiID).
		Return([]models.KepegawaianAlamat{*existing}, nil)
	s.repo.On("DeleteAlamat", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteAlamat(context.Background(), 1, actor)

	s.NoError(err)
}

// alamatFactoryNoWilayah mengembalikan alamat dari factory dengan seluruh
// field hirarki wilayah di-nil-kan, supaya test yang tidak sedang menguji
// validasi hirarki tidak tersandung olehnya secara tidak sengaja.
func alamatFactoryNoWilayah() *models.KepegawaianAlamat {
	m := factories.NewKepegawaianAlamatFactory().Make()
	m.NegaraID = nil
	m.ProvinsiID = nil
	m.KotaKabupatenID = nil
	m.KecamatanID = nil
	m.KelurahanDesaID = nil
	return m
}
