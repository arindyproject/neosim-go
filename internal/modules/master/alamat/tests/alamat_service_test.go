package tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"
	"neosim_go/internal/modules/master/alamat/services"
	"neosim_go/internal/modules/master/alamat/tests/factories"
	"neosim_go/internal/modules/master/alamat/tests/mocks"

	alamatContracts "neosim_go/internal/modules/master/alamat/contracts"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  MasterAlamat Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/master/alamat")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/master/alamat")
	}

	os.Exit(code)
}

type MasterAlamatServiceTestSuite struct {
	suite.Suite
	repo     *mocks.MasterAlamatRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	svc      alamatContracts.Service
}

func (s *MasterAlamatServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MasterAlamatRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	// ─── PENYESUAIAN CACHE ───────────────────────────────────────────────────
	// Kita passing `nil` untuk redis.Client dan `false` untuk cacheEnabled.
	// Tujuannya agar unit test tidak bergantung pada koneksi Redis yang sebenarnya
	// dan fokus menguji logika bisnis, validasi, serta permission seperti semula.
	s.svc = services.NewMasterAlamatService(s.repo, s.rbacRepo, s.authRepo, nil)
}

func TestMasterAlamatService(t *testing.T) {
	suite.Run(t, new(MasterAlamatServiceTestSuite))
}

func superadminActor() alamatContracts.AuthContext {
	return alamatContracts.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() alamatContracts.AuthContext {
	return alamatContracts.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *MasterAlamatServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}

// =====================================================================
// NEGARA
// =====================================================================

func (s *MasterAlamatServiceTestSuite) Test_CreateNegara_Superadmin_Success() {
	req := &dto.CreateNegaraRequest{Code: "ID", Name: "Indonesia"}
	actor := superadminActor()

	s.repo.On("CreateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(nil)

	result, err := s.svc.CreateNegara(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_CreateNegara_WithPermission_Success() {
	req := &dto.CreateNegaraRequest{Code: "ID", Name: "Indonesia"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("CreateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(nil)

	result, err := s.svc.CreateNegara(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateNegara_WithManagePermission_Success() {
	req := &dto.CreateNegaraRequest{Code: "ID", Name: "Indonesia"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("CreateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(nil)

	result, err := s.svc.CreateNegara(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateNegara_Forbidden() {
	req := &dto.CreateNegaraRequest{Code: "ID", Name: "Indonesia"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateNegara(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateNegara_RepoError() {
	req := &dto.CreateNegaraRequest{Code: "ID", Name: "Indonesia"}
	actor := superadminActor()

	s.repo.On("CreateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateNegara(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDNegara_Success() {
	item := factories.NewNegaraFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDNegara", int64(1)).Return(item, nil)

	result, err := s.svc.GetByIDNegara(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDNegara_NotFound() {
	s.repo.On("GetByIDNegara", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDNegara(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDNegara_RepoError() {
	s.repo.On("GetByIDNegara", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDNegara(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_ListNegara_Success() {
	filter := &dto.FilterNegaraRequest{}
	items := []models.MasterAlamatNegara{
		*factories.NewNegaraFactory().Make(),
		*factories.NewNegaraFactory().Make(),
	}

	s.repo.On("ListNegara", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListNegara(1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterAlamatServiceTestSuite) Test_ListNegara_NotFound() {
	filter := &dto.FilterNegaraRequest{}

	s.repo.On("ListNegara", 1, 10, filter).Return([]models.MasterAlamatNegara{}, int64(0), nil)

	result, total, err := s.svc.ListNegara(1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_ListNegara_DefaultPagination() {
	filter := &dto.FilterNegaraRequest{}
	items := []models.MasterAlamatNegara{*factories.NewNegaraFactory().Make()}

	s.repo.On("ListNegara", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListNegara(0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListNegara", 1, 10, filter)
}

func (s *MasterAlamatServiceTestSuite) Test_ListNegara_PageSizeCapped() {
	filter := &dto.FilterNegaraRequest{}
	items := []models.MasterAlamatNegara{*factories.NewNegaraFactory().Make()}

	s.repo.On("ListNegara", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListNegara(1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListNegara", 1, 10, filter)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateNegara_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewNegaraFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateNegaraRequest{Name: &newName}

	s.repo.On("GetByIDNegara", int64(1)).Return(existing, nil)
	s.repo.On("UpdateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(nil)

	result, err := s.svc.UpdateNegara(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateNegara_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateNegaraRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateNegara(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateNegara_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateNegaraRequest{}

	s.repo.On("GetByIDNegara", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateNegara(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateNegara_PartialFields() {
	actor := superadminActor()
	existing := factories.NewNegaraFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newDesc := "New description"

	req := &dto.UpdateNegaraRequest{Description: &newDesc}

	s.repo.On("GetByIDNegara", int64(1)).Return(existing, nil)
	s.repo.On("UpdateNegara", mock.MatchedBy(func(m *models.MasterAlamatNegara) bool {
		return m.Name == originalName && *m.Description == newDesc
	})).Return(nil)

	result, err := s.svc.UpdateNegara(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newDesc, *result.Description)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateNegara_RepoError() {
	actor := superadminActor()
	existing := factories.NewNegaraFactory().Make()
	existing.ID = 1
	req := &dto.UpdateNegaraRequest{}

	s.repo.On("GetByIDNegara", int64(1)).Return(existing, nil)
	s.repo.On("UpdateNegara", mock.AnythingOfType("*models.MasterAlamatNegara")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateNegara(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteNegara_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewNegaraFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDNegara", int64(1)).Return(existing, nil)
	s.repo.On("DeleteNegara", int64(1)).Return(nil)

	err := s.svc.DeleteNegara(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteNegara_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteNegara(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteNegara_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDNegara", int64(999)).Return(nil, nil)

	err := s.svc.DeleteNegara(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteNegara_RepoError() {
	actor := superadminActor()
	existing := factories.NewNegaraFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDNegara", int64(1)).Return(existing, nil)
	s.repo.On("DeleteNegara", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteNegara(1, actor)

	s.Error(err)
}

// =====================================================================
// PROVINSI
// =====================================================================

func (s *MasterAlamatServiceTestSuite) Test_CreateProvinsi_Superadmin_Success() {
	req := &dto.CreateProvinsiRequest{NegaraID: 1, Code: "35", Name: "Jawa Timur"}
	actor := superadminActor()
	negara := factories.NewNegaraFactory().Make()
	negara.ID = 1

	s.repo.On("GetByIDNegara", int64(1)).Return(negara, nil)
	s.repo.On("CreateProvinsi", mock.AnythingOfType("*models.MasterAlamatProvinsi")).Return(nil)

	result, err := s.svc.CreateProvinsi(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_CreateProvinsi_WithPermission_Success() {
	req := &dto.CreateProvinsiRequest{NegaraID: 1, Code: "35", Name: "Jawa Timur"}
	actor := regularActor()
	negara := factories.NewNegaraFactory().Make()
	negara.ID = 1

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByIDNegara", int64(1)).Return(negara, nil)
	s.repo.On("CreateProvinsi", mock.AnythingOfType("*models.MasterAlamatProvinsi")).Return(nil)

	result, err := s.svc.CreateProvinsi(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateProvinsi_Forbidden() {
	req := &dto.CreateProvinsiRequest{NegaraID: 1, Code: "35", Name: "Jawa Timur"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateProvinsi(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateProvinsi_NegaraNotFound() {
	req := &dto.CreateProvinsiRequest{NegaraID: 999, Code: "35", Name: "Jawa Timur"}
	actor := superadminActor()

	s.repo.On("GetByIDNegara", int64(999)).Return(nil, nil)

	result, err := s.svc.CreateProvinsi(req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Negara tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_CreateProvinsi_RepoError() {
	req := &dto.CreateProvinsiRequest{NegaraID: 1, Code: "35", Name: "Jawa Timur"}
	actor := superadminActor()
	negara := factories.NewNegaraFactory().Make()
	negara.ID = 1

	s.repo.On("GetByIDNegara", int64(1)).Return(negara, nil)
	s.repo.On("CreateProvinsi", mock.AnythingOfType("*models.MasterAlamatProvinsi")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateProvinsi(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDProvinsi_Success() {
	item := factories.NewProvinsiFactory().Make()
	item.ID = 1
	item.Negara = models.MasterAlamatNegara{Name: "Indonesia"}

	s.repo.On("GetByIDProvinsi", int64(1)).Return(item, nil)
	s.repo.On("CountKotaByProvinsiID", int64(1)).Return(int64(5), nil)
	s.repo.On("CountKecamatanByProvinsiID", int64(1)).Return(int64(20), nil)
	s.repo.On("CountDesaByProvinsiID", int64(1)).Return(int64(100), nil)

	result, err := s.svc.GetByIDProvinsi(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal("Indonesia", result.NegaraName)
	s.Equal(int64(5), result.TotalKota)
	s.Equal(int64(20), result.TotalKecamatan)
	s.Equal(int64(100), result.TotalDesa)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDProvinsi_NotFound() {
	s.repo.On("GetByIDProvinsi", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDProvinsi(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDProvinsi_CountError() {
	item := factories.NewProvinsiFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDProvinsi", int64(1)).Return(item, nil)
	s.repo.On("CountKotaByProvinsiID", int64(1)).Return(int64(0), fmt.Errorf("db error"))

	result, err := s.svc.GetByIDProvinsi(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_ListProvinsi_Success() {
	filter := &dto.FilterProvinsiRequest{}
	items := []models.MasterAlamatProvinsi{
		*factories.NewProvinsiFactory().Make(),
		*factories.NewProvinsiFactory().Make(),
	}

	s.repo.On("ListProvinsi", 1, 10, (*int64)(nil), filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListProvinsi(1, 10, nil, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterAlamatServiceTestSuite) Test_ListProvinsi_WithNegaraIDFilter() {
	var negaraID int64 = 1
	filter := &dto.FilterProvinsiRequest{}
	items := []models.MasterAlamatProvinsi{*factories.NewProvinsiFactory().Make()}

	s.repo.On("ListProvinsi", 1, 10, &negaraID, filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListProvinsi(1, 10, &negaraID, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterAlamatServiceTestSuite) Test_ListProvinsi_NotFound() {
	filter := &dto.FilterProvinsiRequest{}

	s.repo.On("ListProvinsi", 1, 10, (*int64)(nil), filter).Return([]models.MasterAlamatProvinsi{}, int64(0), nil)

	result, total, err := s.svc.ListProvinsi(1, 10, nil, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateProvinsi_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewProvinsiFactory().Make()
	existing.ID = 1
	newName := "Updated Provinsi"
	req := &dto.UpdateProvinsiRequest{Name: &newName}

	s.repo.On("GetByIDProvinsi", int64(1)).Return(existing, nil)
	s.repo.On("UpdateProvinsi", mock.AnythingOfType("*models.MasterAlamatProvinsi")).Return(nil)

	result, err := s.svc.UpdateProvinsi(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateProvinsi_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateProvinsiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateProvinsi(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateProvinsi_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateProvinsiRequest{}

	s.repo.On("GetByIDProvinsi", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateProvinsi(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateProvinsi_RepoError() {
	actor := superadminActor()
	existing := factories.NewProvinsiFactory().Make()
	existing.ID = 1
	req := &dto.UpdateProvinsiRequest{}

	s.repo.On("GetByIDProvinsi", int64(1)).Return(existing, nil)
	s.repo.On("UpdateProvinsi", mock.AnythingOfType("*models.MasterAlamatProvinsi")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateProvinsi(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteProvinsi_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewProvinsiFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDProvinsi", int64(1)).Return(existing, nil)
	s.repo.On("DeleteProvinsi", int64(1)).Return(nil)

	err := s.svc.DeleteProvinsi(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteProvinsi_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteProvinsi(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteProvinsi_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDProvinsi", int64(999)).Return(nil, nil)

	err := s.svc.DeleteProvinsi(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteProvinsi_RepoError() {
	actor := superadminActor()
	existing := factories.NewProvinsiFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDProvinsi", int64(1)).Return(existing, nil)
	s.repo.On("DeleteProvinsi", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteProvinsi(1, actor)

	s.Error(err)
}

// =====================================================================
// KOTA / KABUPATEN
// =====================================================================

func (s *MasterAlamatServiceTestSuite) Test_CreateKotaKabupaten_Superadmin_Success() {
	req := &dto.CreateKotaKabupatenRequest{ProvinsiID: 1, Code: "35.21", Name: "Surabaya"}
	actor := superadminActor()
	provinsi := factories.NewProvinsiFactory().Make()
	provinsi.ID = 1

	s.repo.On("GetByIDProvinsi", int64(1)).Return(provinsi, nil)
	s.repo.On("CreateKotaKabupaten", mock.AnythingOfType("*models.MasterAlamatKotaKabupaten")).Return(nil)

	result, err := s.svc.CreateKotaKabupaten(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKotaKabupaten_Forbidden() {
	req := &dto.CreateKotaKabupatenRequest{ProvinsiID: 1, Code: "35.21", Name: "Surabaya"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKotaKabupaten(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKotaKabupaten_ProvinsiNotFound() {
	req := &dto.CreateKotaKabupatenRequest{ProvinsiID: 999, Code: "35.21", Name: "Surabaya"}
	actor := superadminActor()

	s.repo.On("GetByIDProvinsi", int64(999)).Return(nil, nil)

	result, err := s.svc.CreateKotaKabupaten(req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Provinsi tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKotaKabupaten_RepoError() {
	req := &dto.CreateKotaKabupatenRequest{ProvinsiID: 1, Code: "35.21", Name: "Surabaya"}
	actor := superadminActor()
	provinsi := factories.NewProvinsiFactory().Make()
	provinsi.ID = 1

	s.repo.On("GetByIDProvinsi", int64(1)).Return(provinsi, nil)
	s.repo.On("CreateKotaKabupaten", mock.AnythingOfType("*models.MasterAlamatKotaKabupaten")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKotaKabupaten(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKotaKabupaten_Success() {
	item := factories.NewKotaKabupatenFactory().Make()
	item.ID = 1
	item.Provinsi = models.MasterAlamatProvinsi{
		NegaraID: 1,
		Name:     "Jawa Timur",
		Negara:   models.MasterAlamatNegara{Name: "Indonesia"},
	}

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(item, nil)
	s.repo.On("CountKecamatanByKotaID", int64(1)).Return(int64(10), nil)
	s.repo.On("CountDesaByKotaID", int64(1)).Return(int64(50), nil)

	result, err := s.svc.GetByIDKotaKabupaten(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Jawa Timur", result.ProvinsiName)
	s.Equal("Indonesia", result.NegaraName)
	s.Equal(int64(10), result.TotalKecamatan)
	s.Equal(int64(50), result.TotalDesa)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKotaKabupaten_NotFound() {
	s.repo.On("GetByIDKotaKabupaten", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDKotaKabupaten(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKotaKabupaten_CountError() {
	item := factories.NewKotaKabupatenFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(item, nil)
	s.repo.On("CountKecamatanByKotaID", int64(1)).Return(int64(0), fmt.Errorf("db error"))

	result, err := s.svc.GetByIDKotaKabupaten(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKotaKabupaten_Success() {
	filter := &dto.FilterKotaKabupatenRequest{}
	items := []models.MasterAlamatKotaKabupaten{
		*factories.NewKotaKabupatenFactory().Make(),
	}

	s.repo.On("ListKotaKabupaten", 1, 10, (*int64)(nil), filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKotaKabupaten(1, 10, nil, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKotaKabupaten_NotFound() {
	filter := &dto.FilterKotaKabupatenRequest{}

	s.repo.On("ListKotaKabupaten", 1, 10, (*int64)(nil), filter).Return([]models.MasterAlamatKotaKabupaten{}, int64(0), nil)

	result, total, err := s.svc.ListKotaKabupaten(1, 10, nil, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKotaKabupaten_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKotaKabupatenFactory().Make()
	existing.ID = 1
	newName := "Updated Kota"
	req := &dto.UpdateKotaKabupatenRequest{Name: &newName}

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKotaKabupaten", mock.AnythingOfType("*models.MasterAlamatKotaKabupaten")).Return(nil)

	result, err := s.svc.UpdateKotaKabupaten(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKotaKabupaten_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKotaKabupatenRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKotaKabupaten(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKotaKabupaten_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKotaKabupatenRequest{}

	s.repo.On("GetByIDKotaKabupaten", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKotaKabupaten(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKotaKabupaten_RepoError() {
	actor := superadminActor()
	existing := factories.NewKotaKabupatenFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKotaKabupatenRequest{}

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKotaKabupaten", mock.AnythingOfType("*models.MasterAlamatKotaKabupaten")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKotaKabupaten(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKotaKabupaten_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKotaKabupatenFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKotaKabupaten", int64(1)).Return(nil)

	err := s.svc.DeleteKotaKabupaten(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKotaKabupaten_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKotaKabupaten(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKotaKabupaten_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDKotaKabupaten", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKotaKabupaten(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKotaKabupaten_RepoError() {
	actor := superadminActor()
	existing := factories.NewKotaKabupatenFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKotaKabupaten", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKotaKabupaten(1, actor)

	s.Error(err)
}

// =====================================================================
// KECAMATAN
// =====================================================================

func (s *MasterAlamatServiceTestSuite) Test_CreateKecamatan_Superadmin_Success() {
	req := &dto.CreateKecamatanRequest{KotaKabupatenID: 1, Code: "35.21.01", Name: "Gubeng"}
	actor := superadminActor()
	kota := factories.NewKotaKabupatenFactory().Make()
	kota.ID = 1

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(kota, nil)
	s.repo.On("CreateKecamatan", mock.AnythingOfType("*models.MasterAlamatKecamatan")).Return(nil)

	result, err := s.svc.CreateKecamatan(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKecamatan_Forbidden() {
	req := &dto.CreateKecamatanRequest{KotaKabupatenID: 1, Code: "35.21.01", Name: "Gubeng"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKecamatan(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKecamatan_KotaNotFound() {
	req := &dto.CreateKecamatanRequest{KotaKabupatenID: 999, Code: "35.21.01", Name: "Gubeng"}
	actor := superadminActor()

	s.repo.On("GetByIDKotaKabupaten", int64(999)).Return(nil, nil)

	result, err := s.svc.CreateKecamatan(req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Kota/Kabupaten tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKecamatan_RepoError() {
	req := &dto.CreateKecamatanRequest{KotaKabupatenID: 1, Code: "35.21.01", Name: "Gubeng"}
	actor := superadminActor()
	kota := factories.NewKotaKabupatenFactory().Make()
	kota.ID = 1

	s.repo.On("GetByIDKotaKabupaten", int64(1)).Return(kota, nil)
	s.repo.On("CreateKecamatan", mock.AnythingOfType("*models.MasterAlamatKecamatan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKecamatan(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKecamatan_Success() {
	item := factories.NewKecamatanFactory().Make()
	item.ID = 1
	item.KotaKabupaten = models.MasterAlamatKotaKabupaten{
		ProvinsiID: 1,
		Name:       "Surabaya",
		Provinsi: models.MasterAlamatProvinsi{
			NegaraID: 1,
			Name:     "Jawa Timur",
			Negara:   models.MasterAlamatNegara{Name: "Indonesia"},
		},
	}

	s.repo.On("GetByIDKecamatan", int64(1)).Return(item, nil)
	s.repo.On("CountDesaByKecamatanID", int64(1)).Return(int64(15), nil)

	result, err := s.svc.GetByIDKecamatan(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Surabaya", result.KotaKabupatenName)
	s.Equal("Jawa Timur", result.ProvinsiName)
	s.Equal("Indonesia", result.NegaraName)
	s.Equal(int64(15), result.TotalDesa)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKecamatan_NotFound() {
	s.repo.On("GetByIDKecamatan", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDKecamatan(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKecamatan_CountError() {
	item := factories.NewKecamatanFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDKecamatan", int64(1)).Return(item, nil)
	s.repo.On("CountDesaByKecamatanID", int64(1)).Return(int64(0), fmt.Errorf("db error"))

	result, err := s.svc.GetByIDKecamatan(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKecamatan_Success() {
	filter := &dto.FilterKecamatanRequest{}
	items := []models.MasterAlamatKecamatan{
		*factories.NewKecamatanFactory().Make(),
	}

	s.repo.On("ListKecamatan", 1, 10, (*int64)(nil), filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKecamatan(1, 10, nil, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKecamatan_NotFound() {
	filter := &dto.FilterKecamatanRequest{}

	s.repo.On("ListKecamatan", 1, 10, (*int64)(nil), filter).Return([]models.MasterAlamatKecamatan{}, int64(0), nil)

	result, total, err := s.svc.ListKecamatan(1, 10, nil, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKecamatan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKecamatanFactory().Make()
	existing.ID = 1
	newName := "Updated Kecamatan"
	req := &dto.UpdateKecamatanRequest{Name: &newName}

	s.repo.On("GetByIDKecamatan", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKecamatan", mock.AnythingOfType("*models.MasterAlamatKecamatan")).Return(nil)

	result, err := s.svc.UpdateKecamatan(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKecamatan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKecamatanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKecamatan(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKecamatan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKecamatanRequest{}

	s.repo.On("GetByIDKecamatan", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKecamatan(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKecamatan_RepoError() {
	actor := superadminActor()
	existing := factories.NewKecamatanFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKecamatanRequest{}

	s.repo.On("GetByIDKecamatan", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKecamatan", mock.AnythingOfType("*models.MasterAlamatKecamatan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKecamatan(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKecamatan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKecamatanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKecamatan", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKecamatan", int64(1)).Return(nil)

	err := s.svc.DeleteKecamatan(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKecamatan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKecamatan(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKecamatan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDKecamatan", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKecamatan(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKecamatan_RepoError() {
	actor := superadminActor()
	existing := factories.NewKecamatanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKecamatan", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKecamatan", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKecamatan(1, actor)

	s.Error(err)
}

// =====================================================================
// KELURAHAN / DESA
// =====================================================================

func (s *MasterAlamatServiceTestSuite) Test_CreateKelurahanDesa_Superadmin_Success() {
	req := &dto.CreateKelurahanDesaRequest{KecamatanID: 1, Code: "35.21.01.2001", Name: "Airlangga"}
	actor := superadminActor()
	kecamatan := factories.NewKecamatanFactory().Make()
	kecamatan.ID = 1

	s.repo.On("GetByIDKecamatan", int64(1)).Return(kecamatan, nil)
	s.repo.On("CreateKelurahanDesa", mock.AnythingOfType("*models.MasterAlamatKelurahanDesa")).Return(nil)

	result, err := s.svc.CreateKelurahanDesa(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKelurahanDesa_Forbidden() {
	req := &dto.CreateKelurahanDesaRequest{KecamatanID: 1, Code: "35.21.01.2001", Name: "Airlangga"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateKelurahanDesa(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKelurahanDesa_KecamatanNotFound() {
	req := &dto.CreateKelurahanDesaRequest{KecamatanID: 999, Code: "35.21.01.2001", Name: "Airlangga"}
	actor := superadminActor()

	s.repo.On("GetByIDKecamatan", int64(999)).Return(nil, nil)

	result, err := s.svc.CreateKelurahanDesa(req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "Kecamatan tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_CreateKelurahanDesa_RepoError() {
	req := &dto.CreateKelurahanDesaRequest{KecamatanID: 1, Code: "35.21.01.2001", Name: "Airlangga"}
	actor := superadminActor()
	kecamatan := factories.NewKecamatanFactory().Make()
	kecamatan.ID = 1

	s.repo.On("GetByIDKecamatan", int64(1)).Return(kecamatan, nil)
	s.repo.On("CreateKelurahanDesa", mock.AnythingOfType("*models.MasterAlamatKelurahanDesa")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateKelurahanDesa(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKelurahanDesa_Success() {
	item := factories.NewKelurahanDesaFactory().Make()
	item.ID = 1
	item.Kecamatan = models.MasterAlamatKecamatan{
		KotaKabupatenID: 1,
		Name:            "Gubeng",
		KotaKabupaten: models.MasterAlamatKotaKabupaten{
			ProvinsiID: 1,
			Name:       "Surabaya",
			Provinsi: models.MasterAlamatProvinsi{
				NegaraID: 1,
				Name:     "Jawa Timur",
				Negara:   models.MasterAlamatNegara{Name: "Indonesia"},
			},
		},
	}

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(item, nil)

	result, err := s.svc.GetByIDKelurahanDesa(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Gubeng", result.KecamatanName)
	s.Equal("Surabaya", result.KotaKabupatenName)
	s.Equal("Jawa Timur", result.ProvinsiName)
	s.Equal("Indonesia", result.NegaraName)
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKelurahanDesa_NotFound() {
	s.repo.On("GetByIDKelurahanDesa", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDKelurahanDesa(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_GetByIDKelurahanDesa_RepoError() {
	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDKelurahanDesa(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKelurahanDesa_Success() {
	filter := &dto.FilterKelurahanDesaRequest{}
	items := []models.MasterAlamatKelurahanDesa{
		*factories.NewKelurahanDesaFactory().Make(),
	}

	s.repo.On("ListKelurahanDesa", 1, 10, (*int64)(nil), filter).Return(items, int64(1), nil)

	result, total, err := s.svc.ListKelurahanDesa(1, 10, nil, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(result, 1)
}

func (s *MasterAlamatServiceTestSuite) Test_ListKelurahanDesa_NotFound() {
	filter := &dto.FilterKelurahanDesaRequest{}

	s.repo.On("ListKelurahanDesa", 1, 10, (*int64)(nil), filter).Return([]models.MasterAlamatKelurahanDesa{}, int64(0), nil)

	result, total, err := s.svc.ListKelurahanDesa(1, 10, nil, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKelurahanDesa_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKelurahanDesaFactory().Make()
	existing.ID = 1
	newName := "Updated Desa"
	req := &dto.UpdateKelurahanDesaRequest{Name: &newName}

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKelurahanDesa", mock.AnythingOfType("*models.MasterAlamatKelurahanDesa")).Return(nil)

	result, err := s.svc.UpdateKelurahanDesa(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKelurahanDesa_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateKelurahanDesaRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateKelurahanDesa(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKelurahanDesa_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateKelurahanDesaRequest{}

	s.repo.On("GetByIDKelurahanDesa", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateKelurahanDesa(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKelurahanDesa_PartialFields() {
	actor := superadminActor()
	existing := factories.NewKelurahanDesaFactory().Make()
	existing.ID = 1
	originalName := existing.Name
	newPostal := "60286"

	req := &dto.UpdateKelurahanDesaRequest{PostalCode: &newPostal}

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKelurahanDesa", mock.MatchedBy(func(m *models.MasterAlamatKelurahanDesa) bool {
		return m.Name == originalName && *m.PostalCode == newPostal
	})).Return(nil)

	result, err := s.svc.UpdateKelurahanDesa(1, req, actor)

	s.NoError(err)
	s.Equal(originalName, result.Name)
	s.Equal(newPostal, *result.PostalCode)
}

func (s *MasterAlamatServiceTestSuite) Test_UpdateKelurahanDesa_RepoError() {
	actor := superadminActor()
	existing := factories.NewKelurahanDesaFactory().Make()
	existing.ID = 1
	req := &dto.UpdateKelurahanDesaRequest{}

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(existing, nil)
	s.repo.On("UpdateKelurahanDesa", mock.AnythingOfType("*models.MasterAlamatKelurahanDesa")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateKelurahanDesa(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKelurahanDesa_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewKelurahanDesaFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKelurahanDesa", int64(1)).Return(nil)

	err := s.svc.DeleteKelurahanDesa(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKelurahanDesa_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteKelurahanDesa(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKelurahanDesa_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDKelurahanDesa", int64(999)).Return(nil, nil)

	err := s.svc.DeleteKelurahanDesa(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterAlamatServiceTestSuite) Test_DeleteKelurahanDesa_RepoError() {
	actor := superadminActor()
	existing := factories.NewKelurahanDesaFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDKelurahanDesa", int64(1)).Return(existing, nil)
	s.repo.On("DeleteKelurahanDesa", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteKelurahanDesa(1, actor)

	s.Error(err)
}
