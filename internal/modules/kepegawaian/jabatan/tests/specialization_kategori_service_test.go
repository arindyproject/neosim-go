package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianJabatanServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, ptrTo, dst) sudah
// didefinisikan di jabatan_service_test.go. File ini HANYA menambah skenario
// test untuk SpecializationKategori, memakai s.svc / s.repo yang SAMA.

// mockSpecKategoriNoDuplicate men-stub kedua pengecekan duplikat (Code dan
// Label) agar dianggap belum dipakai. Dipanggil eksplisit per test, bukan
// dari SetupTest, supaya tidak bentrok dengan test yang menguji skenario 409.
func (s *KepegawaianJabatanServiceTestSuite) mockSpecKategoriNoDuplicate(code, label string) {
	s.repo.On("GetSpecializationKategoriByCode", code).Return(nil, nil).Once()
	s.repo.On("GetSpecializationKategoriByLabel", label).Return(nil, nil).Once()
}

func newCreateSpecKategoriReq() *dto.CreateSpecializationKategoriRequest {
	return &dto.CreateSpecializationKategoriRequest{
		Code:  "BEDAH",
		Label: "Bedah",
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_Superadmin_Success() {
	fhir := "surgical"
	req := newCreateSpecKategoriReq()
	req.FHIRCode = &fhir
	actor := superadminActor()

	s.mockSpecKategoriNoDuplicate(req.Code, req.Label)
	s.repo.On("CreateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
	s.Equal(req.FHIRCode, result.FHIRCode)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_Forbidden() {
	req := newCreateSpecKategoriReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "GetSpecializationKategoriByCode", mock.Anything)
	s.repo.AssertNotCalled(s.T(), "CreateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_DuplicateCode() {
	req := newCreateSpecKategoriReq()
	actor := superadminActor()

	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.Code = req.Code
	s.repo.On("GetSpecializationKategoriByCode", req.Code).Return(existing, nil)

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "GetSpecializationKategoriByLabel", mock.Anything)
	s.repo.AssertNotCalled(s.T(), "CreateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_DuplicateLabel() {
	req := newCreateSpecKategoriReq()
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByCode", req.Code).Return(nil, nil)
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.Label = req.Label
	s.repo.On("GetSpecializationKategoriByLabel", req.Label).Return(existing, nil)

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_CheckCodeRepoError() {
	req := newCreateSpecKategoriReq()
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByCode", req.Code).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	s.repo.AssertNotCalled(s.T(), "CreateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_CheckLabelRepoError() {
	req := newCreateSpecKategoriReq()
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByCode", req.Code).Return(nil, nil)
	s.repo.On("GetSpecializationKategoriByLabel", req.Label).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	s.repo.AssertNotCalled(s.T(), "CreateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_RepoError() {
	req := newCreateSpecKategoriReq()
	actor := superadminActor()

	s.mockSpecKategoriNoDuplicate(req.Code, req.Label)
	s.repo.On("CreateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateSpecializationKategori(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_Success() {
	actor := superadminActor()
	item := factories.NewSpecializationKategoriFactory().Make()
	item.ID = 1

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── List ─────────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecializationKategori_Success() {
	actor := superadminActor()
	filter := &dto.FilterSpecializationKategoriRequest{}
	items := []models.SpecializationKategori{
		*factories.NewSpecializationKategoriFactory().Make(),
		*factories.NewSpecializationKategoriFactory().Make(),
	}

	s.repo.On("ListSpecializationKategori", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListSpecializationKategori(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecializationKategori_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterSpecializationKategoriRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListSpecializationKategori(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── ListSelect ───────────────────────────────────────────────────────────────
// Catatan: ListSelectSpecializationKategori memakai cache (s.cache.Get/SetDefault).
// Karena SetupTest memakai cache.NewManager(nil, false, 0) (cache nonaktif),
// s.cache.Get selalu miss, jadi selalu jatuh ke s.repo.ListSelectSpecializationKategori.
// Kalau signature cache.Manager Anda berbeda dan tidak aman dipakai dengan nil
// client meski disabled, beri tahu saya supaya saya sesuaikan lagi.

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectSpecializationKategori_Success() {
	actor := superadminActor()
	items := []models.SpecializationKategori{
		*factories.NewSpecializationKategoriFactory().Make(),
	}

	s.repo.On("ListSelectSpecializationKategori", "").Return(items, nil)

	result, err := s.svc.ListSelectSpecializationKategori(context.Background(), "", actor)

	s.NoError(err)
	s.Len(result, 1)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectSpecializationKategori_Empty_NotFound() {
	actor := superadminActor()

	s.repo.On("ListSelectSpecializationKategori", "xxx").Return([]models.SpecializationKategori{}, nil)

	result, err := s.svc.ListSelectSpecializationKategori(context.Background(), "xxx", actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusNotFound, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectSpecializationKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.ListSelectSpecializationKategori(context.Background(), "", actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateSpecializationKategoriRequest{Label: &newLabel}

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationKategoriByLabel", newLabel).Return(nil, nil)
	s.repo.On("UpdateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_NoFieldChanged_SkipsDuplicateCheck() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	fhir := "baru"
	req := &dto.UpdateSpecializationKategoriRequest{FHIRCode: &fhir} // Code & Label tidak diubah

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(fhir, *result.FHIRCode)
	s.repo.AssertNotCalled(s.T(), "GetSpecializationKategoriByCode", mock.Anything)
	s.repo.AssertNotCalled(s.T(), "GetSpecializationKategoriByLabel", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_SameCode_SkipsDuplicateCheck() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	sameCode := existing.Code
	req := &dto.UpdateSpecializationKategoriRequest{Code: &sameCode} // dikirim tapi nilainya sama

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(sameCode, result.Code)
	s.repo.AssertNotCalled(s.T(), "GetSpecializationKategoriByCode", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_DuplicateCode() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	newCode := "SUDAH-DIPAKAI"
	req := &dto.UpdateSpecializationKategoriRequest{Code: &newCode}

	other := factories.NewSpecializationKategoriFactory().Make()
	other.ID = 2
	other.Code = newCode

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationKategoriByCode", newCode).Return(other, nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_DuplicateLabel() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	newLabel := "Sudah Dipakai"
	req := &dto.UpdateSpecializationKategoriRequest{Label: &newLabel}

	other := factories.NewSpecializationKategoriFactory().Make()
	other.ID = 2
	other.Label = newLabel

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationKategoriByLabel", newLabel).Return(other, nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateSpecializationKategori", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateSpecializationKategoriRequest{}

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateSpecializationKategoriRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateSpecializationKategori(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteSpecializationKategori", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteSpecializationKategori(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteSpecializationKategori(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteSpecializationKategori(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
