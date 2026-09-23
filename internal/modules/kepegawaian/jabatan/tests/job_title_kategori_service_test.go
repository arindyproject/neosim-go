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
// JobTitleKategori, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleKategori_Superadmin_Success() {
	req := &dto.CreateJobTitleKategoriRequest{Name: "Test JobTitleKategori"}
	actor := superadminActor()

	s.repo.On("CreateJobTitleKategori", mock.AnythingOfType("*models.JobTitleKategori")).Return(nil)

	result, err := s.svc.CreateJobTitleKategori(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleKategori_Forbidden() {
	req := &dto.CreateJobTitleKategoriRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJobTitleKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleKategori_RepoError() {
	req := &dto.CreateJobTitleKategoriRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateJobTitleKategori", mock.AnythingOfType("*models.JobTitleKategori")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJobTitleKategori(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleKategoriByID_Success() {
	actor := superadminActor()
	item := factories.NewJobTitleKategoriFactory().Make()
	item.ID = 1

	s.repo.On("GetJobTitleKategoriByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJobTitleKategoriByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleKategoriByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJobTitleKategoriByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleKategoriByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJobTitleKategoriByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitleKategori_Success() {
	actor := superadminActor()
	filter := &dto.FilterJobTitleKategoriRequest{}
	items := []models.JobTitleKategori{
		*factories.NewJobTitleKategoriFactory().Make(),
		*factories.NewJobTitleKategoriFactory().Make(),
	}

	s.repo.On("ListJobTitleKategori", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJobTitleKategori(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitleKategori_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJobTitleKategoriRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJobTitleKategori(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleKategori_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleKategoriFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateJobTitleKategoriRequest{Name: &newName}

	s.repo.On("GetJobTitleKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitleKategori", mock.AnythingOfType("*models.JobTitleKategori")).Return(nil)

	result, err := s.svc.UpdateJobTitleKategori(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleKategori_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJobTitleKategoriRequest{}

	s.repo.On("GetJobTitleKategoriByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJobTitleKategori(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleKategori_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJobTitleKategoriRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJobTitleKategori(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleKategori_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleKategoriFactory().Make()
	existing.ID = 1

	s.repo.On("GetJobTitleKategoriByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJobTitleKategori", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJobTitleKategori(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleKategori_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleKategoriByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJobTitleKategori(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleKategori_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJobTitleKategori(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
