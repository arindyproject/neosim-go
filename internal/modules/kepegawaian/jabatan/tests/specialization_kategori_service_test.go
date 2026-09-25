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
// SpecializationKategori, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_Superadmin_Success() {
	req := &dto.CreateSpecializationKategoriRequest{Name: "Test SpecializationKategori"}
	actor := superadminActor()

	s.repo.On("CreateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.CreateSpecializationKategori(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_Forbidden() {
	req := &dto.CreateSpecializationKategoriRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateSpecializationKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecializationKategori_RepoError() {
	req := &dto.CreateSpecializationKategoriRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateSpecializationKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_Success() {
	actor := superadminActor()
	item := factories.NewSpecializationKategoriFactory().Make()
	item.ID = 1

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationKategoriByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetSpecializationKategoriByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecializationKategori_Success() {
	actor := superadminActor()
	filter := &dto.FilterSpecializationKategoriRequest{}
	items := []models.SpecializationKategori{
		*factories.NewSpecializationKategoriFactory().Make(),
		*factories.NewSpecializationKategoriFactory().Make(),
	}

	s.repo.On("ListSpecializationKategori", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListSpecializationKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecializationKategori_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterSpecializationKategoriRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListSpecializationKategori(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateSpecializationKategoriRequest{Name: &newName}

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateSpecializationKategori", mock.AnythingOfType("*models.SpecializationKategori")).Return(nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateSpecializationKategoriRequest{}

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateSpecializationKategori(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecializationKategori_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateSpecializationKategoriRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateSpecializationKategori(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetSpecializationKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteSpecializationKategori", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteSpecializationKategori(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationKategoriByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteSpecializationKategori(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecializationKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteSpecializationKategori(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
