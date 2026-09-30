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
// (superadminActor, regularActor, mockNoPermissions, ptrTo, dst) sudah
// didefinisikan di pegawai_service_test.go. File ini HANYA menambah skenario
// test untuk Jenis, memakai s.svc / s.repo yang SAMA.
//
// CATATAN PENTING soal isi test ini: saya belum melihat isi asli
// jenis_service.go untuk modul pegawai. Test di bawah HANYA menyesuaikan
// field Code/Label (dari Name lama) dan memakai pola dasar create/get/list/
// update/delete yang sudah terbukti benar di modul lain (permission check,
// GetByID sebelum update/delete, dst). Saya TIDAK menambahkan skenario
// duplikat Code/Label atau validasi lain (seperti CheckJobTitleKategori di
// modul jabatan) karena tidak tahu apakah CreateJenis/UpdateJenis di modul
// ini benar-benar memvalidasi itu. Kalau ternyata service ini juga memanggil
// GetJenisByCode/GetJenisByLabel sebelum create/update (mengikuti pola
// SpecializationKategori), test create/update di bawah akan panic dengan
// "unexpected call" persis seperti kasus-kasus sebelumnya — kirim isi
// jenis_service.go kalau itu terjadi, supaya saya sesuaikan tanpa tebakan.

func newCreateJenisReq() *dto.CreateJenisRequest {
	return &dto.CreateJenisRequest{
		Code:  "TEST001",
		Label: "Test Tipe",
	}
}

func (s *KepegawaianPegawaiServiceTestSuite) mockJenisNoDuplicate(code, label string) {
	s.repo.On("GetJenisByCode", code).Return(nil, nil)
	s.repo.On("GetJenisByLabel", label).Return(nil, nil)
}

// ── Create ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateJenis_Superadmin_Success() {
	req := newCreateJenisReq()
	actor := superadminActor()

	s.mockJenisNoDuplicate(req.Code, req.Label)
	s.repo.On("CreateJenis", mock.AnythingOfType("*models.Jenis")).Return(nil)

	result, err := s.svc.CreateJenis(context.Background(), req, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateJenis_Forbidden() {
	req := newCreateJenisReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJenis(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateJenis", mock.Anything)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_CreateJenis_RepoError() {
	req := newCreateJenisReq()
	actor := superadminActor()

	s.repo.On("GetJenisByCode", req.Code).Return(nil, nil)
	s.repo.On("GetJenisByLabel", req.Label).Return(nil, nil)
	s.repo.On("CreateJenis", mock.AnythingOfType("*models.Jenis")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJenis(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetJenisByID_Success() {
	actor := superadminActor()
	item := factories.NewJenisFactory().Make()
	item.ID = 1

	s.repo.On("GetJenisByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJenisByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetJenisByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJenisByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJenisByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_GetJenisByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJenisByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── List ─────────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListJenis_Success() {
	actor := superadminActor()
	filter := &dto.FilterJenisRequest{}
	items := []models.Jenis{
		*factories.NewJenisFactory().Make(),
		*factories.NewJenisFactory().Make(),
	}

	s.repo.On("ListJenis", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJenis(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListJenis_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJenisRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJenis(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_ListJenis_RepoError() {
	actor := superadminActor()
	filter := &dto.FilterJenisRequest{}

	s.repo.On("ListJenis", 1, 10, filter).Return(nil, int64(0), fmt.Errorf("db error"))

	result, total, err := s.svc.ListJenis(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateJenis_Success() {
	actor := superadminActor()
	existing := factories.NewJenisFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateJenisRequest{Label: &newLabel}

	s.repo.On("GetJenisByID", int64(1)).Return(existing, nil)
	s.repo.On("GetJenisByLabel", newLabel).Return(nil, nil)
	s.repo.On("UpdateJenis", mock.AnythingOfType("*models.Jenis")).Return(nil)

	result, err := s.svc.UpdateJenis(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateJenis_PartialFields() {
	actor := superadminActor()
	existing := factories.NewJenisFactory().Make()
	existing.ID = 1
	originalCode := existing.Code
	newLabel := "Label Baru"
	req := &dto.UpdateJenisRequest{Label: &newLabel}

	s.repo.On("GetJenisByID", int64(1)).Return(existing, nil)
	s.repo.On("GetJenisByLabel", newLabel).Return(nil, nil)
	s.repo.On("UpdateJenis", mock.MatchedBy(func(m *models.Jenis) bool {
		return m.Code == originalCode && m.Label == newLabel
	})).Return(nil)

	result, err := s.svc.UpdateJenis(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalCode, result.Code)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateJenis_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJenisRequest{}

	s.repo.On("GetJenisByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJenis(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateJenis_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJenisRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJenis(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_UpdateJenis_RepoError() {
	actor := superadminActor()
	existing := factories.NewJenisFactory().Make()
	existing.ID = 1
	req := &dto.UpdateJenisRequest{}

	s.repo.On("GetJenisByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJenis", mock.AnythingOfType("*models.Jenis")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateJenis(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteJenis_Success() {
	actor := superadminActor()
	existing := factories.NewJenisFactory().Make()
	existing.ID = 1

	s.repo.On("GetJenisByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJenis", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJenis(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteJenis_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJenisByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJenis(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPegawaiServiceTestSuite) Test_DeleteJenis_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJenis(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
