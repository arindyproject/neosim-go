package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/services"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/mocks"

	jabatanContracts "neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	pegawaiModels "neosim_go/internal/modules/kepegawaian/pegawai/models"
	pegawaiMock "neosim_go/internal/modules/kepegawaian/pegawai/tests/mocks"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/shared/cache"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"

	masterMocks "neosim_go/internal/modules/master/departemen/tests/mocks"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  KepegawaianJabatan Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/kepegawaian/jabatan")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/kepegawaian/jabatan")
	}

	os.Exit(code)
}

// KepegawaianJabatanServiceTestSuite dipakai bersama oleh SELURUH item di dalam
// sub-module ini — karena hanya ada satu struct service/repository,
// satu suite ini sudah cukup untuk semuanya.
type KepegawaianJabatanServiceTestSuite struct {
	suite.Suite
	repo                 *mocks.KepegawaianJabatanRepositoryMock
	rbacRepo             *mocks.RBACRepositoryMock
	authRepo             *mocks.AuthRepositoryMock
	userRepo             *mocks.UserRepositoryMock
	masterDepartemenRepo *masterMocks.MasterDepartemenRepositoryMock
	pegawaiRepo          *pegawaiMock.KepegawaianPegawaiRepositoryMock
	svc                  jabatanContracts.Service
	cfg                  *config.Config
}

func (s *KepegawaianJabatanServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianJabatanRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.masterDepartemenRepo = new(masterMocks.MasterDepartemenRepositoryMock)
	s.pegawaiRepo = new(pegawaiMock.KepegawaianPegawaiRepositoryMock)
	cacheManager := cache.NewManager(nil, false, 0)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 10,
	}

	s.svc = services.NewKepegawaianJabatanService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.masterDepartemenRepo, s.pegawaiRepo, s.cfg, cacheManager)

	s.pegawaiRepo.On("GetPegawaiByID", mock.Anything, mock.Anything).
		Return(&pegawaiModels.KepegawaianPegawai{ID: 10}, nil).Maybe()

	s.userRepo.On("GetByID", mock.Anything).Return(nil, nil).Maybe()
	s.userRepo.On("GetByIDs", mock.Anything).Return(nil, nil).Maybe()

	s.repo.On("GetJobTitleKategoriByCode", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetJobTitleKategoriByLabel", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	s.repo.On("GetPositionKategoriByCode", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetPositionKategoriByLabel", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	s.repo.On("GetJobTitleRumpunProfesiByCode", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	s.repo.On("GetJobTitleRumpunProfesiByLabel", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	s.repo.On("GetSpecializationByLabel", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
}

func TestKepegawaianJabatanService(t *testing.T) {
	suite.Run(t, new(KepegawaianJabatanServiceTestSuite))
}

// ── helpers ──────────────────────────────────────────────────────────────────

const dateLayout = "2006-01-02"

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func ptrTo[T any](v T) *T { return &v }

// dateOnly membuat *types.DateOnly untuk request DTO.
func dateOnly(year int, month time.Month, day int) *types.DateOnly {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return types.NewDateOnlyPtr(&t)
}

// fmtDate memformat *types.DateOnly ke YYYY-MM-DD untuk assertion.
func fmtDate(d *types.DateOnly) string {
	if t := d.ToTimePtr(); t != nil {
		return t.Format(dateLayout)
	}
	return ""
}

func (s *KepegawaianJabatanServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *KepegawaianJabatanServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

// mockJobTitleCodeAvailable men-stub ExistsByCode agar kode dianggap belum
// dipakai. excludeID 0 untuk create, ID yang diedit untuk update.
func (s *KepegawaianJabatanServiceTestSuite) mockJobTitleCodeAvailable(code string, excludeID int64) {
	s.repo.On("ExistsByCode", code, excludeID).Return(false, nil).Once()
}

// newCreateReq membuat request create yang valid tanpa spesialisasi
// (supaya tidak butuh mock GetSpecializationByID) dan bukan jabatan primer.
func newCreateReq() *dto.CreateKepegawaianJabatanRequest {
	return &dto.CreateKepegawaianJabatanRequest{
		PegawaiID:    1,
		DepartmentID: 1,
		PositionID:   1,
		JobTitleID:   1,
		TanggalMulai: *dateOnly(2024, 1, 15),
	}
}

// newJabatan membuat model dari factory dengan kondisi deterministik:
// tanpa spesialisasi, non-primer, aktif, tanpa tanggal_selesai.
func newJabatan(id int64) *models.KepegawaianJabatan {
	m := factories.NewKepegawaianJabatanFactory().
		With("pegawai_id", int64(1)).
		With("is_primary", false).
		With("is_aktif", true).
		With("tanggal_selesai", nil).
		With("specialization_id", nil).
		Make()
	m.ID = id
	return m
}

// mockCreateSetsID membuat mock CreateJabatan mengisi ID (seperti DB) lalu
// mock reload GetJabatanByID mengembalikan model tersebut.
func (s *KepegawaianJabatanServiceTestSuite) mockCreateSetsID(id int64) {
	s.repo.On("CreateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).
		Run(func(args mock.Arguments) {
			args.Get(0).(*models.KepegawaianJabatan).ID = id
		}).Return(nil)
	s.repo.On("GetJabatanByID", id).Return(nil, nil).Maybe() // fallback ke model hasil create
}

func (s *KepegawaianJabatanServiceTestSuite) assertAppErrCode(err error, code int) {
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(code, appErr.Code)
}

// ── Create ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_Superadmin_Success() {
	req := newCreateReq()
	actor := superadminActor()
	s.mockCreateSetsID(1)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.PegawaiID, result.PegawaiID)
	s.Equal(req.DepartmentID, result.DepartmentID)
	s.Equal(req.PositionID, result.PositionID)
	s.Equal(req.JobTitleID, result.JobTitleID)
	s.Require().NotNil(result.TanggalMulai)
	s.Equal("2024-01-15", fmtDate(result.TanggalMulai))
	s.Nil(result.TanggalSelesai)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_WithPermission_Success() {
	req := newCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.mockCreateSetsID(1)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_WithManagePermission_Success() {
	req := newCreateReq()
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyManage).Return(true, nil)
	s.mockCreateSetsID(1)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_Forbidden() {
	req := newCreateReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusForbidden)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_RepoError() {
	req := newCreateReq()
	actor := superadminActor()

	s.repo.On("CreateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_DefaultIsAktifTrue_WhenNoTanggalSelesai() {
	req := newCreateReq()
	actor := superadminActor()

	s.repo.On("CreateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return m.IsAktif && !m.IsPrimary && m.TanggalSelesai == nil
	})).Run(func(args mock.Arguments) {
		args.Get(0).(*models.KepegawaianJabatan).ID = 1
	}).Return(nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(nil, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.True(result.IsAktif)
	s.False(result.IsPrimary)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_DefaultIsAktifFalse_WhenTanggalSelesaiSet() {
	req := newCreateReq()
	req.TanggalSelesai = dateOnly(2024, 6, 30)
	actor := superadminActor()

	s.mockCreateSetsID(1)
	s.repo.On("CreateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return !m.IsAktif && m.TanggalSelesai != nil
	})).Return(nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
	s.Require().NotNil(result.TanggalSelesai)
	s.Equal("2024-06-30", fmtDate(result.TanggalSelesai))
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_ExplicitIsAktifOverridesDefault() {
	req := newCreateReq()
	req.IsAktif = ptrTo(false)
	actor := superadminActor()

	s.repo.On("CreateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return !m.IsAktif
	})).Run(func(args mock.Arguments) {
		args.Get(0).(*models.KepegawaianJabatan).ID = 1
	}).Return(nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(nil, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_TanggalMulaiRequired() {
	req := newCreateReq()
	req.TanggalMulai = types.DateOnly{} // kosong

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_TanggalSelesaiBeforeMulai() {
	req := newCreateReq()
	req.TanggalMulai = *dateOnly(2024, 6, 1)
	req.TanggalSelesai = dateOnly(2024, 1, 1)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_SpecializationNotFound() {
	req := newCreateReq()
	req.SpecializationID = ptrTo(int64(99))

	s.repo.On("GetSpecializationByID", int64(99)).Return(nil, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_SpecializationJobTitleMismatch() {
	req := newCreateReq() // JobTitleID = 1
	req.SpecializationID = ptrTo(int64(5))

	s.repo.On("GetSpecializationByID", int64(5)).Return(&models.Specialization{
		ID: 5, JobTitleID: ptrTo(int64(2)), // milik job title lain
	}, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_SpecializationMatchesJobTitle_Success() {
	req := newCreateReq() // JobTitleID = 1
	req.SpecializationID = ptrTo(int64(5))

	s.repo.On("GetSpecializationByID", int64(5)).Return(&models.Specialization{
		ID: 5, JobTitleID: ptrTo(int64(1)),
	}, nil)
	s.mockCreateSetsID(1)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_SpecializationRepoError() {
	req := newCreateReq()
	req.SpecializationID = ptrTo(int64(5))

	s.repo.On("GetSpecializationByID", int64(5)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_PrimaryConflict() {
	req := newCreateReq()
	req.IsPrimary = ptrTo(true)

	s.repo.On("ExistsPrimaryAktifByPegawai", req.PegawaiID, int64(0)).Return(true, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusConflict)
	s.repo.AssertNotCalled(s.T(), "CreateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_Primary_Success() {
	req := newCreateReq()
	req.IsPrimary = ptrTo(true)

	s.repo.On("ExistsPrimaryAktifByPegawai", req.PegawaiID, int64(0)).Return(false, nil)
	s.mockCreateSetsID(1)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.NoError(err)
	s.True(result.IsPrimary)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_PrimaryInactive_SkipsPrimaryCheck() {
	req := newCreateReq()
	req.IsPrimary = ptrTo(true)
	req.TanggalSelesai = dateOnly(2024, 6, 30) // otomatis is_aktif = false

	s.repo.On("CreateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).
		Run(func(args mock.Arguments) {
			args.Get(0).(*models.KepegawaianJabatan).ID = 1
		}).Return(nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(nil, nil)

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertNotCalled(s.T(), "ExistsPrimaryAktifByPegawai", mock.Anything, mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJabatan_PrimaryCheckRepoError() {
	req := newCreateReq()
	req.IsPrimary = ptrTo(true)

	s.repo.On("ExistsPrimaryAktifByPegawai", req.PegawaiID, int64(0)).Return(false, fmt.Errorf("db error"))

	result, err := s.svc.CreateJabatan(context.Background(), req, superadminActor())

	s.Nil(result)
	s.Error(err)
}

// Format tanggal salah ditolak saat bind JSON (types.DateOnly.UnmarshalJSON),
// bukan lagi di service.
func (s *KepegawaianJabatanServiceTestSuite) Test_CreateRequest_InvalidTanggalFormat_RejectedOnBind() {
	var req dto.CreateKepegawaianJabatanRequest
	err := json.Unmarshal([]byte(`{"pegawai_id":1,"tanggal_mulai":"15-01-2024"}`), &req)
	s.Error(err)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_Superadmin_Success() {
	actor := superadminActor()
	item := newJabatan(1)

	s.repo.On("GetJabatanByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJabatanByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.PegawaiID, result.PegawaiID)
	s.Equal(item.PositionID, result.PositionID)
	s.Equal(item.JobTitleID, result.JobTitleID)
	s.Require().NotNil(result.TanggalMulai)
	s.Equal(item.TanggalMulai.Format(dateLayout), fmtDate(result.TanggalMulai))
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_MapsRelations() {
	actor := superadminActor()
	item := newJabatan(1)
	item.Position = &models.Position{ID: item.PositionID, Name: "Kepala Instalasi"}
	item.JobTitle = &models.JobTitle{ID: item.JobTitleID, Code: "DR-SPESIALIS", Label: "Dokter Spesialis"}
	item.SpecializationID = ptrTo(int64(3))
	item.Specialization = &models.Specialization{ID: 3, Code: "SP-A", Label: "Spesialis Anak", Gelar: ptrTo("Sp.A")}

	s.repo.On("GetJabatanByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJabatanByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result.Position)
	s.Equal("Kepala Instalasi", result.Position.Name)
	s.Require().NotNil(result.JobTitle)
	s.Equal("DR-SPESIALIS", result.JobTitle.Code)
	s.Require().NotNil(result.Specialization)
	s.Equal("SP-A", result.Specialization.Code)
	s.Equal("Sp.A", *result.Specialization.Gelar)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_WithPermission_Success() {
	actor := regularActor()
	item := newJabatan(1)

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJabatanByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJabatanByID(context.Background(), 1, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusForbidden)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJabatanByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJabatanByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetJabatanByID", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetJabatanByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByPegawaiID ───────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_Superadmin_Success() {
	actor := superadminActor()
	items := []models.KepegawaianJabatan{*newJabatan(1), *newJabatan(2)}

	s.repo.On("GetJabatanByPegawaiID", int64(1), 1, 10).Return(items, int64(2), nil)

	result, total, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 1, 10, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_WithPermission_Success() {
	actor := regularActor()
	items := []models.KepegawaianJabatan{*newJabatan(1)}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("GetJabatanByPegawaiID", int64(1), 1, 10).Return(items, int64(1), nil)

	result, total, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 1, 10, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, total, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 1, 10, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.assertAppErrCode(err, http.StatusForbidden)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_DefaultPagination() {
	actor := superadminActor()

	s.repo.On("GetJabatanByPegawaiID", int64(1), 1, 10).Return([]models.KepegawaianJabatan{}, int64(0), nil)

	result, total, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 0, 0, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_PageSizeCapped() {
	actor := superadminActor()

	s.repo.On("GetJabatanByPegawaiID", int64(1), 1, 10).Return([]models.KepegawaianJabatan{}, int64(0), nil)

	_, _, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 1, 999, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "GetJabatanByPegawaiID", int64(1), 1, 10)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJabatanByPegawaiID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetJabatanByPegawaiID", int64(1), 1, 10).Return(nil, int64(0), fmt.Errorf("db error"))

	result, total, err := s.svc.GetJabatanByPegawaiID(context.Background(), 1, 1, 10, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

// ── List ─────────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}
	items := []models.KepegawaianJabatan{*newJabatan(1), *newJabatan(2)}

	s.repo.On("ListJabatan", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJabatan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_WithPermission_Success() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}
	items := []models.KepegawaianJabatan{*newJabatan(1)}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead).Return(true, nil)
	s.repo.On("ListJabatan", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListJabatan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJabatan(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.assertAppErrCode(err, http.StatusForbidden)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_DefaultPagination() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}

	s.repo.On("ListJabatan", 1, 10, filter).Return([]models.KepegawaianJabatan{}, int64(0), nil)

	result, total, err := s.svc.ListJabatan(context.Background(), 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}

	s.repo.On("ListJabatan", 1, 10, filter).Return([]models.KepegawaianJabatan{}, int64(0), nil)

	_, _, err := s.svc.ListJabatan(context.Background(), 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListJabatan", 1, 10, filter)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_WithFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianJabatanRequest{
		PegawaiID: ptrTo(int64(1)),
		IsAktif:   ptrTo(true),
	}
	items := []models.KepegawaianJabatan{*newJabatan(1)}

	s.repo.On("ListJabatan", 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListJabatan(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
	s.repo.AssertCalled(s.T(), "ListJabatan", 1, 10, filter)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJabatan_RepoError() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianJabatanRequest{}

	s.repo.On("ListJabatan", 1, 10, filter).Return(nil, int64(0), fmt.Errorf("db error"))

	result, total, err := s.svc.ListJabatan(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_Superadmin_Success() {
	actor := superadminActor()
	existing := newJabatan(1)
	req := &dto.UpdateKepegawaianJabatanRequest{
		DepartmentID: ptrTo(int64(7)),
		PositionID:   ptrTo(int64(8)),
		NomorSK:      ptrTo("  SK/2025/001  "),
	}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(int64(7), result.DepartmentID)
	s.Equal(int64(8), result.PositionID)
	s.Require().NotNil(result.NomorSK)
	s.Equal("SK/2025/001", *result.NomorSK) // di-trim
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_WithPermission_Success() {
	actor := regularActor()
	existing := newJabatan(1)
	req := &dto.UpdateKepegawaianJabatanRequest{DepartmentID: ptrTo(int64(3))}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyUpdate).Return(true, nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianJabatanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusForbidden)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianJabatanRequest{}

	s.repo.On("GetJabatanByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_PartialFields() {
	actor := superadminActor()
	existing := newJabatan(1)
	origDept, origPos, origJob := existing.DepartmentID, existing.PositionID, existing.JobTitleID
	origMulai := existing.TanggalMulai
	req := &dto.UpdateKepegawaianJabatanRequest{NomorSK: ptrTo("SK/BARU/1")}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return m.DepartmentID == origDept &&
			m.PositionID == origPos &&
			m.JobTitleID == origJob &&
			m.TanggalMulai.Equal(origMulai) &&
			m.NomorSK != nil && *m.NomorSK == "SK/BARU/1"
	})).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(origDept, result.DepartmentID)
	s.Equal("SK/BARU/1", *result.NomorSK)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_SetsUpdatedBy() {
	actor := superadminActor()
	existing := newJabatan(1)
	req := &dto.UpdateKepegawaianJabatanRequest{DepartmentID: ptrTo(int64(2))}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return m.UpdatedBy != nil && *m.UpdatedBy == actor.UserID &&
			time.Since(m.UpdatedAt) < time.Minute
	})).Return(nil)

	_, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_TanggalSelesai_AutoDeactivates() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.TanggalMulai = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianJabatanRequest{TanggalSelesai: dateOnly(2024, 12, 31)}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return !m.IsAktif && m.TanggalSelesai != nil
	})).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
	s.Require().NotNil(result.TanggalSelesai)
	s.Equal("2024-12-31", fmtDate(result.TanggalSelesai))
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_TanggalSelesai_ExplicitIsAktifWins() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.TanggalMulai = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianJabatanRequest{
		TanggalSelesai: dateOnly(2024, 12, 31),
		IsAktif:        ptrTo(true),
	}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.MatchedBy(func(m *models.KepegawaianJabatan) bool {
		return m.IsAktif
	})).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.True(result.IsAktif)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_TanggalSelesaiBeforeMulai() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.TanggalMulai = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianJabatanRequest{TanggalSelesai: dateOnly(2024, 1, 1)}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "UpdateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_SpecializationJobTitleMismatch() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.JobTitleID = 1
	req := &dto.UpdateKepegawaianJabatanRequest{SpecializationID: ptrTo(int64(5))}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationByID", int64(5)).Return(&models.Specialization{
		ID: 5, JobTitleID: ptrTo(int64(2)),
	}, nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
	s.repo.AssertNotCalled(s.T(), "UpdateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_ChangeJobTitle_ValidatesExistingSpecialization() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.JobTitleID = 1
	existing.SpecializationID = ptrTo(int64(5)) // milik job title 1
	req := &dto.UpdateKepegawaianJabatanRequest{JobTitleID: ptrTo(int64(2))}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationByID", int64(5)).Return(&models.Specialization{
		ID: 5, JobTitleID: ptrTo(int64(1)),
	}, nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusUnprocessableEntity)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_PrimaryConflict() {
	actor := superadminActor()
	existing := newJabatan(1) // non-primer, aktif
	req := &dto.UpdateKepegawaianJabatanRequest{IsPrimary: ptrTo(true)}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsPrimaryAktifByPegawai", existing.PegawaiID, int64(1)).Return(true, nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.assertAppErrCode(err, http.StatusConflict)
	s.repo.AssertNotCalled(s.T(), "UpdateJabatan", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_Primary_ExcludesSelf() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.IsPrimary = true // sudah primer; edit field lain tidak boleh dianggap konflik dengan dirinya
	req := &dto.UpdateKepegawaianJabatanRequest{NomorSK: ptrTo("SK/2")}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsPrimaryAktifByPegawai", existing.PegawaiID, int64(1)).Return(false, nil)
	s.repo.On("UpdateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_ClosePrimary_SkipsPrimaryCheck() {
	actor := superadminActor()
	existing := newJabatan(1)
	existing.IsPrimary = true
	existing.TanggalMulai = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &dto.UpdateKepegawaianJabatanRequest{TanggalSelesai: dateOnly(2024, 12, 31)} // mutasi: tutup jabatan primer

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(nil)

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.False(result.IsAktif)
	s.repo.AssertNotCalled(s.T(), "ExistsPrimaryAktifByPegawai", mock.Anything, mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJabatan_RepoError() {
	actor := superadminActor()
	existing := newJabatan(1)
	req := &dto.UpdateKepegawaianJabatanRequest{}

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJabatan", mock.AnythingOfType("*models.KepegawaianJabatan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateJabatan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJabatan_Superadmin_Success() {
	actor := superadminActor()
	existing := newJabatan(1)

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJabatan", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJabatan(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJabatan_WithPermission_Success() {
	actor := regularActor()
	existing := newJabatan(1)

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyDelete).Return(true, nil)
	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJabatan", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJabatan(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJabatan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJabatan(context.Background(), 1, actor)

	s.assertAppErrCode(err, http.StatusForbidden)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJabatan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJabatanByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJabatan(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJabatan_RepoError() {
	actor := superadminActor()
	existing := newJabatan(1)

	s.repo.On("GetJabatanByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJabatan", int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteJabatan(context.Background(), 1, actor)

	s.Error(err)
}
