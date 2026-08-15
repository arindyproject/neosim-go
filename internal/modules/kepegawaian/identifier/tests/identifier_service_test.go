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

type KepegawaianIdentifierServiceTestSuite struct {
	suite.Suite
	repo     *mocks.KepegawaianIdentifierRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	userRepo *mocks.UserRepositoryMock
	svc      identifierContracts.Service
	cfg      *config.Config
	ctx      context.Context
}

func (s *KepegawaianIdentifierServiceTestSuite) SetupTest() {
	s.repo = new(mocks.KepegawaianIdentifierRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.userRepo = new(mocks.UserRepositoryMock)
	s.cfg = &config.Config{DefaultPageSize: 10, DefaultPageSizeMax: 100}
	s.svc = services.NewKepegawaianIdentifierService(s.repo, s.rbacRepo, s.authRepo, s.userRepo, s.cfg)
	s.ctx = context.Background()

	// Stub default agar buildCreator/buildAuditMaps tidak panic saat memanggil userRepo.
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

func (s *KepegawaianIdentifierServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

// sampleTipe mengembalikan master tipe dummy untuk stub GetTipeByID.
func sampleTipe(id int64, code string) *models.Tipe {
	return &models.Tipe{ID: id, Code: code}
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_Superadmin_Success() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		TipeID:    1,
		Nilai:     "3201010101010001",
		IsPrimary: false,
		IsAktif:   true,
	}

	s.repo.On("GetTipeByID", int64(1)).Return(sampleTipe(1, "NIK"), nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), req.Nilai, int64(0)).Return(false, nil)
	s.repo.On("CreateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Nilai, result.Nilai)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_WithPermission_Success() {
	actor := regularActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		TipeID:    1,
		Nilai:     "3201010101010001",
		IsAktif:   true,
	}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetTipeByID", int64(1)).Return(sampleTipe(1, "NIK"), nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), req.Nilai, int64(0)).Return(false, nil)
	s.repo.On("CreateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_Forbidden() {
	actor := regularActor()
	req := &dto.CreateKepegawaianIdentifierRequest{PegawaiID: 10, TipeID: 1, Nilai: "123"}
	s.mockNoPermissions()

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_TipeTidakDitemukan() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{PegawaiID: 10, TipeID: 999, Nilai: "123"}

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_Duplicate() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{PegawaiID: 10, TipeID: 1, Nilai: "3201010101010001"}

	s.repo.On("GetTipeByID", int64(1)).Return(sampleTipe(1, "NIK"), nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), req.Nilai, int64(0)).Return(true, nil)

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_SetPrimary_UnsetLama() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{
		PegawaiID: 10,
		TipeID:    1,
		Nilai:     "3201010101010001",
		IsPrimary: true,
	}

	s.repo.On("GetTipeByID", int64(1)).Return(sampleTipe(1, "NIK"), nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), req.Nilai, int64(0)).Return(false, nil)
	s.repo.On("UnsetPrimaryByPegawaiIDAndTipe", s.ctx, int64(10), int64(1), actor.UserID).Return(nil)
	s.repo.On("CreateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateIdentifier_RepoError() {
	actor := superadminActor()
	req := &dto.CreateKepegawaianIdentifierRequest{PegawaiID: 10, TipeID: 1, Nilai: "123"}

	s.repo.On("GetTipeByID", int64(1)).Return(sampleTipe(1, "NIK"), nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), req.Nilai, int64(0)).Return(false, nil)
	s.repo.On("CreateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateIdentifier(s.ctx, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_Superadmin_Success() {
	actor := superadminActor()
	item := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("Nilai", "3201010101010001").
		Make()

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(item, nil)

	result, err := s.svc.GetIdentifierByID(s.ctx, 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Nilai, result.Nilai)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_WithPermission_Success() {
	actor := regularActor()
	item := factories.NewKepegawaianIdentifierFactory().With("ID", int64(1)).Make()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyRead, mock.Anything).Return(true, nil)
	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(item, nil)

	result, err := s.svc.GetIdentifierByID(s.ctx, 1, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetIdentifierByID(s.ctx, 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", s.ctx, int64(999)).Return(nil, nil)

	result, err := s.svc.GetIdentifierByID(s.ctx, 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetIdentifierByID_RepoError() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetIdentifierByID(s.ctx, 1, actor)

	s.Nil(result)
	s.Error(err)
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_Superadmin_Success() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().Make(),
		*factories.NewKepegawaianIdentifierFactory().Make(),
	}

	s.repo.On("ListIdentifier", s.ctx, 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListIdentifier(s.ctx, 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListIdentifier(s.ctx, 1, 10, filter, actor)

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

	s.repo.On("ListIdentifier", s.ctx, 1, 10, filter).Return([]models.KepegawaianIdentifier{}, int64(0), nil)

	result, total, err := s.svc.ListIdentifier(s.ctx, 0, 0, filter, actor)

	s.NoError(err)
	s.Equal(int64(0), total)
	s.Empty(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_PageSizeCapped() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{}

	s.repo.On("ListIdentifier", s.ctx, 1, 10, filter).Return([]models.KepegawaianIdentifier{}, int64(0), nil)

	_, _, err := s.svc.ListIdentifier(s.ctx, 1, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListIdentifier", s.ctx, 1, 10, filter)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListIdentifier_WithNilaiFilter() {
	actor := superadminActor()
	filter := &dto.FilterKepegawaianIdentifierRequest{Nilai: "3201"}
	items := []models.KepegawaianIdentifier{*factories.NewKepegawaianIdentifierFactory().Make()}

	s.repo.On("ListIdentifier", s.ctx, 1, 10, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListIdentifier(s.ctx, 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

// ── ListByPegawai ─────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListByPegawai_Superadmin_Success() {
	actor := superadminActor()
	items := []models.KepegawaianIdentifier{
		*factories.NewKepegawaianIdentifierFactory().With("PegawaiID", int64(10)).Make(),
	}

	s.repo.On("FindByPegawaiID", s.ctx, int64(10)).Return(items, nil)

	result, err := s.svc.ListByPegawai(s.ctx, 10, actor)

	s.NoError(err)
	s.Len(result, 1)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListByPegawai_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.ListByPegawai(s.ctx, 10, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("Nilai", "3201010101010001").
		With("IsPrimary", false).
		Make()

	newNilai := "3201010101019999"
	req := &dto.UpdateKepegawaianIdentifierRequest{Nilai: &newNilai}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), newNilai, int64(1)).Return(false, nil)
	s.repo.On("UpdateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newNilai, result.Nilai)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_GantiTipe_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("Nilai", "3201010101010001").
		Make()

	newTipeID := int64(2)
	req := &dto.UpdateKepegawaianIdentifierRequest{TipeID: &newTipeID}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, newTipeID, existing.Nilai, int64(1)).Return(false, nil)
	s.repo.On("GetTipeByID", newTipeID).Return(sampleTipe(2, "STR"), nil)
	s.repo.On("UpdateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_SetPrimary_UnsetLama() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("IsPrimary", false).
		Make()

	isPrimary := true
	req := &dto.UpdateKepegawaianIdentifierRequest{IsPrimary: &isPrimary}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, existing.TipeID, existing.Nilai, int64(1)).Return(false, nil)
	s.repo.On("UnsetPrimaryByPegawaiIDAndTipe", s.ctx, int64(10), int64(1), actor.UserID).Return(nil)
	s.repo.On("UpdateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(nil)

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKepegawaianIdentifierRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKepegawaianIdentifierRequest{}

	s.repo.On("GetIdentifierByID", s.ctx, int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateIdentifier(s.ctx, 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_Duplicate() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("TipeID", int64(1)).
		With("Nilai", "3201010101010001").
		Make()

	newNilai := "3201010101019999"
	req := &dto.UpdateKepegawaianIdentifierRequest{Nilai: &newNilai}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, int64(1), newNilai, int64(1)).Return(true, nil)

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateIdentifier_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().With("ID", int64(1)).Make()
	req := &dto.UpdateKepegawaianIdentifierRequest{}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("ExistsByNilaiAndTipe", s.ctx, existing.TipeID, existing.Nilai, int64(1)).Return(false, nil)
	s.repo.On("UpdateIdentifier", s.ctx, mock.AnythingOfType("*models.KepegawaianIdentifier")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateIdentifier(s.ctx, 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("IsPrimary", false).
		Make()

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("DeleteIdentifier", s.ctx, int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteIdentifier(s.ctx, 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Primary_TanpaPenggantiAktif_Success() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("IsPrimary", true).
		Make()
	existing.ID = 1 // set langsung, hindari kemungkinan bug reflection di With()

	others := []models.KepegawaianIdentifier{*existing}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("FindByPegawaiIDAndTipe", s.ctx, int64(10), int64(1)).Return(others, nil)
	s.repo.On("DeleteIdentifier", s.ctx, int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteIdentifier(s.ctx, 1, actor)

	s.NoError(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Primary_AdaPenggantiAktif_Gagal() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("IsPrimary", true).
		Make()

	lainnya := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(2)).
		With("PegawaiID", int64(10)).
		With("TipeID", int64(1)).
		With("IsAktif", true).
		Make()

	others := []models.KepegawaianIdentifier{*existing, *lainnya}

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("FindByPegawaiIDAndTipe", s.ctx, int64(10), int64(1)).Return(others, nil)

	err := s.svc.DeleteIdentifier(s.ctx, 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "DeleteIdentifier", mock.Anything, mock.Anything, mock.Anything)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteIdentifier(s.ctx, 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_NotFound() {
	actor := superadminActor()

	s.repo.On("GetIdentifierByID", s.ctx, int64(999)).Return(nil, nil)

	err := s.svc.DeleteIdentifier(s.ctx, 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteIdentifier_RepoError() {
	actor := superadminActor()
	existing := factories.NewKepegawaianIdentifierFactory().
		With("ID", int64(1)).
		With("IsPrimary", false).
		Make()

	s.repo.On("GetIdentifierByID", s.ctx, int64(1)).Return(existing, nil)
	s.repo.On("DeleteIdentifier", s.ctx, int64(1), actor.UserID).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteIdentifier(s.ctx, 1, actor)

	s.Error(err)
}

// ── GetExpiringSoon / GetExpired ────────────────────────────────────────────

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetExpiringSoonIdentifier_Success() {
	actor := superadminActor()
	items := []models.KepegawaianIdentifier{*factories.NewKepegawaianIdentifierFactory().Make()}

	s.repo.On("FindExpiringSoonIdentifier", s.ctx, 30).Return(items, nil)

	result, err := s.svc.GetExpiringSoonIdentifier(s.ctx, 0, actor) // 0 → default 30 hari

	s.NoError(err)
	s.Len(result, 1)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetExpiredIdentifier_Success() {
	actor := superadminActor()
	items := []models.KepegawaianIdentifier{*factories.NewKepegawaianIdentifierFactory().Make()}

	s.repo.On("FindExpiredIdentifier", s.ctx).Return(items, nil)

	result, err := s.svc.GetExpiredIdentifier(s.ctx, actor)

	s.NoError(err)
	s.Len(result, 1)
}
