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
// PositionKategori, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePositionKategori_Superadmin_Success() {
	req := &dto.CreatePositionKategoriRequest{Name: "Test PositionKategori"}
	actor := superadminActor()

	s.repo.On("CreatePositionKategori", mock.AnythingOfType("*models.PositionKategori")).Return(nil)

	result, err := s.svc.CreatePositionKategori(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePositionKategori_Forbidden() {
	req := &dto.CreatePositionKategoriRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePositionKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreatePositionKategori_RepoError() {
	req := &dto.CreatePositionKategoriRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreatePositionKategori", mock.AnythingOfType("*models.PositionKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePositionKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionKategoriByID_Success() {
	actor := superadminActor()
	item := factories.NewPositionKategoriFactory().Make()
	item.ID = 1

	s.repo.On("GetPositionKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetPositionKategoriByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionKategoriByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPositionKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetPositionKategoriByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetPositionKategoriByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetPositionKategoriByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListPositionKategori_Success() {
	actor := superadminActor()
	filter := &dto.FilterPositionKategoriRequest{}
	items := []models.PositionKategori{
		*factories.NewPositionKategoriFactory().Make(),
		*factories.NewPositionKategoriFactory().Make(),
	}

	s.repo.On("ListPositionKategori", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPositionKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListPositionKategori_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterPositionKategoriRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListPositionKategori(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePositionKategori_Success() {
	actor := superadminActor()
	existing := factories.NewPositionKategoriFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdatePositionKategoriRequest{Name: &newName}

	s.repo.On("GetPositionKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePositionKategori", mock.AnythingOfType("*models.PositionKategori")).Return(nil)

	result, err := s.svc.UpdatePositionKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePositionKategori_NotFound() {
	actor := superadminActor()
	req := &dto.UpdatePositionKategoriRequest{}

	s.repo.On("GetPositionKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePositionKategori(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdatePositionKategori_Forbidden() {
	actor := regularActor()
	req := &dto.UpdatePositionKategoriRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePositionKategori(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePositionKategori_Success() {
	actor := superadminActor()
	existing := factories.NewPositionKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetPositionKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeletePositionKategori", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeletePositionKategori(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePositionKategori_NotFound() {
	actor := superadminActor()

	s.repo.On("GetPositionKategoriByID", int64(999)).Return(nil, nil)

	err := s.svc.DeletePositionKategori(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeletePositionKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePositionKategori(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
