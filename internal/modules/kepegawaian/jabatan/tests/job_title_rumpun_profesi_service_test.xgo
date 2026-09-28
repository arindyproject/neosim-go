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
// JobTitleRumpunProfesi, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module jabatan).

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleRumpunProfesi_Superadmin_Success() {
	fhir := "dokter"
	req := &dto.CreateJobTitleRumpunProfesiRequest{
		Code:     "DOKTER",
		Label:    "Dokter",
		FHIRCode: &fhir,
	}
	actor := superadminActor()

	s.repo.On("CreateJobTitleRumpunProfesi", mock.AnythingOfType("*models.JobTitleRumpunProfesi")).Return(nil)

	result, err := s.svc.CreateJobTitleRumpunProfesi(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
	s.Equal(req.FHIRCode, result.FHIRCode)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleRumpunProfesi_Forbidden() {
	req := &dto.CreateJobTitleRumpunProfesiRequest{Code: "DOKTER", Label: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJobTitleRumpunProfesi(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateJobTitleRumpunProfesi", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitleRumpunProfesi_RepoError() {
	req := &dto.CreateJobTitleRumpunProfesiRequest{Code: "DOKTER", Label: "Dokter"}
	actor := superadminActor()

	s.repo.On("CreateJobTitleRumpunProfesi", mock.AnythingOfType("*models.JobTitleRumpunProfesi")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJobTitleRumpunProfesi(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleRumpunProfesiByID_Success() {
	actor := superadminActor()
	item := factories.NewJobTitleRumpunProfesiFactory().Make()
	item.ID = 1

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJobTitleRumpunProfesiByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleRumpunProfesiByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJobTitleRumpunProfesiByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleRumpunProfesiByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJobTitleRumpunProfesiByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitleRumpunProfesi_Success() {
	actor := superadminActor()
	filter := &dto.FilterJobTitleRumpunProfesiRequest{}
	items := []models.JobTitleRumpunProfesi{
		*factories.NewJobTitleRumpunProfesiFactory().Make(),
		*factories.NewJobTitleRumpunProfesiFactory().Make(),
	}

	s.repo.On("ListJobTitleRumpunProfesi", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJobTitleRumpunProfesi(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitleRumpunProfesi_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJobTitleRumpunProfesiRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJobTitleRumpunProfesi(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleRumpunProfesi_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleRumpunProfesiFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateJobTitleRumpunProfesiRequest{Label: &newLabel}

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitleRumpunProfesi", mock.AnythingOfType("*models.JobTitleRumpunProfesi")).Return(nil)

	result, err := s.svc.UpdateJobTitleRumpunProfesi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleRumpunProfesi_PartialFields() {
	actor := superadminActor()
	existing := factories.NewJobTitleRumpunProfesiFactory().Make()
	existing.ID = 1
	originalCode := existing.Code
	newLabel := "Label Baru"
	req := &dto.UpdateJobTitleRumpunProfesiRequest{Label: &newLabel}

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitleRumpunProfesi", mock.MatchedBy(func(m *models.JobTitleRumpunProfesi) bool {
		return m.Code == originalCode && m.Label == newLabel
	})).Return(nil)

	result, err := s.svc.UpdateJobTitleRumpunProfesi(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(originalCode, result.Code)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleRumpunProfesi_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJobTitleRumpunProfesiRequest{}

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJobTitleRumpunProfesi(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitleRumpunProfesi_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJobTitleRumpunProfesiRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJobTitleRumpunProfesi(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleRumpunProfesi_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleRumpunProfesiFactory().Make()
	existing.ID = 1

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJobTitleRumpunProfesi", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJobTitleRumpunProfesi(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleRumpunProfesi_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleRumpunProfesiByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJobTitleRumpunProfesi(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitleRumpunProfesi_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJobTitleRumpunProfesi(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
