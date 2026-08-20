package tests

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianIdentifierServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di identifier_service_test.go. File ini HANYA menambah skenario test untuk
// Tipe, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module identifier).
func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateTipe_Superadmin_Success() {
	penerbit := "Dukcapil / Kemendagri"
	fhirSystem := "https://fhir.kemkes.go.id/id/nik"
	req := &dto.CreateTipeRequest{
		Code:       "NIK",
		Label:      "NIK",
		Penerbit:   &penerbit,
		FHIRSystem: &fhirSystem,
		HasExpiry:  false,
		IsNakes:    false,
		IsRequired: true,
	}
	actor := superadminActor()

	// Mock pengecekan keunikan (Code & Label belum digunakan)
	s.repo.On("GetTipeByCode", req.Code).Return(nil, nil)
	s.repo.On("GetTipeByLabel", req.Label).Return(nil, nil)

	s.repo.On("CreateTipe", mock.AnythingOfType("*models.Tipe")).Return(nil)

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateTipe_Forbidden() {
	req := &dto.CreateTipeRequest{
		Code:  "NIK",
		Label: "NIK",
	}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_CreateTipe_RepoError() {
	req := &dto.CreateTipeRequest{
		Code:  "TEST_" + randomString(8),
		Label: "TEST_" + randomString(8),
	}
	actor := superadminActor()

	// 1. Mock pengecekan code & label (mengembalikan nil agar lolos validasi keunikan)
	s.repo.On("GetTipeByCode", req.Code).Return(nil, nil)
	s.repo.On("GetTipeByLabel", req.Label).Return(nil, nil)

	// 2. Mock error saat insert ke DB
	s.repo.On("CreateTipe", mock.AnythingOfType("*models.Tipe")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func randomString(i int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, i)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetTipeByID_Success() {
	actor := superadminActor()
	item := factories.NewTipeFactory().Make()
	item.ID = 1

	s.repo.On("GetTipeByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetTipeByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetTipeByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetTipeByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_GetTipeByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetTipeByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListTipe_Success() {
	actor := superadminActor()
	filter := &dto.FilterTipeRequest{}
	items := []models.Tipe{
		*factories.NewTipeFactory().Make(),
		*factories.NewTipeFactory().Make(),
	}

	s.repo.On("ListTipe", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListTipe(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_ListTipe_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterTipeRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListTipe(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateTipe_Success() {
	actor := superadminActor()
	existing := factories.NewTipeFactory().Make()
	existing.ID = 1

	newCode := "NIK_UPDATED"
	newLabel := "NIK Perbarui"
	req := &dto.UpdateTipeRequest{
		Code:  &newCode,
		Label: &newLabel,
	}

	// 1. Ambil data eksis
	s.repo.On("GetTipeByID", int64(1)).Return(existing, nil)

	// 2. Mock validasi keunikan Code & Label baru (kembalikan nil agar dianggap tidak ada bentrok)
	s.repo.On("GetTipeByCode", newCode).Return(nil, nil)
	s.repo.On("GetTipeByLabel", newLabel).Return(nil, nil)

	// 3. Update ke DB
	s.repo.On("UpdateTipe", mock.AnythingOfType("*models.Tipe")).Return(nil)

	result, err := s.svc.UpdateTipe(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newCode, result.Code)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateTipe_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateTipeRequest{}

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateTipe(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_UpdateTipe_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateTipeRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateTipe(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteTipe_Success() {
	actor := superadminActor()
	existing := factories.NewTipeFactory().Make()
	existing.ID = 1

	s.repo.On("GetTipeByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteTipe", int64(1)).Return(nil)

	err := s.svc.DeleteTipe(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteTipe_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteTipe(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianIdentifierServiceTestSuite) Test_DeleteTipe_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteTipe(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
