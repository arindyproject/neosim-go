package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianPegawaiServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di pegawai_service_test.go. File ini HANYA menambah skenario test untuk
// Status, memakai s.svc / s.repo yang SAMA.
//
// CATATAN PENTING: saya belum melihat isi status_service.go. Test di bawah
// HANYA menyesuaikan field Code/Label (dari Name lama), TANPA menambahkan
// skenario duplikat Code/Label karena belum terkonfirmasi apakah
// CreateStatus/UpdateStatus memvalidasi itu. Kalau ternyata service ini
// memanggil GetStatusByCode/GetStatusByLabel sebelum create/update (seperti
// SpecializationKategori), test create/update sukses di bawah akan panic
// "unexpected call" — kirim isi status_service.go kalau itu terjadi.

func newCreateStatusReq() *dto.CreateStatusRequest {
	return &dto.CreateStatusRequest{
		Code:  "TEST001",
		Label: "Test Status",
	}
}

func (s *KepegawaianPegawaiServiceTestSuite) mockStatusNoDuplicate(code, label string) {
	s.repo.On("GetStatusByCode", code).Return(nil, nil)
	s.repo.On("GetStatusByLabel", label).Return(nil, nil)
}

// ── Create ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateStatus_Superadmin_Success() {
	req := newCreateStatusReq()
	actor := superadminActor()

	s.mockStatusNoDuplicate(req.Code, req.Label)
	s.repo.On("CreateStatus", mock.AnythingOfType("*models.Status")).Return(nil)

	result, err := s.svc.CreateStatus(context.Background(), req, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateStatus_Forbidden() {
	req := newCreateStatusReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateStatus(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateStatus", mock.Anything)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateStatus_RepoError() {
	req := newCreateStatusReq()
	actor := superadminActor()

	s.mockStatusNoDuplicate(req.Code, req.Label)
	s.repo.On("CreateStatus", mock.AnythingOfType("*models.Status")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateStatus(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetStatusByID_Success() {
	actor := superadminActor()
	item := factories.NewStatusFactory().Make()
	item.ID = 1

	s.repo.On("GetStatusByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetStatusByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetStatusByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetStatusByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetStatusByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetStatusByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetStatusByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── List ─────────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListStatus_Success() {
	actor := superadminActor()
	filter := &dto.FilterStatusRequest{}
	items := []models.Status{
		*factories.NewStatusFactory().Make(),
		*factories.NewStatusFactory().Make(),
	}

	s.repo.On("ListStatus", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListStatus(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListStatus_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterStatusRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListStatus(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListStatus_RepoError() {
	actor := superadminActor()
	filter := &dto.FilterStatusRequest{}

	s.repo.On("ListStatus", 1, 10, filter).Return(nil, int64(0), fmt.Errorf("db error"))

	result, total, err := s.svc.ListStatus(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateStatus_Success() {
	actor := superadminActor()
	existing := factories.NewStatusFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateStatusRequest{Label: &newLabel}

	s.repo.On("GetStatusByID", int64(1)).Return(existing, nil)
	s.repo.On("GetStatusByLabel", newLabel).Return(nil, nil)
	s.repo.On("UpdateStatus", mock.AnythingOfType("*models.Status")).Return(nil)

	result, err := s.svc.UpdateStatus(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}
func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateStatus_PartialFields() {
	actor := superadminActor()
	existing := factories.NewStatusFactory().Make()
	existing.ID = 1
	originalCode := existing.Code
	newLabel := "Label Baru"
	req := &dto.UpdateStatusRequest{Label: &newLabel}

	s.repo.On("GetStatusByID", int64(1)).Return(existing, nil)
	s.repo.On("GetStatusByLabel", newLabel).Return(nil, nil)
	s.repo.On("UpdateStatus", mock.MatchedBy(func(m *models.Status) bool {
		return m.Code == originalCode && m.Label == newLabel
	})).Return(nil)

	result, err := s.svc.UpdateStatus(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalCode, result.Code)
	s.Equal(newLabel, result.Label)
}
func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateStatus_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateStatusRequest{}

	s.repo.On("GetStatusByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateStatus(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateStatus_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateStatusRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateStatus(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateStatus_RepoError() {
	actor := superadminActor()
	existing := factories.NewStatusFactory().Make()
	existing.ID = 1
	req := &dto.UpdateStatusRequest{}

	s.repo.On("GetStatusByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateStatus", mock.AnythingOfType("*models.Status")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateStatus(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteStatus_Success() {
	actor := superadminActor()
	existing := factories.NewStatusFactory().Make()
	existing.ID = 1

	s.repo.On("GetStatusByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteStatus", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteStatus(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteStatus_NotFound() {
	actor := superadminActor()

	s.repo.On("GetStatusByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteStatus(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteStatus_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteStatus(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
