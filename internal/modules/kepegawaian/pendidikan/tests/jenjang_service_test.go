package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	"neosim_go/internal/modules/kepegawaian/pendidikan/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianPendidikanServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di pendidikan_service_test.go. File ini HANYA menambah skenario test untuk
// Jenjang, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module pendidikan).

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreateJenjang_Superadmin_Success() {
	req := &dto.CreateJenjangRequest{Code: "TEST", Label: "Test Jenjang"}
	actor := superadminActor()

	s.repo.On("CreateJenjang", mock.AnythingOfType("*models.Jenjang")).Return(nil)

	result, err := s.svc.CreateJenjang(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreateJenjang_Forbidden() {
	req := &dto.CreateJenjangRequest{Code: "TEST", Label: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJenjang(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_CreateJenjang_RepoError() {
	req := &dto.CreateJenjangRequest{Code: "TEST", Label: "Test"}
	actor := superadminActor()

	s.repo.On("CreateJenjang", mock.AnythingOfType("*models.Jenjang")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJenjang(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetJenjangByID_Success() {
	actor := superadminActor()
	item := factories.NewJenjangFactory().Make()
	item.ID = 1

	s.repo.On("GetJenjangByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJenjangByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetJenjangByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJenjangByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJenjangByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_GetJenjangByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJenjangByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListJenjang_Success() {
	actor := superadminActor()
	filter := &dto.FilterJenjangRequest{}
	items := []models.Jenjang{
		*factories.NewJenjangFactory().Make(),
		*factories.NewJenjangFactory().Make(),
	}

	s.repo.On("ListJenjang", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJenjang(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_ListJenjang_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJenjangRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJenjang(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdateJenjang_Success() {
	actor := superadminActor()
	existing := factories.NewJenjangFactory().Make()
	existing.ID = 1
	newLabel := "Updated Name"
	req := &dto.UpdateJenjangRequest{Label: &newLabel}

	s.repo.On("GetJenjangByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJenjang", mock.AnythingOfType("*models.Jenjang")).Return(nil)

	result, err := s.svc.UpdateJenjang(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdateJenjang_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJenjangRequest{}

	s.repo.On("GetJenjangByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJenjang(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_UpdateJenjang_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJenjangRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJenjang(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeleteJenjang_Success() {
	actor := superadminActor()
	existing := factories.NewJenjangFactory().Make()
	existing.ID = 1

	s.repo.On("GetJenjangByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJenjang", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJenjang(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeleteJenjang_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJenjangByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJenjang(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianPendidikanServiceTestSuite) Test_DeleteJenjang_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJenjang(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
