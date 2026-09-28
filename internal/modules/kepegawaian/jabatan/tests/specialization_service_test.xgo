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
// Specialization, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_Superadmin_Success() {
	req := &dto.CreateSpecializationRequest{Name: "Test Specialization"}
	actor := superadminActor()

	s.repo.On("CreateSpecialization", mock.AnythingOfType("*models.Specialization")).Return(nil)

	result, err := s.svc.CreateSpecialization(context.Background(),req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_Forbidden() {
	req := &dto.CreateSpecializationRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateSpecialization(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_RepoError() {
	req := &dto.CreateSpecializationRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateSpecialization", mock.AnythingOfType("*models.Specialization")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateSpecialization(context.Background(),req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationByID_Success() {
	actor := superadminActor()
	item := factories.NewSpecializationFactory().Make()
	item.ID = 1

	s.repo.On("GetSpecializationByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetSpecializationByID(context.Background(),1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetSpecializationByID(context.Background(),999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetSpecializationByID(context.Background(),1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecialization_Success() {
	actor := superadminActor()
	filter := &dto.FilterSpecializationRequest{}
	items := []models.Specialization{
		*factories.NewSpecializationFactory().Make(),
		*factories.NewSpecializationFactory().Make(),
	}

	s.repo.On("ListSpecialization", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListSpecialization(context.Background(),1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecialization_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterSpecializationRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListSpecialization(context.Background(),1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecialization_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateSpecializationRequest{Name: &newName}

	s.repo.On("GetSpecializationByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateSpecialization", mock.AnythingOfType("*models.Specialization")).Return(nil)

	result, err := s.svc.UpdateSpecialization(context.Background(),1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecialization_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateSpecializationRequest{}

	s.repo.On("GetSpecializationByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateSpecialization(context.Background(),999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecialization_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateSpecializationRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateSpecialization(context.Background(),1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecialization_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationFactory().Make()
	existing.ID = 1

	s.repo.On("GetSpecializationByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteSpecialization", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteSpecialization(context.Background(),1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecialization_NotFound() {
	actor := superadminActor()

	s.repo.On("GetSpecializationByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteSpecialization(context.Background(),999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecialization_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteSpecialization(context.Background(),1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
