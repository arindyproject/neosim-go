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
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di jabatan_service_test.go. File ini HANYA menambah skenario test untuk
// Position, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePosition_Superadmin_Success() {
	req := &dto.CreatePositionRequest{Name: "Test Position"}
	actor := superadminActor()

	s.repo.On("CreatePosition", mock.AnythingOfType("*models.Position")).Return(nil)

	result, err := s.svc.CreatePosition(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePosition_Forbidden() {
	req := &dto.CreatePositionRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePosition(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePosition_RepoError() {
	req := &dto.CreatePositionRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreatePosition", mock.AnythingOfType("*models.Position")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePosition(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionByID_Success() {
	actor := superadminActor()
	item := factories.NewPositionFactory().Make()
	item.ID = 1

	s.repo.On("GetPositionByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPositionByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPositionByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetPositionByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetPositionByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListPosition_Success() {
	actor := superadminActor()
	filter := &dto.FilterPositionRequest{}
	items := []models.Position{
		*factories.NewPositionFactory().Make(),
		*factories.NewPositionFactory().Make(),
	}

	s.repo.On("ListPosition", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPosition(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListPosition_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterPositionRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListPosition(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePosition_Success() {
	actor := superadminActor()
	existing := factories.NewPositionFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdatePositionRequest{Name: &newName}

	s.repo.On("GetPositionByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePosition", mock.AnythingOfType("*models.Position")).Return(nil)

	result, err := s.svc.UpdatePosition(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePosition_NotFound() {
	actor := superadminActor()
	req := &dto.UpdatePositionRequest{}

	s.repo.On("GetPositionByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePosition(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePosition_Forbidden() {
	actor := regularActor()
	req := &dto.UpdatePositionRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePosition(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePosition_Success() {
	actor := superadminActor()
	existing := factories.NewPositionFactory().Make()
	existing.ID = 1

	s.repo.On("GetPositionByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePosition", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeletePosition(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePosition_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPositionByID", int64(999)).Return(nil, nil)

	err := s.svc.DeletePosition(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePosition_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePosition(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
