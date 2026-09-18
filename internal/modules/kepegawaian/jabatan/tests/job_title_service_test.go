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
// JobTitle, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_Superadmin_Success() {
	req := &dto.CreateJobTitleRequest{Name: "Test JobTitle"}
	actor := superadminActor()

	s.repo.On("CreateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(nil)

	result, err := s.svc.CreateJobTitle(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_Forbidden() {
	req := &dto.CreateJobTitleRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJobTitle(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_RepoError() {
	req := &dto.CreateJobTitleRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJobTitle(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_Success() {
	actor := superadminActor()
	item := factories.NewJobTitleFactory().Make()
	item.ID = 1

	s.repo.On("GetJobTitleByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJobTitleByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJobTitleByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJobTitleByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitle_Success() {
	actor := superadminActor()
	filter := &dto.FilterJobTitleRequest{}
	items := []models.JobTitle{
		*factories.NewJobTitleFactory().Make(),
		*factories.NewJobTitleFactory().Make(),
	}

	s.repo.On("ListJobTitle", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJobTitle(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitle_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJobTitleRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJobTitle(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateJobTitleRequest{Name: &newName}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(nil)

	result, err := s.svc.UpdateJobTitle(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJobTitleRequest{}

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJobTitle(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJobTitleRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJobTitle(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJobTitle", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJobTitle(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJobTitle(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJobTitle(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
