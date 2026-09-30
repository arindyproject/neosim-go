package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/services"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/mocks"
	masterModels "neosim_go/internal/modules/master/master/models"
	userModels "neosim_go/internal/modules/users/models"

	pegawaiContracts "neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/shared/cache"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"

	masterMock "neosim_go/internal/modules/master/master/tests/mocks"
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
	repo       *mocks.KepegawaianPegawaiRepositoryMock
	rbacRepo   *mocks.RBACRepositoryMock
	authRepo   *mocks.AuthRepositoryMock
	userRepo   *mocks.UserRepositoryMock
	masterRepo *masterMock.MasterRepositoryMock
	svc        pegawaiContracts.Service
	cfg        *config.Config
}

func (s *KepegawaianPegawaiServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianPegawaiRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}
	s.masterRepo = new(masterMock.MasterRepositoryMock)
	cacheManager := cache.NewManager(nil, false, 0)
	s.svc = services.NewKepegawaianPegawaiService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.masterRepo, s.cfg, cacheManager)

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
	// Boleh dipanggil 0 kali atau lebih (.Maybe()) tergantung skenario test.
	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return([]userModels.User{}, nil).Maybe()
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

// ── Helper ────────────────────────────────────────────────────────────────────

func dateOnly(year int, month time.Month, day int) *types.DateOnly {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return types.NewDateOnlyPtr(&t)
}

func newCreateReq() *dto.CreateKepegawaianPegawaiRequest {
	golDarah := int64(1)
	return &dto.CreateKepegawaianPegawaiRequest{
		NIK:                "3577010101900001",
		NomorPegawai:       "PEG-000001",
		NamaLengkap:        "Test Pegawai",
		JenisKelaminID:     1,
		TanggalLahir:       dateOnly(1990, 1, 1),
		TempatLahir:        "Madiun",
		GolonganDarahID:    &golDarah,
		AgamaID:            1,
		StatusPernikahanID: 1,
		TanggalMasuk:       dateOnly(2020, 1, 1),
		JenisID:            1,
		StatusID:           1,
	}
}

// mockCreateDeps men-stub semua dependensi jalur sukses CreatePegawai:
// cek duplikat NIK/nomor, master kepegawaian, dan master umum.
func (s *KepegawaianPegawaiServiceTestSuite) mockCreateDeps() {
	s.repo.On("GetPegawaiByNIK", mock.Anything).Return(nil, nil)
	s.repo.On("GetPegawaiByNomorPegawai", mock.Anything).Return(nil, nil)
	s.repo.On("GetJenisByID", int64(1)).Return(&models.Jenis{ID: 1}, nil)
	s.repo.On("GetStatusByID", int64(1)).Return(&models.Status{ID: 1}, nil)
	s.masterRepo.On("GetByIDJenisKelamin", int64(1)).Return(&masterModels.MasterJenisKelamin{ID: 1}, nil)
	s.masterRepo.On("GetByIDGolonganDarah", int64(1)).Return(&masterModels.MasterGolonganDarah{ID: 1}, nil)
	s.masterRepo.On("GetByIDAgama", int64(1)).Return(&masterModels.MasterAgama{ID: 1}, nil)
	s.masterRepo.On("GetByIDStatusPernikahan", int64(1)).Return(&masterModels.MasterStatusPernikahan{ID: 1}, nil)
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_Superadmin_Success() {
	req := newCreateReq()
	actor := superadminActor()

	s.mockCreateDeps()
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.NamaLengkap, result.NamaLengkap)
	s.Equal(req.NIK, result.NIK)
	s.True(result.IsAktif)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_WithPermission_Success() {
	req := newCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.mockCreateDeps()
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_WithManagePermission_Success() {
	req := newCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.mockCreateDeps()
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_RepoError() {
	actor := superadminActor()

	s.mockCreateDeps()
	s.repo.On("CreatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_DuplicateNIK() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByNIK", mock.Anything).Return(factories.NewKepegawaianPegawaiFactory().Make(), nil)

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreatePegawai", mock.Anything)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_DuplicateNomorPegawai() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByNIK", mock.Anything).Return(nil, nil)
	s.repo.On("GetPegawaiByNomorPegawai", mock.Anything).Return(factories.NewKepegawaianPegawaiFactory().Make(), nil)

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_UserAlreadyLinked() {
	actor := superadminActor()
	req := newCreateReq()
	userID := int64(5)
	req.UserID = &userID

	s.userRepo.ExpectedCalls = nil // ganti stub default agar user "ditemukan"
	s.userRepo.On("GetByID", mock.Anything).Return(&userModels.User{ID: userID}, nil)
	s.userRepo.On("GetByIDs", mock.Anything).Return([]userModels.User{}, nil).Maybe()
	s.repo.On("GetPegawaiByUserID", userID).Return(factories.NewKepegawaianPegawaiFactory().Make(), nil)

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_UserNotFound() {
	actor := superadminActor()
	req := newCreateReq()
	userID := int64(999)
	req.UserID = &userID // stub default GetByID mengembalikan nil, nil

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_TanggalKeluarBeforeMasuk() {
	actor := superadminActor()
	req := newCreateReq()
	req.TanggalKeluar = dateOnly(2019, 12, 31) // masuk 2020-01-01

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusBadRequest, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_TanggalKeluarForcesInactive() {
	actor := superadminActor()
	req := newCreateReq()
	req.TanggalKeluar = dateOnly(2024, 1, 1)
	aktif := true
	req.IsAktif = &aktif // dikirim true, tetapi harus dipaksa false

	s.mockCreateDeps()
	s.repo.On("CreatePegawai", mock.MatchedBy(func(m *models.KepegawaianPegawai) bool {
		return !m.IsAktif && m.TanggalKeluar != nil
	})).Return(nil)

	result, err := s.svc.CreatePegawai(context.Background(), req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_JenisNotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByNIK", mock.Anything).Return(nil, nil)
	s.repo.On("GetPegawaiByNomorPegawai", mock.Anything).Return(nil, nil)
	s.repo.On("GetJenisByID", int64(1)).Return(nil, nil)

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Jenis pegawai tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreatePegawai_AgamaNotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByNIK", mock.Anything).Return(nil, nil)
	s.repo.On("GetPegawaiByNomorPegawai", mock.Anything).Return(nil, nil)
	s.repo.On("GetJenisByID", int64(1)).Return(&models.Jenis{ID: 1}, nil)
	s.repo.On("GetStatusByID", int64(1)).Return(&models.Status{ID: 1}, nil)
	s.masterRepo.On("GetByIDJenisKelamin", int64(1)).Return(&masterModels.MasterJenisKelamin{ID: 1}, nil)
	s.masterRepo.On("GetByIDGolonganDarah", int64(1)).Return(&masterModels.MasterGolonganDarah{ID: 1}, nil)
	s.masterRepo.On("GetByIDAgama", int64(1)).Return(nil, nil)

	result, err := s.svc.CreatePegawai(context.Background(), newCreateReq(), actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Agama tidak ditemukan")
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPegawaiByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.NamaLengkap, result.NamaLengkap)
	s.Equal(item.NIK, result.NIK)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianPegawaiFactory().Make()
	item.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPegawaiByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetPegawaiByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetPegawaiByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetPegawaiByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetPegawaiByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	items := []models.KepegawaianPegawai{
		*factories.NewKepegawaianPegawaiFactory().Make(),
		*factories.NewKepegawaianPegawaiFactory().Make(),
	}

	s.repo.On("ListPegawai", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPegawai(context.Background(), 1, 10, filter, actor)

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

	result, total, err := s.svc.ListPegawai(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListPegawai(context.Background(), 1, 10, filter, actor)

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

	result, total, err := s.svc.ListPegawai(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{}

	s.repo.On("ListPegawai", 1, 10, filter).Return([]models.KepegawaianPegawai{}, int64(0), nil)

	_, _, err := s.svc.ListPegawai(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPegawai", 1, 10, filter)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListPegawai_WithNameFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianPegawaiRequest{Name: "test"}
	items := []models.KepegawaianPegawai{*factories.NewKepegawaianPegawaiFactory().Make()}

	s.repo.On("ListPegawai", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListPegawai(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateKepegawaianPegawaiRequest{NamaLengkap: &newName}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.NamaLengkap)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newName := "Updated"
	req := &dto.UpdateKepegawaianPegawaiRequest{NamaLengkap: &newName}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianPegawaiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

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

	result, err := s.svc.UpdatePegawai(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	originalName := existing.NamaLengkap
	newTempat := "Ngawi"
	req := &dto.UpdateKepegawaianPegawaiRequest{TempatLahir: &newTempat}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.MatchedBy(func(m *models.KepegawaianPegawai) bool {
		return m.NamaLengkap == originalName && m.TempatLahir == newTempat
	})).Return(nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.NamaLengkap)
	s.Equal(newTempat, result.TempatLahir)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_DuplicateNIK() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	newNIK := "3577999999999999"
	req := &dto.UpdateKepegawaianPegawaiRequest{NIK: &newNIK}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("GetPegawaiByNIK", newNIK).Return(factories.NewKepegawaianPegawaiFactory().Make(), nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdatePegawai", mock.Anything)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_SameNIK_SkipDuplicateCheck() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	sameNIK := existing.NIK
	req := &dto.UpdateKepegawaianPegawaiRequest{NIK: &sameNIK}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(nil)

	_, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.NoError(err)
	s.repo.AssertNotCalled(s.T(), "GetPegawaiByNIK", mock.Anything)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_TanggalKeluar_SetsInactive() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	existing.IsAktif = true
	existing.TanggalMasuk = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianPegawaiRequest{TanggalKeluar: dateOnly(2024, 6, 30)}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.MatchedBy(func(m *models.KepegawaianPegawai) bool {
		return m.TanggalKeluar != nil && !m.IsAktif
	})).Return(nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
	s.NotNil(result.TanggalKeluar)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_TanggalKeluarBeforeMasuk() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	existing.TanggalMasuk = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianPegawaiRequest{TanggalKeluar: dateOnly(2019, 1, 1)}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusBadRequest, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_TanggalKeluar_WithActiveTrue_Rejected() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	existing.TanggalMasuk = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	aktif := true
	req := &dto.UpdateKepegawaianPegawaiRequest{TanggalKeluar: dateOnly(2024, 1, 1), IsAktif: &aktif}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.Nil(result)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusBadRequest, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_JenisNotFound() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	existing.JenisID, existing.StatusID = 1, 1
	newJenis := int64(99)
	req := &dto.UpdateKepegawaianPegawaiRequest{JenisID: &newJenis}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("GetJenisByID", int64(99)).Return(nil, nil)

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Jenis pegawai tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdatePegawai_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKepegawaianPegawaiRequest{}

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePegawai", mock.AnythingOfType("*models.KepegawaianPegawai")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdatePegawai(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ────────────────────────────────────────────────────────────────────
// (tidak berubah dari versi sebelumnya)

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeletePegawai(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_WithPermission_Success() {
	actor := regularActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeletePegawai(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePegawai(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPegawaiByID", int64(999)).Return(nil, nil)

	err := s.svc.DeletePegawai(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeletePegawai_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianPegawaiFactory().Make()
	existing.ID = 1

	s.repo.On("GetPegawaiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePegawai", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeletePegawai(context.Background(), 1, actor)

	s.Error(err)
}
