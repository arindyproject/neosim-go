package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	"neosim_go/internal/modules/kepegawaian/kontak/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianKontakServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di kontak_service_test.go. File ini HANYA menambah skenario test untuk
// Tipe, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module kontak).

func (s *KepegawaianKontakServiceTestSuite) Test_CreateTipe_Superadmin_Success() {
	req := &dto.CreateTipeRequest{Code: "TEST001", Label: "Test Tipe"}
	actor := superadminActor()

	s.repo.On("CreateTipe", mock.AnythingOfType("*models.Tipe")).Return(nil)

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateTipe_Forbidden() {
	req := &dto.CreateTipeRequest{Code: "TEST001", Label: "Test Tipe"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianKontakServiceTestSuite) Test_CreateTipe_RepoError() {
	req := &dto.CreateTipeRequest{Code: "TEST001", Label: "Test Tipe"}
	actor := superadminActor()

	s.repo.On("CreateTipe", mock.AnythingOfType("*models.Tipe")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateTipe(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetTipeByID_Success() {
	actor := superadminActor()
	item := factories.NewTipeFactory().Make()
	item.ID = 1

	s.repo.On("GetTipeByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetTipeByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetTipeByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetTipeByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_GetTipeByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetTipeByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_ListTipe_Success() {
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

func (s *KepegawaianKontakServiceTestSuite) Test_ListTipe_Forbidden() {
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

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateTipe_Success() {
	actor := superadminActor()
	existing := factories.NewTipeFactory().Make()
	existing.ID = 1
	newCode := "UPDATED001"
	newLabel := "Updated Label"
	req := &dto.UpdateTipeRequest{Code: &newCode, Label: &newLabel}

	s.repo.On("GetTipeByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateTipe", mock.AnythingOfType("*models.Tipe")).Return(nil)

	result, err := s.svc.UpdateTipe(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newCode, result.Code)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateTipe_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateTipeRequest{}

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateTipe(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_UpdateTipe_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateTipeRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateTipe(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteTipe_Success() {
	actor := superadminActor()
	existing := factories.NewTipeFactory().Make()
	existing.ID = 1

	s.repo.On("GetTipeByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteTipe", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteTipe(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteTipe_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTipeByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteTipe(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianKontakServiceTestSuite) Test_DeleteTipe_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteTipe(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
